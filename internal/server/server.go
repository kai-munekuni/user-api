package server

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/kai-munekuni/user-api/internal/config"
	"github.com/kai-munekuni/user-api/internal/db"
	"github.com/kai-munekuni/user-api/internal/http"
)

// Run bootstrap of server
func Run() {
	os.Exit(run(context.Background()))
}

func run(ctx context.Context) int {
	termCh := make(chan os.Signal, 1)
	signal.Notify(termCh, syscall.SIGTERM, syscall.SIGINT)

	sa := option.WithCredentialsJSON(config.GCPSAKey())
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Printf("%+v", err)
		return 1
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Printf("%+v", err)
		return 1
	}

	defer client.Close()

	d := db.NewFirestore(client)
	s := http.NewServer(8080, d)
	errCh := make(chan error, 1)

	go func() {
		errCh <- s.Start()
	}()

	select {
	case <-termCh:
		if err := s.Stop(ctx); err != nil {
			log.Printf("%+v", err)
		}
		return 0
	case <-errCh:
		if err := s.Stop(ctx); err != nil {
			log.Printf("%+v", err)
		}
		return 1
	}
}
