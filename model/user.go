package model

type Register struct {
	Id      string `json:"userId" db:"id"`
	Name    string `json:"Name" db:"name"`
	Address string `json:"Address" db:"address"`
	Phone   string `json:"Phone" db:"phone"`
	Email   string `json:"Email" db:"email"`
	DOB     string `json:"DOB" db:"dob"`
	Image   string `json:"image" db:"image_id"`
}
