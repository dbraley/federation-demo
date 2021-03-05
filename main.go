package main

// This is a simple multiple runner for getting all four graphql services to run in the background

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	startDgraph = flag.Bool("start-dgraph", false, "Start inventory service with Dgraph")
	stopDgraph  = flag.Bool("stop-dgraph", false, "Start inventory service with Dgraph")
)

func main() {
	flag.Parse()
	if *startDgraph {
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
	}
	if *stopDgraph {
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
	}

	var subpath = []string{
		"/services/accounts/server",
		"/services/inventory/server",
		"/services/products/server",
		"/services/reviews/server",
	}

	var stdoutBuf = []bytes.Buffer{{}, {}, {}, {}}
	var stderrBuf = []bytes.Buffer{{}, {}, {}, {}}

	ctx, done := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)
	// listen - a simple function that listens to the signals channel for interruption signals and then call done() of the errgroup.
	listen := func() error {
		signalChannel := getStopSignalsChannel()
		select {
		case sig := <-signalChannel:
			log.Printf("Received signal: %s\n", sig)
			done()
		case <-ctx.Done():
			log.Printf("closing signal goroutine\n")
			return ctx.Err()
		}

		return nil
	}
	// listen for os interrupts, and cancel context
	g.Go(listen)
	// run all four federated services
	for i := 0; i < 4; i++ {
		g.Go(runner(subpath[i], &stdoutBuf[i], &stderrBuf[i]))
	}
	var gatewayStdOut, gatewayStdErr bytes.Buffer
	g.Go(gatewayRunner(&gatewayStdOut, &gatewayStdErr))

	err := g.Wait()
	if err != nil {
		log.Fatalf("os.Wait() failed with %v\n", err)
	}
	for i := 0; i < 4; i++ {
		outStr, errStr := string(stdoutBuf[i].Bytes()), string(stderrBuf[i].Bytes())
		fmt.Printf("\n---------\nsubpath %s output:\n\nout:\n%s\nerr:\n%s\n", subpath[i], outStr, errStr)
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
			return fmt.Errorf("cmd.Start() failed with %v\n", err)
		}
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("cmd.Wait() failed with %v\n", err)
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
			return fmt.Errorf("cmd.Start() failed with %v\n", err)
		}
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("cmd.Wait() failed with %v\n", err)
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
