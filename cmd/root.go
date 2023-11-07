package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const appName = "tmp2p"

type appState struct {
	Log *zap.Logger
}

func NewRootCmd(log *zap.Logger) *cobra.Command {
	a := &appState{Log: log}

	var rootCmd = &cobra.Command{
		Use:   "tmp2p",
		Short: "tmp2p - a simple CLI to validate tm p2p addresses",
		Long:  `tmp2p is a CLI to validate tendermint p2p peer addresses`,
	}

	rootCmd.AddCommand(
		validateCmd(a),
		versionCmd(a),
	)

	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd(nil)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
