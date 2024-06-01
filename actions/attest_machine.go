// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package actions

import (
	"context"

	"dataverse/storage"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/vms/platformvm/warp"
	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/codec"
	"github.com/ava-labs/hypersdk/state"
	"github.com/ava-labs/hypersdk/utils"
)

var _ chain.Action = (*RegisterMachine)(nil)

type AttestMachine struct {
	MachineAddress      []byte `json:"machine_address"`
	MachineCategory     []byte `json:"machine_category"`
	MachineManufacturer []byte `json:"machine_manufacturer"`
	MachineCID          []byte `json:"machine_cid"`
}

func (*AttestMachine) GetTypeID() uint8 {
	return attestMachineID
}

func (*AttestMachine) StateKeys(_ chain.Auth, txID ids.ID) []string {
	return []string{
		string(storage.AttestMachineKey(txID)),
	}
}

func (*AttestMachine) StateKeysMaxChunks() []uint16 {
	return []uint16{storage.MachineCategoryChunks}
}

func (*AttestMachine) OutputsWarpMessage() bool {
	return false
}

func (c *AttestMachine) Execute(
	ctx context.Context,
	_ chain.Rules,
	mu state.Mutable,
	_ int64,
	auth chain.Auth,
	txID ids.ID,
	_ bool,
) (bool, uint64, []byte, *warp.UnsignedMessage, error) {

	if len(c.MachineAddress) != 44 {
		return false, AttestMachineComputeUnits, OutputInvalidMachineAddressLen, nil, nil
	}

	if len(c.MachineCategory) > 100 {
		return false, AttestMachineComputeUnits, OutputInvalidMachineCategoryLen, nil, nil
	}

	if len(c.MachineManufacturer) > 100 {
		return false, AttestMachineComputeUnits, OutputInvalidMachineManufacturerLen, nil, nil
	}

	if len(c.MachineCID) != 66 {
		return false, AttestMachineComputeUnits, OutputInvalidMachineCIDLen, nil, nil
	}

	// It should only be possible to overwrite an existing asset if there is
	// a hash collision.
	if err := storage.AttestMachine(ctx, mu, txID, c.MachineAddress, c.MachineCategory, c.MachineManufacturer, c.MachineCID); err != nil {
		return false, AttestMachineComputeUnits, utils.ErrBytes(err), nil, nil
	}
	return true, AttestMachineComputeUnits, nil, nil, nil
}

func (*AttestMachine) MaxComputeUnits(chain.Rules) uint64 {
	return AttestMachineComputeUnits
}

func (c *AttestMachine) Size() int {
	// TODO: add small bytes (smaller int prefix)
	return (codec.BytesLen(c.MachineAddress) +
		codec.BytesLen(c.MachineCategory) +
		codec.BytesLen(c.MachineManufacturer) +
		codec.BytesLen(c.MachineCID))

}

func (c *AttestMachine) Marshal(p *codec.Packer) {
	p.PackBytes(c.MachineAddress)
	p.PackBytes(c.MachineCategory)
	p.PackBytes(c.MachineManufacturer)
	p.PackBytes(c.MachineCID)

}

func UnmarshalAttestMachineCID(p *codec.Packer, _ *warp.Message) (chain.Action, error) {

	var create AttestMachine

	p.UnpackBytes(MachineAddressUnits, true, &create.MachineAddress)
	p.UnpackBytes(MachineCategoryUnits, true, &create.MachineCategory)
	p.UnpackBytes(MachineManufacturerUnits, true, &create.MachineManufacturer)
	p.UnpackBytes(MachineCIDUnits, true, &create.MachineCID)

	return &create, p.Err()

}

func (*AttestMachine) ValidRange(chain.Rules) (int64, int64) {
	// Returning -1, -1 means that the action is always valid.
	return -1, -1
}
