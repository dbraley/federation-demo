package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/99designs/gqlgen/graphql/playground"
)

const defaultPort = "4004"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	url, _ := url.Parse("http://localhost:8080")
	http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	http.Handle("/graphql", httputil.NewSingleHostReverseProxy(url))
	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
