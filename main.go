package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"mini-wallet/libraries/config"
	"mini-wallet/libraries/database"
	"mini-wallet/routing"
	"mini-wallet/schema"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

func main() {
	if _, ok := os.LookupEnv("APP_ENV"); !ok {
		config.Setup(".env")
	}

	// =========================================================================
	// Logging
	log := log.New(os.Stdout, "mini-wallet : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	if err := run(log); err != nil {
		log.Printf("error: shutting down: %s", err)
		os.Exit(1)
	}
}

func run(log *log.Logger) error {
	// =========================================================================
	// App Starting

	log.Printf("main : Started")
	defer log.Println("main : Completed")

	// =========================================================================

	// Start Database

	db, err := database.Open()
	if err != nil {
		return fmt.Errorf("connecting to db: %v", err)
	}
	defer db.Close()

	// Handle cli command
	flag.Parse()

	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return fmt.Errorf("applying migrations: %v", err)
		}
		log.Println("Migrations complete")
		return nil

	case "seed":
		if err := schema.Seed(db); err != nil {
			return fmt.Errorf("seeding database: %v", err)
		}
		log.Println("Seed data complete")
		return nil
	}

	// parameter server
	server := http.Server{
		Addr:         ":" + os.Getenv("APP_PORT"),
		Handler:      routing.API(db, log),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	serverErrors := make(chan error, 1)
	// mulai listening server
	go func() {
		log.Println("server listening on", server.Addr)
		serverErrors <- server.ListenAndServe()
	}()

	// Membuat channel untuk mendengarkan sinyal interupsi/terminate dari OS.
	// Menggunakan channel buffered karena paket signal membutuhkannya.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Mengontrol penerimaan data dari channel,
	// jika ada error saat listenAndServe server maupun ada sinyal shutdown yang diterima
	select {
	case err := <-serverErrors:
		return fmt.Errorf("Starting server: %v", err)

	case <-shutdown:
		log.Println("caught signal, shutting down")

		// Jika ada shutdown, meminta tambahan waktu 5 detik untuk menyelesaikan proses yang sedang berjalan.
		const timeout = 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", timeout, err)
			if err := server.Close(); err != nil {
				return fmt.Errorf("could not stop server gracefully: %v", err)
			}
		}
	}

	return nil
}
