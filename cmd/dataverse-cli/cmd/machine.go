package cmd

import (
	"context"
	"dataverse/actions"
	"dataverse/consts"
	"fmt"

	"github.com/ava-labs/hypersdk/codec"
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

		project := &actions.RegisterMachine{
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

var getregisterMachineCID = &cobra.Command{
	Use: "get-machine-cid",
	RunE: func(*cobra.Command, []string) error {

		ctx := context.Background()
		_, _, _, _, _, tcli, err := handler.DefaultActor()
		if err != nil {
			return err
		}

		id, _ := handler.Root().PromptID("register machine txid")

		ID, MachineCID, _ := tcli.MachineCID(ctx, id, false)

		addr, err := codec.AddressBech32(consts.HRP, codec.Address(ID))

		fmt.Println("Id: ", addr, ", MachineCID: ", string(MachineCID))

		return err

	},
}
