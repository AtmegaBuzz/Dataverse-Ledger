package cmd

import (
	"context"
	"dataverse/actions"
	"dataverse/consts"
	"dataverse/storage"
	"fmt"

	"github.com/ava-labs/avalanchego/ids"
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

		machineCID, err := handler.Root().PromptString("Machine CID", 66, 66)
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
		te, _, err := sendAndWait(ctx, nil, project, cli, scli, tcli, factory, true)

		if err != nil {
			fmt.Println("Error occured")
		}

		// fmt.Println(id)
		fmt.Println(te)

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

		ID, MachineCID, err := tcli.MachineCID(ctx, id, false)

		if err != nil {
			return err
		}

		addr, err := codec.AddressBech32(consts.HRP, codec.Address(ID))

		fmt.Println("ID", addr, ", MachineCID: ", string(MachineCID))

		return err

	},
}

var attestMachine = &cobra.Command{
	Use: "attest-machine",
	RunE: func(*cobra.Command, []string) error {

		ctx := context.Background()
		_, _, factory, cli, scli, tcli, err := handler.DefaultActor()
		if err != nil {
			return err
		}

		address, err := handler.Root().PromptString("Machine Address", 44, 44)
		if err != nil {
			return err
		}

		machine_category, err := handler.Root().PromptString("Machine Category", 1, 100)
		if err != nil {
			return err
		}

		machine_manufacturer, err := handler.Root().PromptString("Machine Manufacturer", 1, 100)
		if err != nil {
			return err
		}

		machineCID, err := handler.Root().PromptString("Machine CID", 66, 66)
		if err != nil {
			return err
		}

		// Confirm action
		cont, err := handler.Root().PromptContinue()
		if !cont || err != nil {
			return err
		}

		project := &actions.AttestMachine{
			MachineAddress:      []byte(address),
			MachineCategory:     []byte(machine_category),
			MachineManufacturer: []byte(machine_manufacturer),
			MachineCID:          []byte(machineCID),
		}

		// Generate transaction
		te, _, err := sendAndWait(ctx, nil, project, cli, scli, tcli, factory, true)

		if err != nil {
			fmt.Println("Error occured")
		}

		fmt.Println(te)

		return err

	},
}

var getAttestedachineCID = &cobra.Command{
	Use: "get-attested-machine",
	RunE: func(*cobra.Command, []string) error {

		ctx := context.Background()
		_, _, _, _, _, tcli, err := handler.DefaultActor()
		if err != nil {
			return err
		}

		id, _ := handler.Root().PromptID("attestation txid")

		ID, MachineAddress, MachineCategory, MachineManufacturer, MachineCID, err := tcli.AttestMachine(ctx, id, false)

		if err != nil {
			return err
		}

		addr, err := codec.AddressBech32(consts.HRP, codec.Address(ID))

		fmt.Println("ID", addr, ", MachineAddress: ", string(MachineAddress), ", MachineCategory: ", string(MachineCategory), ", MachineManufacturer: ", string(MachineManufacturer), ", MachineCID: ", string(MachineCID))

		return err

	},
}

var notarizeData = &cobra.Command{
	Use: "notarize",
	RunE: func(*cobra.Command, []string) error {

		ctx := context.Background()
		_, _, factory, cli, scli, tcli, err := handler.DefaultActor()
		if err != nil {
			return err
		}

		attestationTx, err := handler.Root().PromptID("attestation txid")
		if err != nil {
			return err
		}

		creator, err := handler.Root().PromptString("Machine Address", 44, 44)
		if err != nil {
			return err
		}

		notarizeType := "/dataverse.asset.MsgNotarizedAsset"

		dataCid, err := handler.Root().PromptString("Data CID", 59, 59)
		if err != nil {
			return err
		}

		// Confirm action
		cont, err := handler.Root().PromptContinue()
		if !cont || err != nil {
			return err
		}

		project := &actions.NotarizeData{
			MachineAttestTx: storage.NotarizeDataKey(attestationTx),
			DataCID:         []byte(dataCid),
			DataType:        []byte(notarizeType),
			DataOwnerAddr:   []byte(creator),
		}

		// Generate transaction
		te, _, err := sendAndWait(ctx, nil, project, cli, scli, tcli, factory, true)

		if err != nil {
			fmt.Println("Error occured")
		}

		fmt.Println(te)

		return err

	},
}

var getNotarizeData = &cobra.Command{
	Use: "get-notarize-data",
	RunE: func(*cobra.Command, []string) error {

		ctx := context.Background()
		_, _, _, _, _, tcli, err := handler.DefaultActor()
		if err != nil {
			return err
		}

		id, _ := handler.Root().PromptID("notarized txid")

		ID, MachineAttestTx, DataOwnerAddr, DataCID, DataType, err := tcli.NotarizeData(ctx, id, false)

		if err != nil {
			return err
		}

		addr, err := codec.AddressBech32(consts.HRP, codec.Address(ID))

		fmt.Println("ID", addr, ", MachineAttestTx: ", ids.ID(MachineAttestTx), ", DataCID: ", string(DataCID), ", DataType: ", string(DataType), ", DataOwnerAddr: ", string(DataOwnerAddr))

		return err

	},
}
