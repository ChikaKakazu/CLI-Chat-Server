package main

import (
	"fmt"
	"net"

	chat "github.com/ChikaKakazu/CLI-Chat-Server/domain/chat"
	"github.com/ChikaKakazu/CLI-Chat-Server/pb"
	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Failed to listen: ", err)
		return
	}

	s := grpc.NewServer()
	pb.RegisterChatRoomServiceServer(s, &chat.Server{Rooms: make(map[string][]chan *pb.ChatMessage)})
	if err := s.Serve(lis); err != nil {
		fmt.Println("Failed to serve: ", err)
		return
	}

}
