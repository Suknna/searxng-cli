/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"searxng-cli/internal/config"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化默认配置文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := cfgFile
		if path == "" {
			var err error
			path, err = config.DefaultConfigPath()
			if err != nil {
				return err
			}
		}
		if err := config.WriteDefault(path); err != nil {
			return err
		}
		_, err := fmt.Fprintln(cmd.OutOrStdout(), path)
		return err
	},
}

func init() {
	configCmd.AddCommand(initCmd)
}
