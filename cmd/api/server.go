package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) server() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.config.port),
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

		s := <-quit

		app.infoLog.Println("shuting down the server", map[string]string{
			"signal": s.String(),
		})

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		defer cancel()

		err := srv.Shutdown(ctx)

		if err != nil {
			shutdown <- err
			return
		}

		app.infoLog.Println("Completing the background tasks", map[string]string{
			"addr": srv.Addr,
		})

		app.wg.Wait()
		shutdown <- nil
	}()

	app.infoLog.Printf("cultured %s is starting on %s", app.config.env, srv.Addr)
	err := srv.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown

	if err != nil {
		return err
	}

	app.infoLog.Printf("cultured %s is stopping on %s", app.config.env, srv.Addr)

	return nil
}
