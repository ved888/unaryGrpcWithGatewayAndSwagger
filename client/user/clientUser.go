package user

import (
	"context"
	"encoding/json"
	"fmt"
	"gRPC2/middleware"
	"gRPC2/model"
	"gRPC2/pb/pb"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"net/mail"
)

var Conn2, _ = grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))

var Client2 = pb.NewUserServiceClient(Conn2)

// RegisterUser
// @Summary Register user by given particular field.
// @Description Register user with the input payload.
// @Tags User
// @Param user body model.Register true "Create User"
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Failure 500
// @Router /user/register [post]
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.Register
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if user.Address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}
	if user.Phone == "" {
		http.Error(w, "phone is required", http.StatusBadRequest)
		return
	}
	if len(user.Phone) != 10 {
		http.Error(w, "phone number is not valid", http.StatusBadRequest)
		return
	}
	if user.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
		return
	}
	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		http.Error(w, "error: Email is not valid", http.StatusBadRequest)
		return
	}
	if user.DOB == "" {
		http.Error(w, "dob is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	res, err := Client2.RegisterUser(ctx, &pb.CreateRegisterRequest{
		UserRegister: &pb.Register{
			UserName:    user.Name,
			UserAddress: user.Address,
			UserPhone:   user.Phone,
			UserEmail:   user.Email,
			UserDob:     user.DOB,
			ImageId:     user.Image,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(res.UserRegister)
	if err != nil {
		return
	}

}

// GetUser
// @Summary Get details of user by id
// @Description Get details of user by id
// @Tags User
// @Param user_id path string true "get user by id"
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/auth/{user_id} [get]
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]
	if userID == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	res, err := Client2.GetUser(ctx, &pb.GetOneUserRequest{
		UserRegister: &pb.GetOneUser{
			UserId: userID,
		},
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		// Handle encoding error
		fmt.Println("Error encoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// GetUsers
// @Summary Get details of users
// @Description Get details of users
// @Tags User
// @Accept  json
// @Produce  json
// @Success 200
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/auth/ [get]
func GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	res, err := Client2.GetMultipleUsers(ctx, &pb.GetUsersRequest{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(res.Users)
	if err != nil {
		return
	}
}

// LoginUsers
// @Summary login user by given phone and email.
// @Description login user with the input payload.
// @Tags User
// @Param login body model.Login true "login User"
// @Accept json
// @Produce json
// @Success 201
// @Failure 400
// @Failure 500
// @Router /user/login [post]
func LoginUsers(w http.ResponseWriter, r *http.Request) {
	// Create a context
	ctx := context.Background()

	var login model.Login

	err := json.NewDecoder(r.Body).Decode(&login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//validate user mobile number
	if login.UserPhone == "" {
		http.Error(w, "error: user phone is required", http.StatusBadRequest)
		return
	}
	if len(login.UserPhone) != 10 {
		http.Error(w, "phone number is not valid", http.StatusBadRequest)
		return
	}
	//validate user email address
	if login.UserEmail == "" {
		http.Error(w, "error: user email is required", http.StatusBadRequest)
		return
	}
	_, err = mail.ParseAddress(login.UserEmail)
	if err != nil {
		http.Error(w, "error: Email is not valid", http.StatusBadRequest)
		return
	}

	res, err := Client2.LoginUser(ctx, &pb.CreateLoginRequest{
		UserLogin: &pb.Login{
			UserPhone: login.UserPhone,
			UserEmail: login.UserEmail,
		},
	})
	if err != nil {
		http.Error(w, "error: Not login please try again.", http.StatusBadRequest)
		return
	}

	tokenS, err := middleware.GenerateJWTToken(res.UserLogin.UserPhone, res.UserLogin.UserEmail, res.UserLogin.UserId)
	if err != nil {
		http.Error(w, "error: Not login please try again.", http.StatusInternalServerError)
		return
	}
	// Add the token to the response headers
	w.Header().Set("Authorization", "Bearer "+tokenS)
	w.Write([]byte(tokenS))

}

// UpdateUser
// @Summary update user detail by given field.
// @Description update user detail with the given user fields.
// @Tags User
// @Param user body model.Update true "update user"
// @Param user_id path string true "update user by id"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/auth/{user_id} [put]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["user_id"]

	var user model.Update
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}
	if user.Address == "" {
		http.Error(w, "address is required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	res, err := Client2.UpdateUserDetail(ctx, &pb.UpdateUserRequest{
		UpdateUser: &pb.Register{
			UserId:      userID,
			UserName:    user.Name,
			UserAddress: user.Address,
		},
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res.UpdateUser)
	if err != nil {
		return
	}
}

// DeleteUser
// @Summary delete user detail by given id.
// @Description delete user detail with the given user id.
// @Tags User
// @Param user_id path string true "delete user by id"
// @Accept json
// @Produce json
// @Success 204
// @Failure 400
// @Failure 500
// @Security ApiKeyAuth
// @Router /user/auth/{user_id} [delete]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user model.Register

	user.Id = mux.Vars(r)["user_id"]
	if user.Id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	ctx := context.Background()
	res, err := Client2.DeleteUserDetail(ctx, &pb.DeleteUserRequest{
		UserId: user.Id,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(res.Success)
	if err != nil {
		return
	}
}
