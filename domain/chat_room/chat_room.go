package chatroom

import (
	"context"
	"fmt"

	"github.com/ChikaKakazu/CLI-Chat-Server/pb"
)

type Server struct {
	pb.UnimplementedChatRoomServiceServer
	Rooms map[string]bool
}

func (s *Server) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	fmt.Println("Received CreateRoom request: ", req.RoomName)
	if _, exists := s.Rooms[req.RoomName]; exists {
		return &pb.CreateRoomResponse{
			Success: false,
			Message: "Room already exists",
		}, nil
	}

	s.Rooms[req.RoomName] = true
	return &pb.CreateRoomResponse{
		Success: true,
		Message: "Room created",
	}, nil
}

func (s *Server) JoinRoom(ctx context.Context, req *pb.JoinRoomRequest) (*pb.JoinRoomResponse, error) {
	fmt.Println("Received JoinRoom request: ", req.RoomName)
	if _, exists := s.Rooms[req.RoomName]; !exists {
		return &pb.JoinRoomResponse{
			Success: false,
			Message: "Room does not exist",
		}, nil
	}

	return &pb.JoinRoomResponse{
		Success: true,
		Message: "Joined room",
	}, nil
}
