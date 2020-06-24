package serve

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
	"fmt"

	"github.com/nuucactus/sms-api/middleware"
	"github.com/nuucactus/sms-api/router"

	"github.com/rs/zerolog/log"
	"github.com/justinas/alice"
	"github.com/spf13/cobra"
)

// RunServe the main event loop for the service
func RunServe() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {

		ip, err := cmd.Flags().GetString("ip")
		if err != nil {
			log.Warn().Msg("Missing ip")
			os.Exit(1)
		}

		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Warn().Msg("Missing port")
			os.Exit(1)
		}

		addr := fmt.Sprintf("%s:%d", ip, port)

		route := router.NewRouter()

		chain := alice.New( middleware.Context, middleware.Logging )

		srv := &http.Server{
			Addr:         addr,
			// Good practice to set timeouts to avoid Slowloris attacks.
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
			Handler: chain.Then(route), // Pass our instance of gorilla/mux in.
		}

		// Run our server in a goroutine so that it doesn't block.
		go func() {
			log.Info().Msg("Listening on " + addr)
			if err := srv.ListenAndServe(); err != nil {
				log.Error().Err(err)
			}
		}()

		c := make(chan os.Signal, 1)
		// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
		// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
		signal.Notify(c, os.Interrupt)

		// Block until we receive our signal.
		<-c

		// Create a deadline to wait for.
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 15)
		defer cancel()
		// Doesn't block if no connections, but will otherwise wait
		// until the timeout deadline.
		srv.Shutdown(ctx)
		// Optionally, you could run srv.Shutdown in a goroutine and block on
		// <-ctx.Done() if your application should wait for other services
		// to finalize based on context cancellation.
		log.Info().Msg("shutting down")
		os.Exit(0)
	}
}
