package main

// This is a simple multiple runner for getting all four graphql services
// to run in the background.
// It merges all the output to the console

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	startDgraph = flag.Bool("start-dgraph", false, "Start inventory service with Dgraph")
	stopDgraph  = flag.Bool("stop-dgraph", false, "Start inventory service with Dgraph")
	initDgraph  = flag.Bool("init-dgraph", false, "Initialize inventory service with Dgraph")
	useDgraph   = flag.Bool("use-dgraph", false, "Use Dgraph for inventory service ")
)

func main() {
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	flag.Parse()
	switch {
	case *startDgraph:
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("os.Getwd() failed with %v\n", err)
		}
		cmd := exec.Command("docker-compose", "up", "--detach")
		cmd.Dir = wd + "/services/inventory/server"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Starting Dgraph failed with: %v\n", err)
		}
		return
	case *stopDgraph:
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("os.Getwd() failed with %v\n", err)
		}
		cmd := exec.Command("docker-compose", "down")
		cmd.Dir = wd + "/services/inventory/server"
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		if err != nil {
			log.Fatalf("Starting Dgraph failed with: %v\n", err)
		}
		return
	case *initDgraph:
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("os.Getwd() failed with %v\n", err)
		}

		log.Println("curl http://localhost:8080/admin/schema " +
			"--upload-file ./services/inventory/schema.graphql")
		err = putInventorySchema(wd)
		if err != nil {
			log.Fatalf("putInventorySchema failed with %+v\n", err)
		}
		inventoryDataMutation := `{
  "query": "mutation { addProduct(upsert: true, input: [{upc: \"1\", inStock: true}, {upc: \"2\", inStock: false}, {upc: \"3\", inStock: true}]) { product { upc inStock } }}"
}`

		log.Printf("curl --request POST \\\n"+
			"--url http://localhost:8080/graphql \\\n"+
			"--header 'Content-Type: application/json' \\\n"+
			"--data '%s'\n", inventoryDataMutation,
		)

		err = postInventoryData(inventoryDataMutation)
		if err != nil {
			log.Fatalf("postInventoryData failed with %+v\n", err)
		}
		return
	}

	subPaths := []string{
		"/services/accounts/server",
		"/services/products/server",
		"/services/reviews/server",
	}

	if !*useDgraph {
		subPaths = append([]string{"/services/inventory/server"}, subPaths...)
	}

	stdoutBuf := []bytes.Buffer{{}, {}, {}, {}}
	stderrBuf := []bytes.Buffer{{}, {}, {}, {}}

	ctx, done := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	// listen - a simple function that listens to the signals channel for interruption signals and then call done() of
	// the errgroup.
	listen := func() error {
		signalChannel := getStopSignalsChannel()
		select {
		case sig := <-signalChannel:
			log.Printf("received signal: %s\n", sig)
			done()
		case <-ctx.Done():
			log.Printf("closing signal goroutine\n")
			return ctx.Err()
		}
		return nil
	}
	// listen for os interrupts, and cancel context
	g.Go(listen)
	// run all federated services
	for i, subpath := range subPaths {
		g.Go(runner(subpath, &stdoutBuf[i], &stderrBuf[i]))
	}
	var gatewayStdOut, gatewayStdErr bytes.Buffer
	g.Go(gatewayRunner(&gatewayStdOut, &gatewayStdErr))

	err := g.Wait()
	if err != nil {
		log.Fatalf("os.Wait() failed with %v\n", err)
	}
	for i := 0; i < 4; i++ {
		outStr, errStr := stdoutBuf[i].String(), stderrBuf[i].String()
		fmt.Printf("\n---------\nsubpath %s output:\n\nout:\n%s\nerr:\n%s\n", subPaths[i], outStr, errStr)
	}
}

func runner(subPath string, stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer) func() error {
	return func() error {
		cmd := exec.Command("go", "run", "server.go")
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("os.Getwd() failed with %v\n", err)
		}
		cmd.Dir = wd + subPath

		cmd.Stdout = io.MultiWriter(os.Stdout, stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, stderrBuf)

		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("cmd.Start() failed with %v", err)
		}
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("cmd.Wait() failed with %v", err)
		}
		return nil
	}
}

func gatewayRunner(stdoutBuf *bytes.Buffer, stderrBuf *bytes.Buffer) func() error {
	return func() error {
		cmd := exec.Command("node", "gateway.js")
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalf("os.Getwd() failed with %v\n", err)
		}
		cmd.Dir = wd

		cmd.Stdout = io.MultiWriter(os.Stdout, stdoutBuf)
		cmd.Stderr = io.MultiWriter(os.Stderr, stderrBuf)

		// wait an arbitrary amount of time until the federated services are available
		time.Sleep(time.Second * 2)
		err = cmd.Start()
		if err != nil {
			return fmt.Errorf("cmd.Start() failed with %v", err)
		}
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("cmd.Wait() failed with %v", err)
		}
		return nil
	}
}

func getStopSignalsChannel() chan os.Signal {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel,
		os.Interrupt,    // interrupt is syscall.SIGINT, Ctrl+C
		syscall.SIGQUIT, // Ctrl-\
		syscall.SIGHUP,  // "terminal is disconnected"
		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
	)
	return signalChannel
}

func putInventorySchema(wd string) error {
	f, fileErr := os.Open(wd + "/services/inventory/schema.graphql")
	if fileErr != nil {
		return fmt.Errorf("os.Open() failed for inventory/schema.graphql with %w", fileErr)
	}
	defer f.Close()
	req, reqErr := http.NewRequest("PUT", "http://localhost:8080/admin/schema", f)
	if reqErr != nil {
		return fmt.Errorf("unable to create dgraph inventory schema request %w", fileErr)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to PUT dgraph inventory schema %w", fileErr)
	}
	resp.Body.Close()
	return nil
}

func postInventoryData(inventoryDataMutation string) error {
	body := strings.NewReader(inventoryDataMutation)
	req, reqErr := http.NewRequest("POST", "http://localhost:8080/graphql", body)
	if reqErr != nil {
		return fmt.Errorf("unable to make dgraph inventory data request %w", reqErr)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("unable to POST dgraph inventory data %w", reqErr)
	}
	resp.Body.Close()
	return nil
}
