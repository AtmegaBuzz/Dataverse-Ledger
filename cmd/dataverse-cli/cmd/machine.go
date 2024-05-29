package cmd

import (
	"context"
	"dataverse/actions"
	"fmt"

	"github.com/spf13/cobra"
)

var machineCmd = &cobra.Command{
	Use: "machine",
	RunE: func(*cobra.Command, []string) error {
		return ErrMissingSubcommand
	},
}

var registerMachineCID = &cobra.Command{
	Use: "register-machine",
	RunE: func(*cobra.Command, []string) error {

		ctx := context.Background()
		_, _, factory, cli, scli, tcli, err := handler.DefaultActor()
		if err != nil {
			return err
		}

		machineCID, err := handler.Root().PromptString("Machine CID", 1, 256)
		if err != nil {
			return err
		}

		// Confirm action
		cont, err := handler.Root().PromptContinue()
		if !cont || err != nil {
			return err
		}

		project := &actions.CreateProject{
			MachineCID: []byte(machineCID),
		}

		// Generate transaction
		_, id, err := sendAndWait(ctx, nil, project, cli, scli, tcli, factory, true)

		if err != nil {
			fmt.Println("Error occured")
		}

		fmt.Println(id)

		return err

	},
}
