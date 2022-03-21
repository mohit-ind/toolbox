package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	logger "github.com/toolboxlogger"
	middlewares "github.com/toolboxmiddlewares"
)

func main() {
	r := mux.NewRouter()

	log := logger.NewCommonLogger("Test Application", "v1.3.2", "test", "localhost", false)

	r.Use(middlewares.LoggingMiddleware(log))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from Test Application\n"))
	})

	fmt.Println("Test Application listening at http://localhost:8080")

	http.ListenAndServe(":8080", r)
}
