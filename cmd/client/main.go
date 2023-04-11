package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	api "google.golang.org/grpc/examples/calculator/api/proto"
	"log"
	"os"
	"strconv"
	"time"
)

var (
	addr      = flag.String("addr", "localhost:50051", "the address to connect to")
	firstArg  = flag.String("x", "2", "first arg")
	secondArg = flag.String("y", "3", "second arg")

	serverCaCert = flag.String("server-ca-cert", "tls/ca-cert.pem", "Client CA cert path")
	clientCert   = flag.String("client-cert", "tls/client-cert.pem", "Client cert path")
	clientKey    = flag.String("client-key", "tls/client-key.pem", "Client key path")
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pemServerCA, err := os.ReadFile(*serverCaCert)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")

	}

	clientCert, err := tls.LoadX509KeyPair(*clientCert, *clientKey)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{clientCert},
		RootCAs:      certPool,
	}
	return credentials.NewTLS(config), nil
}
func main() {
	flag.Parse()

	x, err := strconv.Atoi(*firstArg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	y, err := strconv.Atoi(*secondArg)
	if err != nil {
		log.Fatalf(err.Error())
	}

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(tlsCredentials))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewCalculatorClient(conn)

	for true {
		// Contact the server and print out its response.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		//defer cancel()

		r, err := c.Add(ctx, &api.AddRequest{
			X: int32(x),
			Y: int32(y),
		})
		if err != nil {
			log.Fatalf("could not calculate sum: %v", err)
		}
		log.Printf("sum: %v", r.GetS())
		time.Sleep(time.Second * 5)
		cancel()
	}
}
