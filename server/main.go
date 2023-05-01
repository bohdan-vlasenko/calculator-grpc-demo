package main

import (
	api "calculator-grpc/api/proto"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"os"
)

var (
	port         = flag.Int("port", 50051, "The server port")
	clientCaCert = flag.String("client-ca-cert", "tls/ca-cert.pem", "Client CA cert path")
	serverCert   = flag.String("servera-cert", "tls/server-cert.pem", "Server cert path")
	serverKey    = flag.String("servera-key", "tls/server-key.pem", "Server key path")
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	pamClientCA, err := os.ReadFile(*clientCaCert)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pamClientCA) {
		return nil, fmt.Errorf("failed to add client  CA's certificate")
	}

	serverCert, err := tls.LoadX509KeyPair(*serverCert, *serverKey)
	if err != nil {
		return nil, err
	}
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return credentials.NewTLS(config), nil
}
func main() {
	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatalf(err.Error())
	}

	s := grpc.NewServer(grpc.Creds(tlsCredentials))

	calculatorSrv := &CalculatorGrpcServer{}
	api.RegisterCalculatorServer(s, calculatorSrv)

	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//INVESTIGATE
// 1. create cluster and install cert manager addon
// 2. create self signed issuer

// REPEATED
// 3. create certificates for server (specify service dns name)
// 4. create certificate for client app (??? investigate DNS service name in case of MTLS)
// 5. create deployments for server and client apps
// 6. mount certs as volumes

// 7. test connection

// INVESTIGATE
// 8* test certificates rotation

// update all grpc clients and servers
