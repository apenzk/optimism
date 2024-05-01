package rollup

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/log"
)

type DAConfig struct {
	Log         log.Logger
	DB          *bolt.DB
	NamespaceId *NamespaceId // Namespace ID as []byte
}

type NamespaceId struct {
	Id []byte
}

// error: cannot use daCfg.Namespace (variable of type rollup.Namespace) as cnc.Namespace value in argument to daCfg.Client.NamespacedData

func NewDAConfig(namespaceIdStr string) (*DAConfig, error) {
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

	fmt.Printf("... creating DB path.\n")
	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Construct the database file path using the current directory
	dbPath := filepath.Join(currentDir, "mydatabase.db")
	fmt.Printf("... created DB path: %s\n", dbPath)

	fmt.Printf("... creating or opening DB")
	// Open or create the BoltDB database
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}
	fmt.Printf("... created or opening DB")

	_ = db
	return &DAConfig{
		DB:          db,
		NamespaceId: namespaceId,
	}, nil
}

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
	updateFunc := func(tx *bolt.Tx) error {
		// Retrieve or create the bucket for the namespace
		bucket, err := tx.CreateBucketIfNotExists(da.NamespaceId.Id)
		if err != nil {
			return err
		}
		key := constructKey(height, index)

		// print height, index and key for debugging
		fmt.Printf("print 2:    height: %d, index: %d, key: %v\n", height, index, key)

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
