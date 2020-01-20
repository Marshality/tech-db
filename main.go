package main

import (
	"database/sql"
	"fmt"
	"github.com/Marshality/tech-db/middleware"
	"github.com/Marshality/tech-db/tools"
	"github.com/Marshality/tech-db/user/delivery"
	_userRepo "github.com/Marshality/tech-db/user/repository"
	_userUcase "github.com/Marshality/tech-db/user/usecase"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	config, err := tools.LoadConfiguration("config.json")

	if err != nil {
		log.Fatal(err)
	}

	dbConn, err := sql.Open("postgres",
		fmt.Sprintf("host=%s port=%s dbname=%s sslmode=disable",
			config.Database.Host,
			config.Database.Port,
			config.Database.Name),
	)

	if err != nil {
		log.Fatal(err)
	}

	if err := dbConn.Ping(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := dbConn.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	e := echo.New()
	mw := middleware.NewGoMiddleware()
	e.Use(mw.CORS, mw.AccessLog)

	// Repositories
	userRepo := _userRepo.NewUserRepository(dbConn)

	// Usecases
	userUcase := _userUcase.NewUserUsecase(userRepo)

	// Handlers
	delivery.ConfigureUserHandler(e, userUcase)

	log.Fatal(e.Start(config.Server.Host + ":" + config.Server.Port))
}
