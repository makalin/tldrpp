package main

import (
	"fmt"
	"os"

	"github.com/makalin/tldrpp/internal/app"
	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "unknown"
	date    = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "tldrpp",
		Short: "Interactive cheat-sheets with fuzzy search and inline editing",
		Long: `tldr++ is a terminal UI that lets you fuzzy-search pages, edit placeholders inline, 
then paste or execute the final command.`,
		Version: fmt.Sprintf("%s (commit: %s, built: %s)", version, commit, date),
	}

	var initCmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize tldr++ by downloading page index",
		Run: func(cmd *cobra.Command, args []string) {
			if err := app.Initialize(); err != nil {
				fmt.Fprintf(os.Stderr, "Error initializing tldr++: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("tldr++ initialized successfully!")
		},
	}

	var updateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update tldr pages cache",
		Run: func(cmd *cobra.Command, args []string) {
			if err := app.UpdateCache(); err != nil {
				fmt.Fprintf(os.Stderr, "Error updating cache: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("Cache updated successfully!")
		},
	}

	var renderCmd = &cobra.Command{
		Use:   "render [command]",
		Short: "Render command with placeholders filled",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vars, _ := cmd.Flags().GetStringToString("vars")
			if err := app.RenderCommand(args[0], vars); err != nil {
				fmt.Fprintf(os.Stderr, "Error rendering command: %v\n", err)
				os.Exit(1)
			}
		},
	}
	renderCmd.Flags().StringToString("vars", nil, "Variables to substitute in placeholders")

	var execCmd = &cobra.Command{
		Use:   "exec [command]",
		Short: "Execute command with placeholders filled",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			vars, _ := cmd.Flags().GetStringToString("vars")
			if err := app.ExecuteCommand(args[0], vars); err != nil {
				fmt.Fprintf(os.Stderr, "Error executing command: %v\n", err)
				os.Exit(1)
			}
		},
	}
	execCmd.Flags().StringToString("vars", nil, "Variables to substitute in placeholders")

	var pluginCmd = &cobra.Command{
		Use:   "plugin",
		Short: "Plugin commands",
	}

	var submitCmd = &cobra.Command{
		Use:   "submit",
		Short: "Submit current example to tldr-pages",
		Run: func(cmd *cobra.Command, args []string) {
			if err := app.SubmitToTldr(); err != nil {
				fmt.Fprintf(os.Stderr, "Error submitting to tldr: %v\n", err)
				os.Exit(1)
			}
		},
	}

	pluginCmd.AddCommand(submitCmd)

	// Global flags
	rootCmd.PersistentFlags().StringP("platform", "p", "", "Platform filter (common, linux, osx, sunos, windows, android)")
	rootCmd.PersistentFlags().StringP("theme", "t", "dark", "Theme (light, dark, solarized)")
	rootCmd.PersistentFlags().BoolP("dev", "d", false, "Development mode")

	rootCmd.AddCommand(initCmd, updateCmd, renderCmd, execCmd, pluginCmd)

	// Default action: run the TUI
	rootCmd.Run = func(cmd *cobra.Command, args []string) {
		platform, _ := cmd.Flags().GetString("platform")
		theme, _ := cmd.Flags().GetString("theme")
		dev, _ := cmd.Flags().GetBool("dev")

		var searchQuery string
		if len(args) > 0 {
			searchQuery = args[0]
		}

		if err := app.RunTUI(searchQuery, platform, theme, dev); err != nil {
			fmt.Fprintf(os.Stderr, "Error running tldr++: %v\n", err)
			os.Exit(1)
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}