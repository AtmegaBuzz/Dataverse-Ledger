// Copyright (C) 2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package actions

var (
	OutputValueZero              = []byte("value is zero")
	OutputMemoTooLarge           = []byte("memo is too large")
	OutputAssetIsNative          = []byte("cannot mint native asset")
	OutputAssetAlreadyExists     = []byte("asset already exists")
	OutputAssetMissing           = []byte("asset missing")
	OutputInTickZero             = []byte("in rate is zero")
	OutputOutTickZero            = []byte("out rate is zero")
	OutputSupplyZero             = []byte("supply is zero")
	OutputSupplyMisaligned       = []byte("supply is misaligned")
	OutputOrderMissing           = []byte("order is missing")
	OutputUnauthorized           = []byte("unauthorized")
	OutputWrongIn                = []byte("wrong in asset")
	OutputWrongOut               = []byte("wrong out asset")
	OutputWrongOwner             = []byte("wrong owner")
	OutputInsufficientInput      = []byte("insufficient input")
	OutputInsufficientOutput     = []byte("insufficient output")
	OutputValueMisaligned        = []byte("value is misaligned")
	OutputSymbolEmpty            = []byte("symbol is empty")
	OutputSymbolIncorrect        = []byte("symbol is incorrect")
	OutputSymbolTooLarge         = []byte("symbol is too large")
	OutputDecimalsIncorrect      = []byte("decimal is incorrect")
	OutputDecimalsTooLarge       = []byte("decimal is too large")
	OutputMetadataEmpty          = []byte("metadata is empty")
	OutputMetadataTooLarge       = []byte("metadata is too large")
	OutputSameInOut              = []byte("same asset used for in and out")
	OutputConflictingAsset       = []byte("warp has same asset as another")
	OutputAnycast                = []byte("anycast output")
	OutputNotWarpAsset           = []byte("not warp asset")
	OutputWarpAsset              = []byte("warp asset")
	OutputWrongDestination       = []byte("wrong destination")
	OutputMustFill               = []byte("must fill request")
	OutputWarpVerificationFailed = []byte("warp verification failed")
	OutputInvalidDestination     = []byte("invalid destination")

	OutputProjectCreated             = []byte("Project created")
	OutputProjectDescriptionNotGiven = []byte("Project Description not provided")
	OutputProjectNameNotGiven        = []byte("Project Name not provided")
	OutputProjectInvalidOwner        = []byte("Project Owner Invalid format")

	OutputProjectTxIdNotProvided          = []byte("Project Txid not provided")
	OutputUpdateExecutableHashNotProvided = []byte("Update Executable Hash not provided")
	OutputUpdateExecutableIPFSNotProvided = []byte("Update Executable IPFS url Not Provided")
	OutputForDeviceNameNotProvided        = []byte("Update Device Name Not Provided")
	OutputUpdateVersionNotProvided        = []byte("Update Version Not Provided")

	OutputRegisterMachineNotProvided = []byte("Machine CID Not Provided")
)
