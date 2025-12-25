package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	httpserver "github.com/katatrina/food-delivery/services/restaurant/internal/infra/http"
	"github.com/katatrina/food-delivery/services/restaurant/internal/infra/postgres"
	"github.com/katatrina/food-delivery/services/restaurant/internal/service"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the HTTP server",
	RunE:  runServe,
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "HTTP server port")
}

func runServe(cmd *cobra.Command, args []string) error {
	port, _ := cmd.Flags().GetInt("port")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer pool.Close()

	store := postgres.NewStore(pool)
	restaurantSvc := service.NewRestaurantService(store)
	restaurantHandler := httpserver.NewRestaurantHandler(restaurantSvc)
	server := httpserver.NewServer(restaurantHandler)

	log.Printf("Starting app on :%d", port)
	if err = http.ListenAndServe(fmt.Sprintf(":%d", port), server.Router()); err != nil {
		return fmt.Errorf("app failed: %w", err)
	}

	return nil
}
