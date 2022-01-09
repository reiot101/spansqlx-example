package grpc

import (
	"context"
	"log"
	"net"
	"strconv"

	"github.com/reiot777/spansqlx"
	"github.com/reiot777/spansqlx-example/grpc/packet"
	"github.com/reiot777/spansqlx-example/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	Host string `long:"host" description:"The IP to listen on" default:"0.0.0.0" env:"HOST"`
	Port int    `long:"port" description:"The port to listen on for insecure connections, defaults to a random value" default:"8888" env:"PORT"`

	SpannerDatabase string `long:"spanner-database" description:"The database for spanner" default:"projects/sandbox/instances/sandbox/databases/sandbox" env:"SPANNER_DATABASE"`
}

func (s *Server) Serve(ctx context.Context) {
	// Open the spansqlx.DB instance
	db, err := spansqlx.Open(ctx, spansqlx.WithDatabase(s.SpannerDatabase))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// grpc service
	svc := Service{
		Store: store.New(db),
	}

	// grpc server
	srv := grpc.NewServer()

	// register service
	packet.RegisterAccountServiceServer(srv, &svc)
	packet.RegisterTodoServiceServer(srv, &svc)
	reflection.Register(srv)

	// listener
	lis, err := net.Listen("tcp", net.JoinHostPort(s.Host, strconv.Itoa(s.Port)))
	if err != nil {
		log.Fatal(err)
	}
	defer lis.Close()

	log.Println("gRPC server serving at", lis.Addr())

	// grpc server serving
	if err := srv.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
