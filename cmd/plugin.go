package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"shield-cli/plugin"

	"github.com/spf13/cobra"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin",
	Short: "Manage Shield plugins",
	Long:  "Install, list, and remove Shield plugins that extend protocol support.",
}

var pluginFrom string

var pluginAddCmd = &cobra.Command{
	Use:     "add <name>",
	Short:   "Install a plugin",
	Example: "  shield plugin add mysql\n  shield plugin add postgres\n  shield plugin add mysql --from ./shield-plugin-mysql",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		fmt.Printf("  Installing plugin %q...\n", name)

		var info *plugin.PluginInfo
		var err error

		if pluginFrom != "" {
			info, err = plugin.InstallFromLocal(name, pluginFrom)
		} else {
			info, err = plugin.Install(name)
		}
		if err != nil {
			return err
		}

		fmt.Printf("  \033[32m✓ Installed %s %s\033[0m\n", info.Name, info.Version)
		fmt.Printf("    Protocols: %v\n", info.Protocols)
		fmt.Printf("    Binary:    %s\n", info.Binary)
		return nil
	},
}

var pluginListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed plugins",
	RunE: func(cmd *cobra.Command, args []string) error {
		reg, err := plugin.LoadRegistry()
		if err != nil {
			return err
		}

		if len(reg.Plugins) == 0 {
			fmt.Println("  No plugins installed.")
			fmt.Println()
			fmt.Printf("  Available plugins: %s\n", plugin.AvailablePluginNames())
			fmt.Println("  Install with: shield plugin add <name>")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintln(w, "  NAME\tVERSION\tPROTOCOLS\tINSTALLED")
		for _, p := range reg.Plugins {
			protocols := ""
			for i, pr := range p.Protocols {
				if i > 0 {
					protocols += ", "
				}
				protocols += pr
			}
			fmt.Fprintf(w, "  %s\t%s\t%s\t%s\n", p.Name, p.Version, protocols, p.InstalledAt)
		}
		w.Flush()
		return nil
	},
}

var pluginUpgradeCmd = &cobra.Command{
	Use:     "upgrade [name]",
	Short:   "Upgrade installed plugins to latest version",
	Example: "  shield plugin upgrade mysql\n  shield plugin upgrade  # upgrade all",
	Args:    cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 1 {
			name := args[0]
			fmt.Printf("  Checking for updates: %s...\n", name)
			result, err := plugin.Upgrade(name)
			if err != nil {
				return err
			}
			if result.Upgraded {
				fmt.Printf("  \033[32m✓ Upgraded %s: %s → %s\033[0m\n", result.Name, result.CurrentVersion, result.LatestVersion)
			} else {
				fmt.Printf("  %s is already up to date (%s)\n", result.Name, result.CurrentVersion)
			}
			return nil
		}

		// Upgrade all
		fmt.Println("  Checking for updates...")
		results, err := plugin.UpgradeAll()
		if err != nil {
			return err
		}
		if len(results) == 0 {
			fmt.Println("  No plugins installed.")
			return nil
		}
		for _, r := range results {
			if r.Upgraded {
				fmt.Printf("  \033[32m✓ Upgraded %s: %s → %s\033[0m\n", r.Name, r.CurrentVersion, r.LatestVersion)
			} else {
				fmt.Printf("  %s: %s (up to date)\n", r.Name, r.CurrentVersion)
			}
		}
		return nil
	},
}

var pluginRemoveCmd = &cobra.Command{
	Use:     "remove <name>",
	Aliases: []string{"rm"},
	Short:   "Remove an installed plugin",
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		reg, err := plugin.LoadRegistry()
		if err != nil {
			return err
		}

		if err := reg.Remove(name); err != nil {
			return err
		}

		fmt.Printf("  \033[32m✓ Removed plugin %q\033[0m\n", name)
		return nil
	},
}

func init() {
	pluginAddCmd.Flags().StringVar(&pluginFrom, "from", "", "Install from a local binary path")
	pluginCmd.AddCommand(pluginAddCmd)
	pluginCmd.AddCommand(pluginListCmd)
	pluginCmd.AddCommand(pluginRemoveCmd)
	pluginCmd.AddCommand(pluginUpgradeCmd)
	rootCmd.AddCommand(pluginCmd)
}
