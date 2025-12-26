package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	httpserver "github.com/katatrina/food-delivery/services/restaurant/internal/infra/http"
	"github.com/katatrina/food-delivery/services/restaurant/internal/infra/postgres"
	"github.com/katatrina/food-delivery/services/restaurant/internal/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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

	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		var configErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configErr) {
			return fmt.Errorf("failed to load config: %w\n", err)
		}
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("failed to create db connection pool: %w", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

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
