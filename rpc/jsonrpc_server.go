// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package rpc

import (
	"net/http"

	"github.com/ava-labs/avalanchego/ids"

	"dataverse/consts"
	"dataverse/genesis"
	"dataverse/orderbook"

	"github.com/ava-labs/hypersdk/chain"
	"github.com/ava-labs/hypersdk/codec"
)

type JSONRPCServer struct {
	c Controller
}

func NewJSONRPCServer(c Controller) *JSONRPCServer {
	return &JSONRPCServer{c}
}

type GenesisReply struct {
	Genesis *genesis.Genesis `json:"genesis"`
}

func (j *JSONRPCServer) Genesis(_ *http.Request, _ *struct{}, reply *GenesisReply) (err error) {
	reply.Genesis = j.c.Genesis()
	return nil
}

type TxArgs struct {
	TxID ids.ID `json:"txId"`
}

type TxReply struct {
	Timestamp int64            `json:"timestamp"`
	Success   bool             `json:"success"`
	Units     chain.Dimensions `json:"units"`
	Fee       uint64           `json:"fee"`
}

func (j *JSONRPCServer) Tx(req *http.Request, args *TxArgs, reply *TxReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Tx")
	defer span.End()

	found, t, success, units, fee, err := j.c.GetTransaction(ctx, args.TxID)
	if err != nil {
		return err
	}
	if !found {
		return ErrTxNotFound
	}
	reply.Timestamp = t
	reply.Success = success
	reply.Units = units
	reply.Fee = fee
	return nil
}

type AssetArgs struct {
	Asset ids.ID `json:"asset"`
}

type AssetReply struct {
	Symbol   []byte `json:"symbol"`
	Decimals uint8  `json:"decimals"`
	Metadata []byte `json:"metadata"`
	Supply   uint64 `json:"supply"`
	Owner    string `json:"owner"`
	Warp     bool   `json:"warp"`
}

func (j *JSONRPCServer) Asset(req *http.Request, args *AssetArgs, reply *AssetReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Asset")
	defer span.End()

	exists, symbol, decimals, metadata, supply, owner, warp, err := j.c.GetAssetFromState(ctx, args.Asset)
	if err != nil {
		return err
	}
	if !exists {
		return ErrAssetNotFound
	}
	reply.Symbol = symbol
	reply.Decimals = decimals
	reply.Metadata = metadata
	reply.Supply = supply
	reply.Owner = codec.MustAddressBech32(consts.HRP, owner)
	reply.Warp = warp
	return err
}

type BalanceArgs struct {
	Address string `json:"address"`
	Asset   ids.ID `json:"asset"`
}

type BalanceReply struct {
	Amount uint64 `json:"amount"`
}

func (j *JSONRPCServer) Balance(req *http.Request, args *BalanceArgs, reply *BalanceReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Balance")
	defer span.End()

	addr, err := codec.ParseAddressBech32(consts.HRP, args.Address)
	if err != nil {
		return err
	}
	balance, err := j.c.GetBalanceFromState(ctx, addr, args.Asset)
	if err != nil {
		return err
	}
	reply.Amount = balance
	return err
}

type OrdersArgs struct {
	Pair string `json:"pair"`
}

type OrdersReply struct {
	Orders []*orderbook.Order `json:"orders"`
}

func (j *JSONRPCServer) Orders(req *http.Request, args *OrdersArgs, reply *OrdersReply) error {
	_, span := j.c.Tracer().Start(req.Context(), "Server.Orders")
	defer span.End()

	reply.Orders = j.c.Orders(args.Pair, ordersToSend)
	return nil
}

type GetOrderArgs struct {
	OrderID ids.ID `json:"orderID"`
}

type GetOrderReply struct {
	Order *orderbook.Order `json:"order"`
}

func (j *JSONRPCServer) GetOrder(req *http.Request, args *GetOrderArgs, reply *GetOrderReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.GetOrder")
	defer span.End()

	exists, in, inTick, out, outTick, remaining, owner, err := j.c.GetOrderFromState(ctx, args.OrderID)
	if err != nil {
		return err
	}
	if !exists {
		return ErrOrderNotFound
	}
	reply.Order = &orderbook.Order{
		ID:        args.OrderID,
		Owner:     codec.MustAddressBech32(consts.HRP, owner),
		InAsset:   in,
		InTick:    inTick,
		OutAsset:  out,
		OutTick:   outTick,
		Remaining: remaining,
	}
	return nil
}

type LoanArgs struct {
	Destination ids.ID `json:"destination"`
	Asset       ids.ID `json:"asset"`
}

type LoanReply struct {
	Amount uint64 `json:"amount"`
}

func (j *JSONRPCServer) Loan(req *http.Request, args *LoanArgs, reply *LoanReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Loan")
	defer span.End()

	amount, err := j.c.GetLoanFromState(ctx, args.Asset, args.Destination)
	if err != nil {
		return err
	}
	reply.Amount = amount
	return nil
}

type ProjectArgs struct {
	Project ids.ID `json:"project"`
}

type ProjectReply struct {
	ID                 []byte `json:"ID"`
	ProjectName        []byte `json:"name"`
	ProjectDescription []byte `json:"description"`
	ProjectOwner       []byte `json:"owner"`
	Logo               []byte `json:"logo"`
}

