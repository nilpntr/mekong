package main

import (
	"github.com/nilpntr/mekong/pkg/action"
	"github.com/nilpntr/mekong/pkg/cli"
	"github.com/spf13/cobra"
)

var settings = cli.New()

func newRootCmd(args []string) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:          "secretary",
		Short:        "Secretary is a tool to sync applicable secrets to other namespaces in k8s",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return action.NewRun(settings)
		},
	}
	flags := cmd.PersistentFlags()
	settings.AddFlags(flags)

	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.Parse(args)

	return cmd, nil
}