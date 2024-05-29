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

type RegisterMachine struct {
	MachineCID []byte `json:"machine_cid"`
}

func (*RegisterMachine) GetTypeID() uint8 {
	return registerMachineCIDID
}

func (*RegisterMachine) StateKeys(_ chain.Auth, txID ids.ID) []string {
	return []string{
		string(storage.RegisterMachineCIDKey(txID)),
	}
}

func (*RegisterMachine) StateKeysMaxChunks() []uint16 {
	return []uint16{storage.MachineCIDChunks}
}

func (*RegisterMachine) OutputsWarpMessage() bool {
	return false
}

func (c *RegisterMachine) Execute(
	ctx context.Context,
	_ chain.Rules,
	mu state.Mutable,
	_ int64,
	auth chain.Auth,
	txID ids.ID,
	_ bool,
) (bool, uint64, []byte, *warp.UnsignedMessage, error) {

	if len(c.MachineCID) == 0 {
		return false, RegisterMachineComputeUnits, OutputRegisterMachineNotProvided, nil, nil
	}

	// It should only be possible to overwrite an existing asset if there is
	// a hash collision.
	if err := storage.SetMachineCID(ctx, mu, txID, c.MachineCID); err != nil {
		return false, RegisterMachineComputeUnits, utils.ErrBytes(err), nil, nil
	}
	return true, RegisterMachineComputeUnits, nil, nil, nil
}

func (*RegisterMachine) MaxComputeUnits(chain.Rules) uint64 {
	return RegisterMachineComputeUnits
}

func (c *RegisterMachine) Size() int {
	// TODO: add small bytes (smaller int prefix)
	return (codec.BytesLen(c.MachineCID))

}

func (c *RegisterMachine) Marshal(p *codec.Packer) {
	p.PackBytes(c.MachineCID)
}

func UnmarshalRegisterMachineCID(p *codec.Packer, _ *warp.Message) (chain.Action, error) {

	var create RegisterMachine

	p.UnpackBytes(MachineCIDUnits, true, &create.MachineCID)

	return &create, p.Err()

}

func (*RegisterMachine) ValidRange(chain.Rules) (int64, int64) {
	// Returning -1, -1 means that the action is always valid.
	return -1, -1
}
