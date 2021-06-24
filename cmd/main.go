package main

import (
	"gora/pkg/db"
	"gora/pkg/handlers"
	"gora/pkg/service"
	"log"
	"net/http"
	"time"
)

func main() {
	sqliteDB := db.NewSQLiteDB()
	srvc := service.NewService(sqliteDB)
	handler := handlers.NewHandler(srvc)
	srv := &http.Server{
		Addr:              ":8080",
		Handler:           handler.InitRoutes(),
		MaxHeaderBytes:    1 << 20, // 1 MB
		WriteTimeout:      10 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}
