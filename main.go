package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/warpcomdev/simgr/internal"
)

const (
	READ_TIMEOUT      = 5 * time.Second
	WRITE_TIMEOUT     = 5 * time.Second
	IDLE_TIMEOUT      = 5 * time.Minute
	GRACEFUL_SHUTDOWN = 30 * time.Second
	FILE_PATH         = "build"
	HTTP_PREFIX       = "/simgr"
)

func main() {

	config, err := internal.NewConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		flag.Usage()
		os.Exit(-1)
	}

	logger := internal.NewLogger(config.Quiet)
	server := http.Server{
		ErrorLog:     logger.Logger(),
		Addr:         fmt.Sprintf(":%d", config.Port),
		ReadTimeout:  READ_TIMEOUT,
		WriteTimeout: WRITE_TIMEOUT,
		IdleTimeout:  IDLE_TIMEOUT,
		Handler: func() http.Handler {
			internalServer := internal.NewServer(logger, os.DirFS(FILE_PATH), config.URL, true)
			// Pongo el servidor debajo de un prefijo
			prefixedServer := http.StripPrefix(HTTP_PREFIX, internalServer)
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Y redirijo al prefijo si se usa la ruta /
				if r.Method == http.MethodGet && (r.URL.Path == "/" || r.URL.Path == HTTP_PREFIX) {
					http.Redirect(w, r, HTTP_PREFIX+"/", http.StatusTemporaryRedirect)
					return
				}
				prefixedServer.ServeHTTP(w, r)
			})
		}(),
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("Listening on address %s", server.Addr)
		err := server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	ctx, cancelFunc := context.WithDeadline(context.Background(), time.Now().Add(GRACEFUL_SHUTDOWN))
	defer cancelFunc()
	if err := server.Shutdown(ctx); err != nil {
		panic(err)
	}
	wg.Wait()
}
