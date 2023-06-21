package dbhelper

import (
	"fmt"

	"userGrpc/model"
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
            VALUES ($1,$2,$3,$4,$5)
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

// GetUserById get user by id
func (d DAO) GetUserById(userId string) (*model.Register, error) {
	// SQL:=language
	SQL := `Select 
                  name,
                   address,
                   phone,
                   email,
                   dob
             From users
             Where id=$1::uuid
             `

	var user model.Register
	err := d.DB.Get(&user, SQL, userId)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUsers get user by id
func (d DAO) GetUsers() ([]model.Register, error) {
	// SQL:=language
	SQL := `Select 
                   id,
                   name,
                   address,
                   phone,
                   email,
                   dob
             From users
             `

	var user []model.Register
	err := d.DB.Select(&user, SQL)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (d DAO) GetUserByEmailAndPhone(phone, email string) (*model.Register, error) {
	// SQL=language
	SQL := `Select 
                   id,
                   name,
                   address,
                   phone,
                   email,
                   dob
             From users
             where phone=$1 and email=$2
             `

	args := []interface{}{
		phone,
		email,
	}
	var user model.Register
	err := d.DB.Get(&user, SQL, args...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser update user by id with given fields
func (d DAO) UpdateUser(user model.Register, userId string) error {
	// SQL query
	SQL := `
		UPDATE users
		SET
		    name = COALESCE($1, name),
			address = COALESCE($2, address),
			updated_at = NOW()
		WHERE id = $3::uuid
	`
	args := []interface{}{
		user.Name,
		user.Address,
		userId,
	}
	_, err := d.DB.Exec(SQL, args...)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser delete user by given user id
func (d DAO) DeleteUser(userId string) error {
	// SQL=language
	SQL := `
		UPDATE users
		SET
			archived_at = NOW()
		WHERE id = $1::uuid
	`

	_, err := d.DB.Exec(SQL, userId)
	if err != nil {
		return err
	}

	return nil
}
