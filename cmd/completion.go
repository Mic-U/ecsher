/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// completionCmd represents the completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(ecsher completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ ecsher completion bash > /etc/bash_completion.d/ecsher
  # macOS:
  $ ecsher completion bash > /usr/local/etc/bash_completion.d/ecsher

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ ecsher completion zsh > "${fpath[1]}/_ecsher"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ ecsher completion fish | source

  # To load completions for each session, execute once:
  $ ecsher completion fish > ~/.config/fish/completions/ecsher.fish

PowerShell:

  PS> ecsher completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> ecsher completion powershell > ecsher.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		err := runCompletion(cmd, args[0])
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}

func runCompletion(cmd *cobra.Command, arg string) error {
	switch arg {
	case "bash":
		return cmd.Root().GenBashCompletion(os.Stdout)
	case "zsh":
		return cmd.Root().GenZshCompletion(os.Stdout)
	case "fish":
		return cmd.Root().GenFishCompletion(os.Stdout, true)
	case "powershell":
		return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
	default:
		return nil
	}
}
