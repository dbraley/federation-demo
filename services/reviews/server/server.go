package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/marwan-at-work/federation-demo/services/reviews"
)

const defaultPort = "4002"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/graphql"))
	http.Handle("/graphql", handler.GraphQL(reviews.NewExecutableSchema(reviews.Config{Resolvers: &reviews.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
