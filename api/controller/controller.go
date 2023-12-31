package controller

import (
	"api/models"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	provider "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

var key = []byte("aBcDeFgHiJkLmNOpQrStUvWxYz")
var Tracer trace.Tracer

func SetTracerProvider(t *provider.TracerProvider) {
	fmt.Println("SetTracerProvider")

	Tracer = t.Tracer("testing-shit")
}

func PatchUser(c *gin.Context) {
	ctx, sp := Tracer.Start(c.Request.Context(), "patch-user")
	defer sp.End()

	db := getDatabase(ctx)

	ctx, sp = Tracer.Start(ctx, "loading-user")
	defer sp.End()

	email := utils.GetEmailFromContext(c)

	var user models.User
	db.WithContext(ctx).Where("email = ?", email).First(&user)

	if err := c.BindJSON(&user); err != nil {
		sp.SetStatus(codes.Error, err.Error())
		c.String(http.StatusBadRequest, "Invalid data")
	} else {
		if email != user.Email {
			sp.SetStatus(codes.Error, "Mismatched email")
			c.String(http.StatusBadRequest, "Mismatched email")
		} else {
			db.WithContext(ctx).Save(&user)
			c.IndentedJSON(http.StatusOK, user)
		}
	}
}

func GetUser(c *gin.Context) {
	ctx, sp := Tracer.Start(c.Request.Context(), "get-user")
	defer sp.End()

	db := getDatabase(ctx)

	ctx, sp = Tracer.Start(ctx, "loading-user")
	defer sp.End()

	email := utils.GetEmailFromContext(c)

	userChan := make(chan *models.User)
	go loadUserByEmail(db, ctx, email, userChan)
	user := <-userChan

	c.IndentedJSON(http.StatusOK, user)
}

func PostDonation(c *gin.Context) {
	ctx, span := Tracer.Start(c.Request.Context(), "donate")
	defer span.End()

	var payload donationPayload
	if err := c.BindJSON(&payload); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	span.SetAttributes(attribute.String("issue", payload.Issue.String()), attribute.Float64("amount", float64(payload.Amount)))

	db := getDatabase(ctx)
	email := utils.GetEmailFromContext(c)

	userChan := make(chan *models.User)
	go loadUserByEmail(db, ctx, email, userChan)

	issueChan := make(chan *models.Issue)
	go loadIssueByGUID(db, ctx, payload.Issue, issueChan)
	user, issue := <-userChan, <-issueChan

	if issue.ID == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("Issue not found: %s", payload.Issue.String())))
		return
	}
	if user.ID == 0 {
		c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("User not found: %s", email)))
		return
	}

	donation := models.Donation{
		GUID:   uuid.New(),
		User:   user.ID,
		Issue:  issue.ID,
		Amount: payload.Amount,
		Status: "created",
	}

	ctx, span = Tracer.Start(ctx, "processing-donation")
	defer span.End()

	transaction := db.WithContext(ctx).Begin()
	transaction.WithContext(ctx).Save(&donation)
	if transaction.Error != nil {
		rollbackAndReturn(transaction, ctx, span, c)
		return
	}
	transaction.WithContext(ctx).Model(&models.Issue{}).Where("ID", issue.ID).Update("pool", gorm.Expr("pool + ?", donation.Amount))
	if transaction.Error != nil {
		rollbackAndReturn(transaction, ctx, span, c)
		return
	}
	transaction.WithContext(ctx).Commit()

	c.IndentedJSON(http.StatusCreated, donation)
}

func DeleteDonation(c *gin.Context) {
	ctx, span := Tracer.Start(c.Request.Context(), "cancel-donation")
	defer span.End()

	db := getDatabase(ctx)
	if guid, err := uuid.Parse(c.Param("guid")); err == nil {
		donationChan := make(chan *models.Donation)
		go loadDonationByGUID(db, ctx, guid, donationChan)
		donation := <-donationChan

		if donation.ID == 0 {
			c.AbortWithError(http.StatusBadRequest, errors.New(fmt.Sprintf("No donation found: %s", guid.String())))
			return
		}
		if donation.Deleted {
			c.IndentedJSON(http.StatusOK, donation)
			return
		}

		ctx, span = Tracer.Start(ctx, "canceling-donation")
		defer span.End()

		donation.Status = "deleted"
		donation.Deleted = true

		transaction := db.WithContext(ctx).Begin()
		transaction.WithContext(ctx).Save(&donation)
		if transaction.Error != nil {
			rollbackAndReturn(transaction, ctx, span, c)
			return
		}
		transaction.WithContext(ctx).Model(&models.Issue{}).Where("ID", donation.Issue).Update("pool", gorm.Expr("pool - ?", donation.Amount))
		if transaction.Error != nil {
			rollbackAndReturn(transaction, ctx, span, c)
			return
		}
		transaction.WithContext(ctx).Commit()

		c.String(http.StatusOK, "OK")
		c.IndentedJSON(http.StatusOK, donation)
	} else {
		span.AddEvent("error", trace.WithAttributes(
			attribute.String("message", err.Error()),
		))
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadGateway)
	}
}

