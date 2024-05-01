package main

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/boltdb/bolt"
)

func main() {
	// Decode command-line argument as hex string
	data, _ := hex.DecodeString(os.Args[1])

	// Parse the decoded data to retrieve height and index
	buf := bytes.NewBuffer(data)
	var height, index int64
	err := binary.Read(buf, binary.BigEndian, &height)
	if err != nil {
		panic(err)
	}
	err = binary.Read(buf, binary.BigEndian, &index)
	if err != nil {
		panic(err)
	}
	fmt.Printf("celestia block height: %v; tx index: %v\n", height, index)
	fmt.Println("-----------------------------------------")

	// Create a new DAConfig instance
	daConfig, err := NewDAConfig("http://localhost:26659", "e8e5f679bf7116cb")
	if err != nil {
		panic(err)
	}

	// Retrieve data from the DAConfig
	blockData, err := daConfig.GetData(height, index)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Block data at height %v and index %v: %x\n", height, index, blockData)
}

// -----------------------------------------------------

// have to later clean this and import rollup

// -----------------------------------------------------

type DAConfig struct {
	DB          *bolt.DB
	NamespaceId *NamespaceId // Namespace ID as []byte
}

type NamespaceId struct {
	Id []byte
}

// error: cannot use daCfg.Namespace (variable of type rollup.Namespace) as cnc.Namespace value in argument to daCfg.Client.NamespacedData

func NewDAConfig(rpc string, namespaceIdStr string) (*DAConfig, error) {
	var nid [8]byte

	n, err := hex.DecodeString(namespaceIdStr)
	if err != nil {
		return &DAConfig{}, err
	}
	copy(nid[:], n)
	// daClient, err := cnc.NewClient(rpc, cnc.WithTimeout(30*time.Second))
	if err != nil {
		return &DAConfig{}, err
	}

	nid2 := append([]byte(nil), nid[:]...)
	namespaceId := &NamespaceId{
		Id: nid2,
	}

	return &DAConfig{
		DB:          nil,
		NamespaceId: namespaceId,
	}, nil
}

// func NewDAConfig(namespaceIdStr string) (*DAConfig, error) {
// 	// Create a BoltDB database in memory
// 	db, err := bolt.Open(":memory:", 0600, nil)
// 	if err != nil {
// 		return nil, fmt.Errorf("error creating BoltDB database: %v", err)
// 	}
// 	namespaceId, err := hex.DecodeString(namespaceIdStr)
// 	if err != nil {
// 		db.Close()
// 		return nil, fmt.Errorf("error decoding namespace ID: %v", err)
// 	}

// 	return &DAConfig{
// 		DB:          db,
// 		NamespaceId: namespaceId,
// 	}, nil
// }

// Retrieve data based on height and index
func (da *DAConfig) GetData(height int64, index int64) ([]byte, error) {
	var blockData []byte
	err := da.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(da.NamespaceId.Id) // Use the namespace as the bucket name
		if bucket == nil {
			return bolt.ErrBucketNotFound
		}
		key := constructKey(height, index)
		blockData = bucket.Get(key)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blockData, nil
}

// PutData writes data to the BoltDB database
func (da *DAConfig) PutData(height int64, index int64, data []byte) error {
	// Define a function to handle the update logic
	updateFunc := func(tx *bolt.Tx) error {
		// Retrieve or create the bucket for the namespace
		bucket, err := tx.CreateBucketIfNotExists(da.NamespaceId.Id)
		if err != nil {
			return err
		}
		key := constructKey(height, index)

		// // Convert the height to bytes to use as the key
		// key := heightToBytes(height)
		// // Combine the blocks into a single byte slice
		// data := combineBlocks(blocks)

		// Store the combined data in the bucket
		return bucket.Put(key, data)
	}

	// Open a transaction to update the database
	err := da.DB.Update(updateFunc)
	return err
}

// Construct a key using height and index
func constructKey(height int64, index int64) []byte {
	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(height))

	// indexBytes := make([]byte, 4)
	// binary.BigEndian.PutUint64(heightBytes, uint64(height))
	indexBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(indexBytes, uint64(index))

	key := append(heightBytes, indexBytes...)
	return key
}
