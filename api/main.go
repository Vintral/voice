package main

import (
	"api/controller"
	"api/models"
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan:
			db, err := models.Database(false)
			if err != nil {
				sql, err := db.DB()
				if err != nil {
					sql.Close()
				}
			}

			cancel()
			os.Exit(1)
		}
	}()

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println(err)
	}

	db, err := models.Database(false)
	if err != nil {
		panic(err)
	}
	models.RunMigrations(db)

	fmt.Println("Env Variables:")
	fmt.Println(os.Environ())

	otelShutdown, tp, err := setupOTelSDK(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		fmt.Println("In shutdown")
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	controller.SetTracerProvider(tp)

	PORT := os.Getenv("PORT")
	HOST := os.Getenv("HOST")
	router := setupRoutes()
	router.Run(HOST + ":" + PORT)
	fmt.Println("Running on: " + HOST + ":" + PORT)
}

func setupRoutes() (router *gin.Engine) {
	router = gin.Default()
	router.Use(otelgin.Middleware(getEnv("OTEL_SERVICE_NAME", "default-service")))
	router.SetTrustedProxies(nil)

	router.GET("/", controller.HealthCheck)
	router.GET("/auth", controller.GenerateToken)

	adminRoutes := router.Group("/admin", Authorizer(), AdminAuthorizer())
	{
		adminRoutes.GET("/issues", controller.GetIssues)
	}

	userRoutes := router.Group("/user", Authorizer())
	{
		userRoutes.GET("/user", controller.GetUser)
		userRoutes.PATCH("/user", controller.PatchUser)

		userRoutes.POST("/user/donation", controller.PostDonation)
		userRoutes.GET("/user/donations", controller.GetDonations)
		userRoutes.DELETE("/user/donation/:guid", controller.DeleteDonation)
	}

	return
}
