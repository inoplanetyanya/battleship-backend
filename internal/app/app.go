package app

import (
	"battleship/pkg/handler"
	"battleship/pkg/repository"
	"battleship/pkg/service"
	wsserver "battleship/pkg/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
	"golang.org/x/net/websocket"
)

func Run() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Ошибка загрузки .env файла: %s", err.Error())
	}

	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")
	dbuser := os.Getenv("DB_USER")
	dbpassword := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	dbsslmode := os.Getenv("DB_SSLMODE")

	config := repository.Config{
		Host:     dbhost,
		Port:     dbport,
		User:     dbuser,
		Password: dbpassword,
		DBname:   dbname,
		SSLmode:  dbsslmode,
	}

	db, err := repository.NewPostgresDB(config)

	if err != nil {
		log.Fatalf("failed to initalize db: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	portws := os.Getenv("PORT_WS")
	serverws := wsserver.NewServer(1024)
	http.Handle("/ws", websocket.Handler(serverws.HandleWS))

	porthttp := os.Getenv("PORT_HTTP")
	serverhttp := &http.Server{
		Addr:    fmt.Sprintf("192.168.0.69:%s", porthttp),
		Handler: handlers.InitRoutes(),
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		log.Printf("WebSocket-сервер запущен на ws://192.168.0.69:%s", portws)
		if err := http.ListenAndServe(fmt.Sprintf("192.168.0.69:%s", portws), nil); err != nil {
			log.Fatalf("Ошибка запуска WebSocket-сервера: %s", err.Error())
		}
	}()

	go func() {
		defer wg.Done()
		log.Printf("HTTP-сервер запущен на http://192.168.0.69:%s", porthttp)
		if err := serverhttp.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска HTTP-сервера: %s", err.Error())
		}
	}()

	wg.Wait()
}
