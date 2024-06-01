package storage

type ProjectData struct {
	Key                string `json:"key"`
	ProjectName        []byte `json:"name"`
	ProjectDescription []byte `json:"description"`
	ProjectOwner       []byte `json:"owner"`
	Logo               []byte `json:"url"`
}

type UpdateData struct {
	Key                  string `json:"key"`
	ProjectTxID          []byte `json:"project_id"` // reference to Project
	UpdateExecutableHash []byte `json:"executable_hash"`
	UpdateIPFSUrl        []byte `json:"executable_ipfs_url"`
	ForDeviceName        []byte `json:"for_device_name"`
	UpdateVersion        uint8  `json:"version"`
	SuccessCount         uint8  `json:"success_count"`
}

type RegisterMachineCIDData struct {
	Key        string `json:"key"`
	MachineCID []byte `json:"machine_cid"`
}

type AttestMachineData struct {
	Key                 string `json:"key"`
	MachineAddress      []byte `json:"machine_address"`
	MachineCategory     []byte `json:"machine_category"`
	MachineManufacturer []byte `json:"machine_manufacturer"`
	MachineCID          []byte `json:"machine_cid"`
}

type NotarizeDataData struct {
	Key             string `json:"key"`
	AttestMachineTx []byte `json:"attest_machine_tx"`
	DataOwnerAddr   []byte `json:"data_owner_address"`
	DataCID         []byte `json:"data_cid"`
	DataType        []byte `json:"data_type"`
}
