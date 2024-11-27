package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RohithBN/netflix-api/routes"
)

func main() {
	fmt.Println("NETFLIX_API")
	fmt.Println("SERVER GETTING STARTED")
	r := routes.Router()
	log.Fatal(http.ListenAndServe(":4000", r))
	fmt.Println("Listening on port 4000")

}
