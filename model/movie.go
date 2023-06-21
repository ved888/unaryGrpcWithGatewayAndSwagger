package model

type Movie struct {
	ID    string `json:"id" db:"id"`
	Title string `json:"title" db:"title"`
	Genre string `json:"genre" db:"genre"`
}