func GetDonations(c *gin.Context) {
	ctx, span := Tracer.Start(c.Request.Context(), "get-donations")
	defer span.End()

	db := getDatabase(ctx)
	email := utils.GetEmailFromContext(c)

	userChan := make(chan *models.User)
	go loadUserByEmail(db, ctx, email, userChan)
	user := <-userChan

	page := getQueryParam(c, "page", "1", 1, 100)
	per := getQueryParam(c, "per", "10", 1, 50)
	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("per", per),
	)

	var donations []models.Donation
	if ret := db.WithContext(ctx).Where("user = ? AND deleted = false", user.ID).Limit(per).Offset((page-1)*per).Order("id desc").Select("guid", "status", "amount", "created_at").Find(&donations); ret.Error != nil {
		c.AbortWithError(http.StatusBadGateway, ret.Error)
	} else {
		c.IndentedJSON(http.StatusOK, donations)
	}
}

func GetIssues(c *gin.Context) {
	ctx, span := Tracer.Start(c.Request.Context(), "get-issues")
	defer span.End()

	db := getDatabase(ctx)

	page := getQueryParam(c, "page", "1", 1, 100)
	per := getQueryParam(c, "per", "10", 1, 50)
	span.SetAttributes(
		attribute.Int("page", page),
		attribute.Int("per", per),
	)

	var donations []models.Issue
	if ret := db.WithContext(ctx).Limit(per).Offset((page - 1) * per).Order("id desc").Find(&donations); ret.Error != nil {
		c.AbortWithError(http.StatusBadGateway, ret.Error)
	} else {
		c.IndentedJSON(http.StatusOK, donations)
	}
}

func HealthCheck(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func GenerateToken(c *gin.Context) {
	ctx, sp := Tracer.Start(c.Request.Context(), "generate-token")
	defer sp.End()

	db := getDatabase(ctx)

	ctx, sp = Tracer.Start(ctx, "loading-user")
	defer sp.End()

	var user models.User
	db.WithContext(ctx).First(&user)

	s, err := createToken(ctx, user.Email)
	if err != nil {
		c.AbortWithError(http.StatusBadGateway, errors.New("Error generating token"))
	} else {
		c.String(http.StatusOK, s)
	}
}

func loadUserByEmail(db *gorm.DB, ctx context.Context, email string, c chan *models.User) {
	ctx, span := Tracer.Start(ctx, "loading-user-by-email")
	defer span.End()

	span.SetAttributes(attribute.String("email", email))

	var user *models.User
	db.WithContext(ctx).Where("email = ?", email).First(&user)

	c <- user
}

func loadDonationByGUID(db *gorm.DB, ctx context.Context, guid uuid.UUID, c chan *models.Donation) {
	ctx, span := Tracer.Start(ctx, "loading-transaction-by-guid")
	defer span.End()

	span.SetAttributes(attribute.String("guid", guid.String()))

	var transaction *models.Donation
	db.WithContext(ctx).Where("guid = ?", guid).First(&transaction)

	c <- transaction
}

func loadIssueByGUID(db *gorm.DB, ctx context.Context, guid uuid.UUID, c chan *models.Issue) {
	ctx, span := Tracer.Start(ctx, "loading-issue-by-guid")
	defer span.End()

	span.SetAttributes(attribute.String("guid", guid.String()))

	var issue *models.Issue
	db.WithContext(ctx).Where("guid = ?", guid).First(&issue)

	c <- issue
}

func createSpan(ctx context.Context, name string) (context.Context, trace.Span) {
	return Tracer.Start(ctx, name)
}

func getDatabase(ctx context.Context) *gorm.DB {
	_, span := createSpan(ctx, "getting-database")
	defer span.End()

	db, err := models.Database(false)
	if err != nil {
		panic(err)
	}

	return db
}

func createToken(ctx context.Context, email string) (string, error) {
	ctx, sp := Tracer.Start(ctx, "creating-token")
	defer sp.End()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": email,
	})

	return token.SignedString(key)
}

func rollbackAndReturn(transaction *gorm.DB, ctx context.Context, span trace.Span, c *gin.Context) {
	fmt.Println(transaction.Error)

	transaction.WithContext(ctx).Rollback()

	span.AddEvent("transaction rolledback", trace.WithAttributes(
		attribute.String("message", transaction.Error.Error()),
	))

	c.AbortWithStatus(http.StatusBadGateway)
}

func getQueryParam(c *gin.Context, name string, def string, min int, max int) int {
	ret, err := strconv.Atoi(c.DefaultQuery(name, def))
	if err != nil {
		return 0
	}

	switch {
	case ret > max:
		return max
	case ret < min:
		return min
	default:
		return ret
	}
}
