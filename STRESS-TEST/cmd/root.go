package cmd

import (
	"stresser/requester"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "stresser",
		Short: "stressed is a simple CLI designed to load test an REST API",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd:  true,
			DisableDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()

			url, _ := cmd.Flags().GetString("url")
			totalRequestCount, _ := cmd.Flags().GetInt("requests")
			concurrency, _ := cmd.Flags().GetInt("concurrency")

			requester.Request(ctx, url, concurrency, totalRequestCount)
		},
	}
)

func init() {
	rootCmd.PersistentFlags().String("url", "", "set target URL")
	rootCmd.MarkPersistentFlagRequired("url")

	rootCmd.PersistentFlags().Int("requests", 100, "set number of requests")
	rootCmd.PersistentFlags().Int("concurrency", 10, "set number of concurrent requests")
	rootCmd.MarkFlagsRequiredTogether("requests", "concurrency")

	for _, command := range []*cobra.Command{
		version(),
	} {
		rootCmd.AddCommand(command)
	}
}

func Execute() {
	rootCmd.Execute()
}

func version() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("stresser: %s\n", "v0.0.1")
		},
	}
}
