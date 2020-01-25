package main

import (
	"database/sql"
	"fmt"
	_forumDelivery "github.com/Marshality/tech-db/forum/delivery"
	_forumRepo "github.com/Marshality/tech-db/forum/repository"
	_forumUcase "github.com/Marshality/tech-db/forum/usecase"
	"github.com/Marshality/tech-db/middleware"
	_postDelivery "github.com/Marshality/tech-db/post/delivery"
	_postRepo "github.com/Marshality/tech-db/post/repository"
	_postUcase "github.com/Marshality/tech-db/post/usecase"
	_serviceDelivery "github.com/Marshality/tech-db/service/delivery"
	_serviceRepo "github.com/Marshality/tech-db/service/repository"
	_serviceUcase "github.com/Marshality/tech-db/service/usecase"
	_threadDelivery "github.com/Marshality/tech-db/thread/delivery"
	_threadRepo "github.com/Marshality/tech-db/thread/repository"
	_threadUcase "github.com/Marshality/tech-db/thread/usecase"
	"github.com/Marshality/tech-db/tools"
	_userDelivery "github.com/Marshality/tech-db/user/delivery"
	_userRepo "github.com/Marshality/tech-db/user/repository"
	_userUcase "github.com/Marshality/tech-db/user/usecase"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func main() {
	var config *tools.Config
	var connection string
	var err error

	mode := os.Getenv("APP_ENV")

	if mode == "dev" {
		config, err = tools.LoadConfiguration("config.json")
		connection = "host=%s port=%s dbname=%s sslmode=disable"
	} else {
		config, err = tools.LoadConfiguration("docker.json")
		connection = "user=%s password=%s host=%s port=%s dbname=%s sslmode=disable"
	}

	if err != nil {
		log.Fatal(err)
	}

	var dbConn *sql.DB

	if mode == "dev" {
		dbConn, err = sql.Open("postgres",
			fmt.Sprintf(connection,
				config.Database.Host,
				config.Database.Port,
				config.Database.Name),
		)
	} else {
		dbConn, err = sql.Open("postgres",
			fmt.Sprintf(connection,
				config.Database.User,
				config.Database.Password,
				config.Database.Host,
				config.Database.Port,
				config.Database.Name),
		)
	}

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
	threadRepo := _threadRepo.NewThreadRepository(dbConn)
	postRepo := _postRepo.NewPostRepository(dbConn)
	serviceRepo := _serviceRepo.NewServiceRepository(dbConn)

	// Usecases
	userUcase := _userUcase.NewUserUsecase(userRepo)
	forumUcase := _forumUcase.NewForumUsecase(forumRepo, userUcase)
	threadUcase := _threadUcase.NewThreadUsecase(threadRepo, userUcase, forumUcase)
	postUcase := _postUcase.NewPostUsecase(postRepo, threadUcase)
	serviceUcase := _serviceUcase.NewServiceUsecase(serviceRepo)

	// Handlers
	_userDelivery.ConfigureUserHandler(e, userUcase)
	_forumDelivery.ConfigureForumHandler(e, forumUcase, threadUcase)
	_threadDelivery.ConfigureThreadHandler(e, threadUcase, postUcase)
	_postDelivery.ConfigurePostHandler(e, postUcase)
	_serviceDelivery.ConfigureServiceHandler(e, serviceUcase)

	log.Fatal(e.Start(config.Server.Host + ":" + config.Server.Port))
}
