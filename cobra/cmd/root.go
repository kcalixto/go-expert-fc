package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "print",
	Short: "just prints",
	Long:  `just prints`,
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		please, _ := cmd.Flags().GetBool("please")
		if !please {
			cmd.Println("no.")
			return
		}

		cmd.Printf("of course! Your name is: %s\n", name)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("name", "n", "", "Your name")
	rootCmd.PersistentFlags().BoolP("please", "p", false, "Polite flag")

	rootCmd.MarkFlagRequired("please")
}
