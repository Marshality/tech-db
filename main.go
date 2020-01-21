package main

import (
	"database/sql"
	"fmt"
	_forumDelivery "github.com/Marshality/tech-db/forum/delivery"
	_forumRepo "github.com/Marshality/tech-db/forum/repository"
	_forumUcase "github.com/Marshality/tech-db/forum/usecase"
	"github.com/Marshality/tech-db/middleware"
	"github.com/Marshality/tech-db/tools"
	_userDelivery "github.com/Marshality/tech-db/user/delivery"
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
	forumRepo := _forumRepo.NewForumRepository(dbConn)

	// Usecases
	userUcase := _userUcase.NewUserUsecase(userRepo)
	forumUcase := _forumUcase.NewForumUsecase(forumRepo, userUcase)

	// Handlers
	_userDelivery.ConfigureUserHandler(e, userUcase)
	_forumDelivery.ConfigureForumHandler(e, forumUcase)

	log.Fatal(e.Start(config.Server.Host + ":" + config.Server.Port))
}
