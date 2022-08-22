package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/zufzuf/cake-store/server"
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	defer cancel()

	if err := server.NewHTTPServer().Run(ctx); err != nil {
		log.Fatalf("failed starting server, err : \n%+v", err)
	}
}
