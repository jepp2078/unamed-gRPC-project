package main

import (
	"golang.org/x/net/context"
	pb "github.com/Jepp2078/myprotos/account"
	UUID "github.com/google/uuid"
	time2 "time"
	"flag"
	"net"
	"fmt"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/testdata"
	"google.golang.org/grpc/credentials"
)

var (
	tls        = flag.Bool("tls", true, "Connection uses TLS if true, else plain TCP")
	certFile   = flag.String("cert_file", "", "The TLS cert file")
	keyFile    = flag.String("key_file", "", "The TLS key file")
	port       = flag.Int("port", 10000, "The server port")
)

type accountsServer struct {

}

func (s *accountsServer) GetAccount(ctx context.Context, AccountID *pb.AccountID) (*pb.Account, error) {
	return &pb.Account{Id: AccountID.Id, Email: "test@test.com", Username: "test", Password:"test"}, nil
}

func (s *accountsServer) GetAuthObject(ctx context.Context, AccountID *pb.AccountID) (*pb.AuthObject, error) {
	now := time2.Now().Unix()
	valid := time2.Now().Add(time2.Hour * 8).Unix()
	uuid := UUID.New().String()
	return &pb.AuthObject{Token:uuid, IssuedAt:now, ValidUntill:valid}, nil
}

func newServer() *accountsServer {
	s := &accountsServer{}
	return s
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	if *tls {
		if *certFile == "" {
			*certFile = testdata.Path("server1.pem")
		}
		if *keyFile == "" {
			*keyFile = testdata.Path("server1.key")
		}
		creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
		if err != nil {
			log.Fatalf("Failed to generate credentials %v", err)
		}
		opts = []grpc.ServerOption{grpc.Creds(creds)}
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAccountsServer(grpcServer, newServer())
	log.Printf("Server listning on: %s", lis.Addr().String())
	grpcServer.Serve(lis)
}