package cmd

import (
	"IMP/app/database/seeders"
	"fmt"
	"github.com/spf13/cobra"
	"reflect"
)

var seedListCmd = &cobra.Command{
	Use:   "seed:list",
	Short: "",
	Long:  "",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		showSeedersList()
	},
}

func init() {
	rootCmd.AddCommand(seedListCmd)
}

// showSeedersList returns list of all available
func showSeedersList() {
	models := seeders.AvailableModels

	fmt.Println("Available models for seeding:")
	for _, model := range models {
		fmt.Printf("- %s\n", reflect.TypeOf(model.Model()).Name())
	}
}
