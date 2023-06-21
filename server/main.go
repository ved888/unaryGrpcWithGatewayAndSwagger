package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"userGrpc/database"

	"github.com/google/uuid"
	"log"
	"userGrpc/dbhelper"
	"userGrpc/model"
	"userGrpc/pb/pb"
)

var (
	port = flag.Int("port", 50052, "gRPC server port")
)

type server struct {
	DAO dbhelper.DAO
	pb.UnimplementedUserServiceServer
}

func main() {
	flag.Parse()

	db, err := database.DbConnection()
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	defer db.Close()

	fmt.Println("gRPC server running ...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	pb.RegisterUserServiceServer(s, &server{DAO: dbhelper.DAO{
		DB: db,
	},
	})

	log.Printf("Server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) RegisterUser(ctx context.Context, req *pb.CreateRegisterRequest) (*pb.CreateRegisterResponse, error) {
	fmt.Println("Create User")
	user := req.GetUserRegister()
	user.UserId = uuid.New().String()

	data := model.Register{
		Id:      user.UserId,
		Name:    user.UserName,
		Address: user.UserAddress,
		Phone:   user.UserPhone,
		Email:   user.UserEmail,
		DOB:     user.UserDob,
		Image:   user.UserPhoto,
	}

	userId, err := s.DAO.CreateUser(&data)
	if err != nil {
		return nil, err
	}

	return &pb.CreateRegisterResponse{
		UserRegister: &pb.Register{
			UserId:      *userId,
			UserName:    user.GetUserName(),
			UserAddress: user.GetUserAddress(),
			UserPhone:   user.GetUserPhone(),
			UserEmail:   user.GetUserEmail(),
			UserDob:     user.GetUserDob(),
			UserPhoto:   user.GetUserPhoto(),
		},
	}, nil
}
func (s *server) GetUser(ctx context.Context, req *pb.GetOneUserRequest) (*pb.GetOneUserResponse, error) {
	user, err := s.DAO.GetUserById(req.GetUserRegister().GetUserId())
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, fmt.Errorf("user not found with ID: %s", req.GetUserRegister().GetUserId())
	}

	return &pb.GetOneUserResponse{
		UserRegister: &pb.Register{
			UserId:      user.Id,
			UserName:    user.Name,
			UserAddress: user.Address,
			UserPhone:   user.Phone,
			UserEmail:   user.Email,
			UserDob:     user.DOB,
		},
	}, nil
}

func (s *server) GetMultipleUsers(ctx context.Context, req *pb.GetUsersRequest) (*pb.GetUsersResponse, error) {
	users, err := s.DAO.GetUsers()
	if err != nil {
		return nil, err
	}
	pbUsers := make([]*pb.Register, len(users))
	for i, user := range users {
		pbUsers[i] = &pb.Register{
			UserId:      user.Id,
			UserName:    user.Name,
			UserAddress: user.Address,
			UserPhone:   user.Phone,
			UserEmail:   user.Email,
			UserDob:     user.DOB,
		}
	}
	return &pb.GetUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *server) LoginUser(ctx context.Context, req *pb.CreateLoginRequest) (*pb.CreateLoginResponse, error) {
	fmt.Println("login user...")
	LoginUser := req.GetUserLogin()
	data := model.Register{
		Phone: LoginUser.GetUserPhone(),
		Email: LoginUser.GetUserEmail(),
	}
	userdata, err := s.DAO.GetUserByEmailAndPhone(data.Phone, data.Email)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if userdata.Id == "" {
		return nil, errors.New("user is not exist with this phone number.please register first... ")
	}
	return &pb.CreateLoginResponse{
		UserLogin: &pb.Register{
			UserId:      userdata.Id,
			UserName:    userdata.Name,
			UserAddress: userdata.Address,
			UserPhone:   userdata.Phone,
			UserEmail:   userdata.Email,
			UserDob:     userdata.DOB,
		},
	}, nil
}

func (s *server) UpdateUserDetail(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	user := req.GetUpdateUser()
	data := model.Register{
		Name:    user.UserName,
		Address: user.UserAddress,
	}
	err := s.DAO.UpdateUser(data, user.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{
		UpdateUser: &pb.Register{
			UserName:    user.GetUserName(),
			UserAddress: user.GetUserAddress(),
		},
	}, nil

}

func (s *server) DeleteUserDetail(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.DAO.DeleteUser(req.GetUserId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{
			Success: true,
		},
		nil
}
