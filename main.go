package main

import (
	"github.com/xiwen1/mini-tiktok-comment/Comment/idl/comment"
	"github.com/xiwen1/mini-tiktok-comment/Comment/service"
	"google.golang.org/grpc"
	"log"
	"net"
)

var port = "50051"

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("fail to listen: %v", err)
	}
	err = Comment.InitComment()
	if err != nil {
		log.Println("unable to init comment service")
	}
	s := grpc.NewServer()
	comment.RegisterCommentActionServer(s, &Comment.CommentActionServer{})
	log.Printf("comment server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err.Error())
	}
}
