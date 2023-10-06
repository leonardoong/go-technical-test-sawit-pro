package main

import (
	"log"
	"os"
	"time"

	"github.com/SawitProRecruitment/UserService/config"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/handler"
	"github.com/SawitProRecruitment/UserService/repository"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	var server generated.ServerInterface = newServer()

	generated.RegisterHandlers(e, server)
	e.Logger.Fatal(e.Start(":1323"))
}

func newServer() *handler.Server {
	cfg := initConfig()

	dbDsn := os.Getenv("DATABASE_URL")
	var repo repository.RepositoryInterface = repository.NewRepository(repository.NewRepositoryOptions{
		Dsn: dbDsn,
	})
	opts := handler.NewServerOptions{
		Repository: repo,
		Config:     cfg,
	}
	return handler.NewServer(opts)
}

func initConfig() *config.Config {
	prvKey, err := os.ReadFile("cert/id_rsa")
	if err != nil {
		time.Sleep(time.Hour)
		log.Fatalln(err)
	}

	pubKey, err := os.ReadFile("cert/id_rsa.pub")
	if err != nil {
		log.Fatalln(err)
	}

	jwtToken := config.NewJWT(prvKey, pubKey)

	return &config.Config{
		JWT: jwtToken,
	}
}
