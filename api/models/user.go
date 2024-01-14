package models

type User struct {
	BaseModel

	Email     string `gorm:"uniqueIndex,size:64" json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Address1  string `json:"address_1"`
	Address2  string `json:"address_2"`
	City      string `json:"city"`
	State     string `json:"state"`
	Zip       string `json:"zip"`
	Admin     bool   `gorm:"default:false" json:"-"`
}
