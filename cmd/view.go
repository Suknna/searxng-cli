/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"searxng-cli/internal/config"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// viewCmd represents the view command
var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View effective configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		_, cfg, err := config.LoadEffective(globalOverrides())
		if err != nil {
			return err
		}
		b, err := yaml.Marshal(cfg)
		if err != nil {
			return err
		}
		_, err = fmt.Fprint(cmd.OutOrStdout(), string(b))
		return err
	},
}

func init() {
	configCmd.AddCommand(viewCmd)
}
