package cmd

import (
	"auth/internal/auth"
	"auth/internal/global"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"log"
	"net"
	authEntity "shared/entity/auth"
	"shared/interfaces"
	"shared/rpc/pb"

	"google.golang.org/grpc"
)

var ()

type server struct {
	pb.UnimplementedAuthServer
	authService interfaces.AuthService
}

func StartGRPCServer() {
	grpcPort := global.Config.Auth.Grpc.Port
	lis, err := net.Listen("tcp", ":"+fmt.Sprint(grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterAuthServer(s, server{
		authService: auth.NewService(),
	})
	log.Printf("GRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s server) GetSession(_ context.Context, req *pb.SessionReq) (*pb.AuthResponse, error) {
	session, err := s.authService.GetSession(&authEntity.SessionRequest{
		AccessToken: req.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	var user []byte
	user, err = bson.Marshal(session.User)

	return &pb.AuthResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		UserId:       session.UserId.String(),
		User:         user,
	}, nil
}

//func (s server) GetProfile(_ context.Context, req *pb.)
