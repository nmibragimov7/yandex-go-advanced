package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
	"yandex-go-advanced/internal/config"
	"yandex-go-advanced/internal/handlers"
)

type Config struct {
	Files []string `env:"FILES" envSeparator:":"`
	Home  string   `env:"HOME"`
	// required требует, чтобы переменная TASK_DURATION была определена
	TaskDuration time.Duration `env:"TASK_DURATION,required"`
}

func main() {
	//mux := http.NewServeMux()x
	//mux.HandleFunc(`/{id}`, handlers.IDPage)
	//mux.HandleFunc(`/`, handlers.MainPage)

	server := config.Init()
	fmt.Println("server", server)

	log.Fatal(http.ListenAndServe(server, handlers.Router()))
}
