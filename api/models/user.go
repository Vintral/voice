package models

type User struct {
	BaseModel

	Email     string `gorm:"unique"`
	FirstName string
	LastName  string
	Address1  string
	Address2  string
	City      string
	State     string
	Zip       string
	Admin     bool `gorm:"default:false"`
}
