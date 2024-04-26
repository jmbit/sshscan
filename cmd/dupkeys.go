/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/jmbit/sshscan/internal/duplicates"
	"github.com/spf13/cobra"
)

// dupkeysCmd represents the dupkeys command
var dupkeysCmd = &cobra.Command{
	Use:   "dupkeys",
	Short: "Finds duplicate ssh keys and prints them",
	Long:  `takes a subnet and iterates through its IP addresses`,
	Run: func(cmd *cobra.Command, args []string) {
		dups, byKey := duplicates.FindDuplicates(args)
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Print("\n----\nResults\n")
		fmt.Fprintf(w, "#\tHost\tKey\n")
		for i, result := range dups {
			fmt.Fprintf(w, "%d\t%s\t%s\n", i, result.IP, result.Key)

		}
		w.Flush()
		w = tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		fmt.Fprintf(w, "Key\tHosts\n")
		for key, list := range byKey {
			fmt.Fprintf(w, "%s\t%v\n", key, list)
		}
		w.Flush()
	},
}

func init() {
	rootCmd.AddCommand(dupkeysCmd)
}
