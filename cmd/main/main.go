package main

import (
    "fmt"
    "net/http"
    "github.com/gambee/flight_path_tracker/internal/handlers"
)

const port = `:8080`

func main() {
    http.HandleFunc("/calculate", handlers.Calculate)

    fmt.Printf("Listening on %s...\n", port)
    if err := http.ListenAndServe(port, nil); err != nil {
        panic(err)
    }
}
