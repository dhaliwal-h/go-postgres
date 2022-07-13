package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dhaliwal-h/go-postgres/router"
)

func main() {
	r := router.Router()

	fmt.Println("Starting Server on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
