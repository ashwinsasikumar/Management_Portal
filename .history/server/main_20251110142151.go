package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Backend is running!")
    })
    fmt.Println("Server started at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}