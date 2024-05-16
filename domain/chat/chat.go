package chat

import (
	"context"
	"fmt"
	"io"
	"sync"

	"github.com/ChikaKakazu/CLI-Chat-Server/pb"
)

type Server struct {
	pb.UnimplementedChatRoomServiceServer
	Rooms   map[string][]chan *pb.ChatMessage
	History map[string][]*pb.ChatMessage
	mu      sync.Mutex
}

func (s *Server) CreateRoom(ctx context.Context, req *pb.CreateRoomRequest) (*pb.CreateRoomResponse, error) {
	fmt.Println("Received CreateRoom request: ", req.RoomName)
	if _, exists := s.Rooms[req.RoomName]; exists {
		return &pb.CreateRoomResponse{
			Success: false,
			Message: "Room already exists",
		}, nil
	}

	s.Rooms[req.RoomName] = make([]chan *pb.ChatMessage, 0)
	s.History[req.RoomName] = make([]*pb.ChatMessage, 0)
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

func (s *Server) Chat(steam pb.ChatRoomService_ChatServer) error {
	fmt.Println("Received Chat request")
	msg, err := steam.Recv()
	if err != nil {
		return err
	}

	roomName := msg.GetRoomName()
	s.mu.Lock()
	// 部屋が存在しない場合はエラーを返す
	if _, exists := s.Rooms[roomName]; !exists {
		s.mu.Unlock()
		fmt.Printf("Room does not exist")
		return fmt.Errorf("room does not exist")
	}

	ch := make(chan *pb.ChatMessage, 10)
	s.Rooms[roomName] = append(s.Rooms[roomName], ch)
	history := s.History[roomName]
	s.mu.Unlock()

	for _, m := range history {
		if err := steam.Send(m); err != nil {
			fmt.Println("Failed to send a message: ", err)
			return err
		}
	}

	go func() {
		for {
			in, err := steam.Recv()
			if err == io.EOF {
				return
			}
			if err != nil {
				fmt.Println("Failed to receive a message: ", err)
				return
			}

			s.mu.Lock()
			s.History[roomName] = append(s.History[roomName], in)
			for _, c := range s.Rooms[roomName] {
				c <- in
			}
			s.mu.Unlock()
		}
	}()

	for {
		select {
		case out := <-ch:
			if err := steam.Send(out); err != nil {
				fmt.Println("Failed to send a message: ", err)
				return err
			}
		}
	}
}

func (s *Server) ListRooms(ctx context.Context, req *pb.ListRoomsRequest) (*pb.ListRoomsResponse, error) {
	roomNames := make([]string, 0, len(s.Rooms))
	for roomName := range s.Rooms {
		roomNames = append(roomNames, roomName)
	}
	return &pb.ListRoomsResponse{
		RoomNames: roomNames,
	}, nil
}
