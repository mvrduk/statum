package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mvrduk/statum/services/user-service/config"
	appuser "github.com/mvrduk/statum/services/user-service/internal/application/service/user"
	"github.com/mvrduk/statum/services/user-service/internal/infrastructure/repository/postgres"
	"github.com/mvrduk/statum/services/user-service/internal/infrastructure/security"
	httptransport "github.com/mvrduk/statum/services/user-service/internal/transport/http"
	"github.com/mvrduk/statum/services/user-service/internal/transport/http/handler"
)

func main() {
	cfg, err := config.Load("config/app.ini")

	if err != nil {
		log.Fatal("config error", err)
	}

	db, err := sql.Open("pgx", cfg.DatabaseURL)

	if err != nil {
		log.Fatal("database open error: ", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("database close error", err)
		}
	}()

	if err := db.Ping(); err != nil {
		log.Fatal("database ping error", err)
	}

	userRepository := postgres.NewUserRepository(db)

	userService := appuser.NewUserService(userRepository, security.NewBcryptHasher())
	userHandler := handler.NewUserHandler(userService)

	server := &http.Server{
		Addr:         cfg.HTTPAddr,
		Handler:      httptransport.New(userHandler),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	log.Println("user service started on", cfg.HTTPAddr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("http server error: ", err)
	}
}
