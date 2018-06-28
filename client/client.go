package main

import (
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"google.golang.org/grpc/testdata"
	"google.golang.org/grpc/credentials"
	pb "github.com/Jepp2078/myprotos/account"
	UUID "github.com/google/uuid"
)

var (
	tls                = flag.Bool("tls", true, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "", "The file containning the CA root cert file")
	serverAddr         = flag.String("server_addr", "127.0.0.1:10000", "The server address in the format of host:port")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

// getAccount gets the account for UUID
func getAccount(client pb.AccountsClient, accountId *pb.AccountID) {
	log.Printf("Getting account for id (%d)", accountId.Id)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	account, err := client.GetAccount(ctx, accountId)
	if err != nil {
		log.Fatalf("%v.GetAccount(_) = _, %v: ", client, err)
	}
	log.Println(account)
}

func main() {
	flag.Parse()
	var opts []grpc.DialOption
	if *tls {
		if *caFile == "" {
			*caFile = testdata.Path("ca.pem")
		}
		creds, err := credentials.NewClientTLSFromFile(*caFile, *serverHostOverride)
		if err != nil {
			log.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewAccountsClient(conn)

	accountId := &pb.AccountID{Id:UUID.New().String()}
	getAccount(client, accountId)
}