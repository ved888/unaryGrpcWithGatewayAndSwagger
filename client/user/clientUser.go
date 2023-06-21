package user

import (
	"context"
	"encoding/json"
	"fmt"
	"gRPC2/jwtUser"
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

// RegisterUser register user with these payload
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.Register
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	user.Name = r.PostFormValue("name")
	user.Address = r.PostFormValue("address")
	user.Phone = r.PostFormValue("phone")
	user.Email = r.PostFormValue("email")
	user.DOB = r.PostFormValue("dob")
	user.Image = r.PostFormValue("image")

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
	if user.Email == "" {
		http.Error(w, "email is required", http.StatusBadRequest)
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
			UserPhoto:   user.Image,
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

// GetUser get user details by user id
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]
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

// GetUsers get all user details
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

// LoginUsers login user by email and phone
func LoginUsers(w http.ResponseWriter, r *http.Request) {
	// Create a context
	ctx := context.Background()
	//defer Conn2.Close()
	userPhone := r.PostFormValue("phone")
	userEmail := r.PostFormValue("email")

	//validate user mobile number
	if userPhone == "" {
		http.Error(w, "error: user phone is required", http.StatusBadRequest)
		return
	}
	if len(userPhone) != 10 {
		http.Error(w, "phone number is not valid", http.StatusBadRequest)
		return
	}
	//validate user email address
	if userEmail == "" {
		http.Error(w, "error: user email is required", http.StatusBadRequest)
		return
	}
	_, err := mail.ParseAddress(userEmail)
	if err != nil {
		http.Error(w, "error: Email is not valid", http.StatusBadRequest)
		return
	}

	res, err := Client2.LoginUser(ctx, &pb.CreateLoginRequest{
		UserLogin: &pb.Login{
			UserPhone: userPhone,
			UserEmail: userEmail,
		},
	})
	if err != nil {
		http.Error(w, "error: Not login please try again.", http.StatusBadRequest)
		return
	}

	tokenS, err := jwtUser.GenerateJWTToken(res.UserLogin.UserPhone, res.UserLogin.UserEmail, res.UserLogin.UserId)
	if err != nil {
		http.Error(w, "error: Not login please try again.", http.StatusInternalServerError)
		return
	}
	// Add the token to the response headers
	w.Header().Set("Authorization", "Bearer "+tokenS)
	w.Write([]byte(tokenS))

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id"]

	var user model.Register
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user.Name = r.PostFormValue("name")
	user.Address = r.PostFormValue("address")

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

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var user model.Register

	user.Id = mux.Vars(r)["id"]
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
