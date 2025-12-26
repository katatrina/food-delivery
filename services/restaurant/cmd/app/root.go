package app

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Short: "Restaurant service for food delivery platform",
}

func Execute() error {
	return rootCmd.Execute()
}
