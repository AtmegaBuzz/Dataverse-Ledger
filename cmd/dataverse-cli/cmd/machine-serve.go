package cmd

import (
	"context"
	"dataverse/actions"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ipfs/go-cid"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/multiformats/go-multibase"
	mh "github.com/multiformats/go-multihash"
)

type RegisteredMachine struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	MachineCID string `gorm:"unique"`
	Txid       string
}

type AttestedMachine struct {
	ID             uint   `gorm:"primaryKey;autoIncrement"`
	MachineAddress string `gorm:"unique"`
	Txid           string
}

type NotarizeData struct {
	ID      uint `gorm:"primaryKey;autoIncrement"`
	Owner   string
	DataCid string
	Txid    string `gorm:"unique"`
}

var DB *gorm.DB

func RegisterMachineCID(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		_, _, factory, cli, scli, tcli, _ := handler.DefaultActor()

		machineCID := r.URL.Query().Get("machinecid")

		var registerMachine RegisteredMachine
		machinecid := r.URL.Query().Get("machinecid")
		result := DB.Where("machine_c_id = ?", machinecid)

		fmt.Println(result.RowsAffected > 0, registerMachine.MachineCID, machineCID)

		if result.Error != nil {
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Machine Already exists", http.StatusBadRequest)
			return
		}

		project := &actions.RegisterMachine{
			MachineCID: []byte(machineCID),
		}

		// Generate transaction
		te, tx, _ := sendAndWait(ctx, nil, project, cli, scli, tcli, factory, true)

		if !te {
			http.Error(w, "Tx failed", http.StatusBadRequest)
			return
		}

		registerMahcine := RegisteredMachine{MachineCID: machineCID, Txid: tx.String()}
		DB.Create(&registerMahcine)
		DB.Commit()

		response := map[string]interface{}{
			"MachineRegisterTx": trimNullChars(tx.String()),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

func GetregisterMachineCID(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var registerMachine RegisteredMachine
		machinecid := r.URL.Query().Get("machinecid")
		result := DB.Where("machine_c_id = ?", machinecid).First(&registerMachine)
		if result.Error != nil {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "Machine not registered", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(registerMachine.Txid)
	}

}

type AttestedMachineInfo struct {
	MachineAddress      string `json:"machine_address"`
	MachineManufacturer string `json:"machine_manufacturer"`
	MachineCID          string `json:"machine_cid"`
	MachineCategory     string `json:"machine_category"`
}

func AttestMachine(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		_, _, factory, cli, scli, tcli, _ := handler.DefaultActor()

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var attestMachine AttestedMachineInfo
		if err := json.Unmarshal(body, &attestMachine); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		// check if tx exists onchain
		var exists bool
		var registerMachine RegisteredMachine
		_ = DB.Model(registerMachine).Select("count(*) > 0").Where("machine_c_id = ?", attestMachine.MachineCID).Find(&exists).Error

		if !exists {
			http.Error(w, "Machine not registered", http.StatusBadRequest)
			return
		}

		project := &actions.AttestMachine{
			MachineAddress:      []byte(attestMachine.MachineAddress),
			MachineCategory:     []byte(attestMachine.MachineCategory),
			MachineManufacturer: []byte(attestMachine.MachineManufacturer),
			MachineCID:          []byte(attestMachine.MachineCID),
		}

		// Generate transaction
		_, id, err := sendAndWait(ctx, nil, project, cli, scli, tcli, factory, true)

		response := ""
		if err != nil {
			http.Error(w, "Error while Attestation", http.StatusInternalServerError)
		}

		registerMahcine := AttestedMachine{MachineAddress: string(project.MachineAddress), Txid: id.String()}
		DB.Create(&registerMahcine)
		DB.Commit()

		response = id.String()

		w.Header().Set("Content-Type", "application/json")

		// w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

type NotarizeDataArgs struct {
	Owner   string `json:"machine_address"`
	DataCid string `json:"data_cid"`
	Data    string `json:"data"`
}

func NotarizeDataView(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		_, _, factory, cli, scli, tcli, _ := handler.DefaultActor()

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var attestMachine NotarizeDataArgs
		if err := json.Unmarshal(body, &attestMachine); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		// check if tx exists onchain
		// var exists bool
		var attestedMachineDB AttestedMachine
		// _ = DB.Model(attestedMachineDB).Select("count(*) > 0").Where("machine_address = ?", attestMachine.Owner).Find(&exists).Error
		// fmt.Println(exists)
		// if !exists {
		// 	http.Error(w, "Machine not registered", http.StatusBadRequest)
		// 	return
		// }

		_ = DB.First(&attestedMachineDB, "machine_address = ?", attestMachine.Owner)

		fmt.Println(attestedMachineDB.Txid, attestMachine.Owner)

		notarizedata := &actions.NotarizeData{
			MachineAttestTx: []byte("plmnt1qh05ghszrxfh8taksfqhvyfgewleq6u5ru9xlg"),
			DataCID:         []byte(attestMachine.DataCid),
			DataType:        []byte("/dataverse.asset.MsgNotarizedAsset"),
			DataOwnerAddr:   []byte(attestMachine.Owner),
		}

		// Generate transaction
		_, id, err := sendAndWait(ctx, nil, notarizedata, cli, scli, tcli, factory, true)

		response := ""
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Error while Attestation", http.StatusInternalServerError)
			return
		}

		registerMahcine := NotarizeData{DataCid: attestMachine.DataCid, Owner: attestMachine.Owner, Txid: id.String()}
		DB.Create(&registerMahcine)
		DB.Commit()

		response = id.String()

		w.Header().Set("Content-Type", "application/json")

		// w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

type VerifyNotarizeDataArgs struct {
	Owner string `json:"machine_address"`
	Data  string `json:"data"`
}

func VerifyNotarizeDataView(ctx context.Context) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		body, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		var verifyArgs VerifyNotarizeDataArgs
		if err := json.Unmarshal(body, &verifyArgs); err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}

		bytes := []byte(verifyArgs.Data)
		hash, _ := mh.Sum(bytes, mh.SHA2_256, -1)

		c := cid.NewCidV1(cid.Raw, hash)
		encodedCID, _ := c.StringOfBase(multibase.Base32)

		fmt.Println(encodedCID)

		var exists bool
		var notarizeData NotarizeData
		_ = DB.Model(&notarizeData).Select("count(*) > 0").Where("data_cid = ?", encodedCID).Find(&exists).Error

		response := exists

		w.Header().Set("Content-Type", "application/json")

		// w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)

	}

}

var serverDataverseCmd = &cobra.Command{
	Use: "serveDataverse",
	RunE: func(*cobra.Command, []string) error {

		ctx := context.Background()

		db, err := gorm.Open(sqlite.Open("index.db"), &gorm.Config{})
		if err != nil {
			fmt.Println("failed to connect database:", err)
		}
		DB = db
		DB.AutoMigrate(&RegisteredMachine{})
		DB.AutoMigrate(&AttestedMachine{})
		DB.AutoMigrate(&NotarizeData{})

		http.HandleFunc("/register-machine", RegisterMachineCID(ctx)) //machinecid
		http.HandleFunc("/get-register-machine", GetregisterMachineCID(ctx))
		http.HandleFunc("/attest-machine", AttestMachine(ctx))
		http.HandleFunc("/notarize-data", NotarizeDataView(ctx))
		http.HandleFunc("/verify", VerifyNotarizeDataView(ctx))

		// Start the HTTP server on port 8080
		fmt.Println("Server is listening on port 8080...")
		// err_http := http.ListenAndServe(":8080", nil)
		// fmt.Println("Server Ended")

		// if err_http != nil {
		// 	return err_http
		// }

		return http.ListenAndServe(":8080", nil)
	},
}
