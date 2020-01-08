package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hashicorp/http-echo/version"
)

var (
	listenFlag  = flag.String("listen", ":5678", "address and port to listen")
	textFlag    = flag.String("text", "", "text to put on the webpage")
	versionFlag = flag.Bool("version", false, "display version information")
	lag         = flag.Int("lag", 5, "response lag")
	enableLag   = flag.Bool("enableLag", true, "enable response lag")
	statusCode  = flag.Int("statusCode", 200, "the status code to return")

	// stdoutW and stderrW are for overriding in test.
	stdoutW = os.Stdout
	stderrW = os.Stderr
)

func main() {
	flag.Parse()

	// Asking for the version?
	if *versionFlag {
		fmt.Fprintln(stderrW, version.HumanVersion)
		os.Exit(0)
	}

	// Validation
	if *textFlag == "" {
		fmt.Fprintln(stderrW, "Missing -text option!")
		os.Exit(127)
	}

	args := flag.Args()
	if len(args) > 0 {
		fmt.Fprintln(stderrW, "Too many arguments!")
		os.Exit(127)
	}

	// Flag gets printed as a page
	mux := http.NewServeMux()
	mux.HandleFunc("/", httpLog(stdoutW, withAppHeaders(httpEcho(*textFlag))))

	// Health endpoint
	mux.HandleFunc("/health", withAppHeaders(httpHealth()))

	server := &http.Server{
		Addr:    *listenFlag,
		Handler: mux,
	}
	serverCh := make(chan struct{})
	go func() {
		log.Printf("[INFO] server is listening on %s\n", *listenFlag)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[ERR] server exited with: %s", err)
		}
		close(serverCh)
	}()

	doneCh := make(chan struct{})
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGUSR1)

	go handleSignal(doneCh, signalCh)
	// Wait for interrupt
	<-doneCh

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERR] failed to shutdown server: %s", err)
	}

	// If we got this far, it was an interrupt, so don't exit cleanly
	os.Exit(2)
}

func httpEcho(v string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, v)
	}
}

func httpHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"status":"ok"}`)
	}
}

func handleSignal(doneCh chan struct{}, sigCh chan os.Signal) {
	for {
		switch <-sigCh {
		case os.Interrupt:
			log.Printf("[INFO] received interrupt, shutting down...")
			doneCh <- struct{}{}
		case syscall.SIGUSR1:
			*enableLag = !*enableLag
			log.Printf("[INFO] received USR!. EnableLag: %s", strconv.FormatBool(*enableLag))
		}
	}
}
