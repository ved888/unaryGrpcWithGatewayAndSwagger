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

type Image struct {
	Id         string `json:"id" db:"id"`
	BucketName string `json:"bucketName" db:"bucket_name"`
	Path       string `json:"path" db:"path"`
}

type Login struct {
	UserPhone string `json:"UserPhone" db:"phone"`
	UserEmail string `json:"UserEmail" db:"email"`
}

type Update struct {
	Name    string `json:"Name" db:"name"`
	Address string `json:"Address" db:"address"`
}
