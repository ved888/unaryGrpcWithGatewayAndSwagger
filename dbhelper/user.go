package dbhelper

import (
	"fmt"
	"gRPC2/model"
)

// CreateUser create user by given payload
func (d DAO) CreateUser(user *model.Register) (*string, error) {
	//SQL==language
	SQL := `INSERT INTO users (
                   name,
                   address,
                   phone,
                   email,
                   dob
                   ) 
            VALUES ($1, $2,$3,$4,$5)
            RETURNING id::uuid ;`

	arg := []interface{}{
		user.Name,
		user.Address,
		user.Phone,
		user.Email,
		user.DOB,
	}
	var userId string

	err := d.DB.QueryRowx(SQL, arg...).Scan(&userId)
	if err != nil {
		fmt.Println("error in userId:", err.Error())
		return nil, err
	}

	return &userId, nil
}
