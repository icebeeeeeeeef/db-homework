package grpcx

import (
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Addr string
}

func NewServer(server *grpc.Server, addr string) *Server {
	return &Server{Server: server, Addr: addr}
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", s.Addr)
	if err != nil {
		panic(err)
	}
	err = s.Server.Serve(l)
	return err
}
