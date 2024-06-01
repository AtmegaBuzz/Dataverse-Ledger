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

var _ chain.Action = (*NotarizeData)(nil)

type NotarizeData struct {
	MachineAttestTx []byte `json:"machine_attest_tx"`
	DataCID         []byte `json:"data_cid"`
	DataType        []byte `json:"data_type"`
	DataOwnerAddr   []byte `json:"data_owner_addr"`
}

func (*NotarizeData) GetTypeID() uint8 {
	return notarizeDataID
}

func (*NotarizeData) StateKeys(_ chain.Auth, txID ids.ID) []string {
	return []string{
		string(storage.AttestMachineKey(txID)),
	}
}

func (*NotarizeData) StateKeysMaxChunks() []uint16 {
	return []uint16{storage.DataCIDChunks}
}

func (*NotarizeData) OutputsWarpMessage() bool {
	return false
}

func (c *NotarizeData) Execute(
	ctx context.Context,
	_ chain.Rules,
	mu state.Mutable,
	_ int64,
	auth chain.Auth,
	txID ids.ID,
	_ bool,
) (bool, uint64, []byte, *warp.UnsignedMessage, error) {

	if len(c.DataOwnerAddr) == 0 {
		return false, AttestMachineComputeUnits, OutputInvalidMachineAddressLen, nil, nil
	}

	// It should only be possible to overwrite an existing asset if there is
	// a hash collision.
	if err := storage.NotarizeData(ctx, mu, txID, c.MachineAttestTx, c.DataOwnerAddr, c.DataCID, c.DataType); err != nil {
		return false, AttestMachineComputeUnits, utils.ErrBytes(err), nil, nil
	}
	return true, AttestMachineComputeUnits, nil, nil, nil
}

func (*NotarizeData) MaxComputeUnits(chain.Rules) uint64 {
	return NotarizeDataComputeUnits
}

func (c *NotarizeData) Size() int {
	// TODO: add small bytes (smaller int prefix)
	return (codec.BytesLen(c.MachineAttestTx) +
		codec.BytesLen(c.DataOwnerAddr) +
		codec.BytesLen(c.DataCID) +
		codec.BytesLen(c.DataType))

}

func (c *NotarizeData) Marshal(p *codec.Packer) {
	p.PackBytes(c.MachineAttestTx)
	p.PackBytes(c.DataOwnerAddr)
	p.PackBytes(c.DataCID)
	p.PackBytes(c.DataType)

}

func UnmarshalNotarizeData(p *codec.Packer, _ *warp.Message) (chain.Action, error) {

	var create NotarizeData

	p.UnpackBytes(MachineAttestTxUnits, true, &create.MachineAttestTx)
	p.UnpackBytes(DataOwnerAddrUnits, true, &create.DataOwnerAddr)
	p.UnpackBytes(DataCIDUnits, true, &create.DataCID)
	p.UnpackBytes(DataTypeUnits, true, &create.DataType)

	return &create, p.Err()

}

func (*NotarizeData) ValidRange(chain.Rules) (int64, int64) {
	// Returning -1, -1 means that the action is always valid.
	return -1, -1
}
