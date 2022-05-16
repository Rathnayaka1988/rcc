package cmd

import (
	"os"

	"github.com/robocorp/rcc/common"
	"github.com/robocorp/rcc/pretty"

	"github.com/spf13/cobra"
)

var (
	enableShared bool
)

var holotreeSharedCommand = &cobra.Command{
	Use:   "shared",
	Short: "Enable shared holotree usage.",
	Long:  "Enable shared holotree usage.",
	Run: func(cmd *cobra.Command, args []string) {
		if common.DebugFlag {
			defer common.Stopwatch("Enabling shared holotree lasted").Report()
		}
		if os.Geteuid() > 0 {
			pretty.Warning("Running this command might need sudo/root access rights. Still, trying ...")
		}
		osSpecificHolotreeSharing(enableShared)
		pretty.Ok()
	},
}

func init() {
	holotreeSharedCommand.Flags().BoolVarP(&enableShared, "enable", "e", false, "Enable shared holotree environments between users. Currently cannot be undone.")
	holotreeCmd.AddCommand(holotreeSharedCommand)
}