func (j *JSONRPCServer) Project(req *http.Request, args *ProjectArgs, reply *ProjectReply) error {
	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Project")
	defer span.End()

	exists, project, err := j.c.GetProjectFromState(ctx, args.Project)

	if err != nil {
		return err
	}
	if !exists {
		return ErrProjectNotFound
	}
	reply.ID = []byte(project.Key)
	reply.ProjectName = project.ProjectName
	reply.ProjectDescription = project.ProjectDescription
	reply.ProjectOwner = project.ProjectOwner
	reply.Logo = project.Logo
	return err

}

type UpdateArgs struct {
	Update ids.ID `json:"update"`
}

type UpdateReply struct {
	ID                   []byte `json:"ID"`
	ProjectTxID          []byte `json:"project_id"` // reference to Project
	UpdateExecutableHash []byte `json:"executable_hash"`
	UpdateIPFSUrl        []byte `json:"executable_ipfs_url"`
	ForDeviceName        []byte `json:"for_device_name"`
	UpdateVersion        uint8  `json:"version"`
	SuccessCount         uint8  `json:"success_count"`
}

func (j *JSONRPCServer) Update(req *http.Request, args *UpdateArgs, reply *UpdateReply) error {

	ctx, span := j.c.Tracer().Start(req.Context(), "Server.Update")
	defer span.End()

	exists, update, err := j.c.GetUpdateFromState(ctx, args.Update)

	if err != nil {
		return err
	}
	if !exists {
		return ErrUpdateNotFound
	}

	reply.ID = []byte(update.Key)
	reply.ProjectTxID = []byte(update.ProjectTxID)
	reply.UpdateExecutableHash = []byte(update.UpdateExecutableHash)
	reply.UpdateIPFSUrl = []byte(update.UpdateIPFSUrl)
	reply.ForDeviceName = []byte(update.ForDeviceName)
	reply.UpdateVersion = uint8(update.UpdateVersion)
	reply.SuccessCount = uint8(update.SuccessCount)

	return err

}

type RegisterMachineCIDArgs struct {
	MachineCIDID ids.ID `json:"MachineCID"`
}

type RegisterMachineCIDReply struct {
	ID         []byte `json:"ID"`
	MachineCID []byte `json:"machine_cid"`
}

func (j *JSONRPCServer) MachineCID(req *http.Request, args *RegisterMachineCIDArgs, reply *RegisterMachineCIDReply) error {

	ctx, span := j.c.Tracer().Start(req.Context(), "Server.MachineCID")
	defer span.End()

	exists, update, err := j.c.GetMachineCID(ctx, args.MachineCIDID)

	if err != nil {
		return err
	}
	if !exists {
		return ErrMachineCIDNotFound
	}

	reply.ID = []byte(update.Key)
	reply.MachineCID = []byte(update.MachineCID)

	return err

}

type AttestMachineArgs struct {
	Tx ids.ID `json:"Tx"`
}

type AttestMachineReply struct {
	ID                  []byte `json:"ID"`
	MachineAddress      []byte `json:"machine_address"`
	MachineCategory     []byte `json:"machine_category"`
	MachineManufacturer []byte `json:"machine_manufacturer"`
	MachineCID          []byte `json:"machine_cid"`
}

func (j *JSONRPCServer) AttestMachine(req *http.Request, args *AttestMachineArgs, reply *AttestMachineReply) error {

	ctx, span := j.c.Tracer().Start(req.Context(), "Server.AttestMachine")
	defer span.End()

	exists, attestmachine, err := j.c.GetAttestMachine(ctx, args.Tx)

	if err != nil {
		return err
	}
	if !exists {
		return ErrAttestMachineNotFound
	}

	reply.ID = []byte(attestmachine.Key)
	reply.MachineAddress = []byte(attestmachine.MachineAddress)
	reply.MachineCategory = []byte(attestmachine.MachineCategory)
	reply.MachineManufacturer = []byte(attestmachine.MachineManufacturer)
	reply.MachineCID = []byte(attestmachine.MachineCID)

	return err

}

type NotarizeDataArgs struct {
	Tx ids.ID `json:"Tx"`
}

type NotarizeDataReply struct {
	ID              []byte `json:"ID"`
	AttestMachineTx []byte `json:"attest_machine_tx"`
	DataOwnerAddr   []byte `json:"data_owner_address"`
	DataCID         []byte `json:"data_cid"`
	DataType        []byte `json:"data_type"`
}

func (j *JSONRPCServer) NotarizeData(req *http.Request, args *NotarizeDataArgs, reply *NotarizeDataReply) error {

	ctx, span := j.c.Tracer().Start(req.Context(), "Server.NotarizeData")
	defer span.End()

	exists, notarizeddata, err := j.c.GetNotarizeData(ctx, args.Tx)

	if err != nil {
		return err
	}
	if !exists {
		return ErrNotarizedDataNotFound
	}

	reply.ID = []byte(notarizeddata.Key)
	reply.AttestMachineTx = []byte(notarizeddata.AttestMachineTx)
	reply.DataOwnerAddr = []byte(notarizeddata.DataOwnerAddr)
	reply.DataCID = []byte(notarizeddata.DataCID)
	reply.DataType = []byte(notarizeddata.DataType)

	return err

}
