package cmd

import (
	"IMP/app/api"
	"github.com/spf13/cobra"
	"log"
	"net/http"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Starts the server",
	Long:  "Starts the Http server",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		mux := api.Serve()

		log.Println("Starting server on :8080")
		log.Fatal(http.ListenAndServe(":8080", mux))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
