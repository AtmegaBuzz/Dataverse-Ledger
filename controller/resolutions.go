// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package controller

import (
	"context"

	"dataverse/genesis"
	"dataverse/orderbook"
	"dataverse/storage"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/trace"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/codec"
)

func (c *Controller) Genesis() *genesis.Genesis {
	return c.genesis
}

func (c *Controller) Logger() logging.Logger {
	return c.inner.Logger()
}

func (c *Controller) Tracer() trace.Tracer {
	return c.inner.Tracer()
}

func (c *Controller) GetTransaction(
	ctx context.Context,
	txID ids.ID,
) (bool, int64, bool, chain.Dimensions, uint64, error) {
	return storage.GetTransaction(ctx, c.metaDB, txID)
}

func (c *Controller) GetAssetFromState(
	ctx context.Context,
	asset ids.ID,
) (bool, []byte, uint8, []byte, uint64, codec.Address, bool, error) {
	return storage.GetAssetFromState(ctx, c.inner.ReadState, asset)
}

func (c *Controller) GetBalanceFromState(
	ctx context.Context,
	addr codec.Address,
	asset ids.ID,
) (uint64, error) {
	return storage.GetBalanceFromState(ctx, c.inner.ReadState, addr, asset)
}

func (c *Controller) Orders(pair string, limit int) []*orderbook.Order {
	return c.orderBook.Orders(pair, limit)
}

func (c *Controller) GetOrderFromState(
	ctx context.Context,
	orderID ids.ID,
) (
	bool, // exists
	ids.ID, // in
	uint64, // inTick
	ids.ID, // out
	uint64, // outTick
	uint64, // remaining
	codec.Address, // owner
	error,
) {
	return storage.GetOrderFromState(ctx, c.inner.ReadState, orderID)
}

func (c *Controller) GetLoanFromState(
	ctx context.Context,
	asset ids.ID,
	destination ids.ID,
) (uint64, error) {
	return storage.GetLoanFromState(ctx, c.inner.ReadState, asset, destination)
}

func (c *Controller) GetProjectFromState(
	ctx context.Context,
	project ids.ID,

) (bool, storage.ProjectData, error) {
	return storage.GetProjectFromState(ctx, c.inner.ReadState, project)
}

func (c *Controller) GetUpdateFromState(
	ctx context.Context,
	update ids.ID,

) (bool, storage.UpdateData, error) {
	return storage.GetUpdateFromState(ctx, c.inner.ReadState, update)
}

func (c *Controller) GetMachineCID(
	ctx context.Context,
	machinCIDID ids.ID,

) (bool, storage.RegisterMachineCIDData, error) {
	return storage.GetMachineCID(ctx, c.inner.ReadState, machinCIDID)
}

func (c *Controller) GetAttestMachine(
	ctx context.Context,
	tx ids.ID,

) (bool, storage.AttestMachineData, error) {
	return storage.GetAttestMachine(ctx, c.inner.ReadState, tx)
}
