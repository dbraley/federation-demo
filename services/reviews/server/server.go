package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/marwan-at-work/federation-demo/services/reviews"
	"github.com/vektah/gqlparser"
)

const defaultPort = "4002"

func main() {
	gqlparser.MustLoadSchema()
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", handler.New(reviews.NewExecutableSchema(reviews.Config{Resolvers: &reviews.Resolver{}})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
