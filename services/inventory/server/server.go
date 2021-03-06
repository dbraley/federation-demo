package main

// This can be run either as a proxy to dgraph, or as a service in it's own right

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gorilla/websocket"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/StevenACoffman/federation-demo/services/inventory"
	"github.com/vektah/gqlparser"
)

const (
	defaultPort   = "4004"
	dgraphAddress = "http://localhost:8080"
)

var proxyDgraph = flag.Bool("proxy-dgraph", false, "Proxy inventory service with Dgraph")

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	flag.Parse()
	if *proxyDgraph {
		url, _ := url.Parse(dgraphAddress)
		http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		http.Handle("/graphql", httputil.NewSingleHostReverseProxy(url))
	} else {
		port = "8080"
		gqlparser.MustLoadSchema()
		srv := handler.New(inventory.NewExecutableSchema(inventory.Config{Resolvers: &inventory.Resolver{}}))
		srv.AddTransport(transport.POST{})
		srv.AddTransport(transport.Websocket{
			KeepAlivePingInterval: 10 * time.Second,
			Upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true
				},
			},
		})
		srv.Use(extension.Introspection{})
		http.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
		http.Handle("/graphql", srv)
	}
	log.Printf("connect to http://localhost:%s/ for inventory GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
