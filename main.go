package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"golang-api-auth-template/data"
	"golang-api-auth-template/http/router"
	"log"
	"net/http"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("ERR: godotenv.Load(): Error loading .env file")
	}

	data.SetDBConfig(&data.DBconfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASS"),
		Name:     os.Getenv("DB_NAME"),
		SSLmode:  os.Getenv("DB_SSL"),
	})

	r := router.NewRouter()

	addr := os.Getenv("HOST")
	port := os.Getenv("PORT")

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%v", addr, port),
		Handler: r,
	}

	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("ERR: ListenAndServe: %s\n", err)
	}
}
