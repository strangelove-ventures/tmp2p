package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const flagJSON = "json"

var (
	// Version defines the application version (defined at compile time)
	Version = ""
	Commit  = ""
	Dirty   = ""
)

type versionInfo struct {
	Version  string `json:"version" yaml:"version"`
	Commit   string `json:"commit" yaml:"commit"`
	CometBFT string `json:"cometbft" yaml:"cometbft"`
	Go       string `json:"go" yaml:"go"`
}

func versionCmd(a *appState) *cobra.Command {
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Print version info",
		Args:    cobra.NoArgs,
		Example: strings.TrimSpace(fmt.Sprintf(`
$ %s version --json
$ %s v`,
			appName, appName,
		)),
		RunE: func(cmd *cobra.Command, args []string) error {
			jsn, err := cmd.Flags().GetBool(flagJSON)
			if err != nil {
				return err
			}

			cometBFT := "(unable to determine)"
			if bi, ok := debug.ReadBuildInfo(); ok {
				for _, dep := range bi.Deps {
					if dep.Path == "github.com/cometbft/cometbft" {
						cometBFT = dep.Version
						break
					}
				}
			}

			commit := Commit
			if Dirty != "0" {
				commit += " (dirty)"
			}

			verInfo := versionInfo{
				Version:  Version,
				Commit:   commit,
				CometBFT: cometBFT,
				Go:       fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
			}

			var bz []byte
			if jsn {
				bz, err = json.Marshal(verInfo)
			} else {
				bz, err = yaml.Marshal(&verInfo)
			}

			fmt.Fprintln(cmd.OutOrStdout(), string(bz))
			return err
		},
	}

	versionCmd.Flags().BoolP(flagJSON, "j", false, "return the response in json format")

	return versionCmd
}
