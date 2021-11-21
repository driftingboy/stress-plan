package main

import (
	"log"
	"stress-plan/cmd/run"

	"github.com/spf13/cobra"
)

func main() {
	mainCmd := &cobra.Command{
		Use:   "stp",
		Short: "Stressing Test Plan",
		Long:  "Stressing Test Plan",
	}

	mainCmd.AddCommand(run.RunCMD())

	if err := mainCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
