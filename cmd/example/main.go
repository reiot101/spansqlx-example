package main

import (
	"context"
	"os"

	"github.com/jessevdk/go-flags"
	"github.com/reiot777/spansqlx-example/grpc"
)

func main() {
	var srv grpc.Server

	parser := flags.NewParser(&srv, flags.Default)

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	srv.Serve(context.Background())
}
