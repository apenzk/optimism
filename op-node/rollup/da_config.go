package rollup

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dgraph-io/badger/v4"
	"github.com/ethereum/go-ethereum/log"
)

type DAConfig struct {
	Log         log.Logger
	DB          *badger.DB
	NamespaceId *NamespaceId // Namespace ID as []byte
}

type NamespaceId struct {
	Id []byte
}

func NewDAConfig(namespaceIdStr string) (*DAConfig, error) {
	var nid [8]byte

	n, err := hex.DecodeString(namespaceIdStr)
	if err != nil {
		return &DAConfig{}, err
	}
	copy(nid[:], n)
	if err != nil {
		return &DAConfig{}, err
	}

	nid2 := append([]byte(nil), nid[:]...)
	namespaceId := &NamespaceId{
		Id: nid2,
	}

	fmt.Printf("... creating DB path.\n")
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// Construct the database directory path using the current directory
	dbPath := filepath.Join(currentDir, "mydatabase")
	// Check if the directory exists, create it if not
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		fmt.Printf("DB path does not exist. Creating directory: %s\n", dbPath)
		err := os.MkdirAll(dbPath, 0755)
		if err != nil {
			fmt.Printf("Failed to create directory: %v\n", err)
		}
	}
	fmt.Printf("... created DB path: %s\n", dbPath)

	fmt.Printf("... creating or opening DB...\n")
	// Open or create the BadgerDB database
	// opts := badger.DefaultOptions(dbPath)
	db, err := badger.Open(badger.DefaultOptions(dbPath).WithVerifyValueChecksum(true))
	if err != nil {
		fmt.Printf("... error: creating or opening DB...\n")
		return nil, err
	}
	fmt.Printf("... created or opened DB\n")

	fmt.Printf("... returning DAConfig\n")
	return &DAConfig{
		DB:          db,
		NamespaceId: namespaceId,
	}, nil
}

// Retrieve data based on height and index
func (da *DAConfig) GetData(height int64, index int64) ([]byte, error) {
	var blockData []byte
	err := da.DB.View(func(tx *badger.Txn) error {
		key := constructKey(height, index)
		item, err := tx.Get(key)
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			blockData = append(blockData, val...)
			return nil
		})
	})
	if err != nil {
		return nil, err
	}
	return blockData, nil
}

// PutData writes data to the BadgerDB database
func (da *DAConfig) PutData(height int64, index int64, data []byte) error {
	fmt.Println("...... PutData ......")
	key := constructKey(height, index)
	fmt.Printf("... print 1: height: %d, index: %d, key: %v\n", height, index, key)
	fmt.Printf("... da.NamespaceId.Id: %v\n", da.NamespaceId.Id)
	// da.Log.Info("... log 1: height: %d, index: %d, key: %v\n", height, index, key)

	// da.NamespaceId.Id = e8e5f679bf7116cb
	if da.NamespaceId.Id == nil {
		fmt.Printf("... error: PutData: namespace ID is nil\n")
		return fmt.Errorf("namespace ID is nil")
	}

	// Define a function to handle the update logic
	updateFunc := func(tx *badger.Txn) error {
		// Store the combined data in the database
		return tx.Set(key, data)
	}

	// Open a transaction to update the database
	err := da.DB.Update(updateFunc)
	return err
}

// Construct a key using height and index
func constructKey(height int64, index int64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, uint64(index))

	key := append(heightBytes, indexBytes...)
	return key
}
