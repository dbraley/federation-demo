package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/handler"
	"github.com/marwan-at-work/federation-demo/services/products"
	"github.com/vektah/gqlparser"
)

const defaultPort = "4003"

func main() {
	gqlparser.MustLoadSchema()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", handler.Playground("GraphQL playground", "/graphql"))
	http.Handle("/graphql", handler.GraphQL(products.NewExecutableSchema(products.Config{Resolvers: &products.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
