package rollup

import (
	"bytes"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/log"
	_ "github.com/mattn/go-sqlite3"
)

type DAConfig struct {
	Log         log.Logger
	DB          *sql.DB
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

	nid2 := append([]byte(nil), nid[:]...)
	namespaceId := &NamespaceId{
		Id: nid2,
	}

	dbPath := filepath.Join("/tmp", "mydatabase.db")

	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return &DAConfig{}, err
	}

	// Create the necessary table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS data (
		key1 INTEGER,
		key2 INTEGER,
		actualdata BLOB,
		PRIMARY KEY (key1, key2)
	);`)
	if err != nil {
		return nil, err
	}

	return &DAConfig{
		DB:          db,
		NamespaceId: namespaceId,
	}, nil
}

// Retrieve data based on height and index
func (da *DAConfig) GetData(height int64, index int64) ([]byte, error) {
	var blockData []byte
	fmt.Printf("... GetData: height: %d, index: %d\n", height, index)
	row := da.DB.QueryRow("SELECT actualdata FROM data WHERE key1 = ? AND key2 = ?", height, index)
	fmt.Printf("... GetData: SELECT actualdata FROM data WHERE key1 = %d AND key2 = %d\n", height, index)

	err := row.Scan(&blockData)
	if err != nil {
		fmt.Printf("Error in GetData: %v\n", err)
		return nil, err
	}
	fmt.Printf("... GetData: data: %v\n", blockData)

	return blockData, nil
}

// PutData writes data to the SQLite database
func (da *DAConfig) PutData(height int64, index int64, data []byte) error {
	fmt.Printf("... SetData: height: %d, index: %d\n", height, index)
	_, err := da.DB.Exec("INSERT OR REPLACE INTO data (key1, key2, actualdata) VALUES (?, ?, ?)", height, index, data)
	if err != nil {
		fmt.Printf("Error in PutData: %v\n", err)
		return err
	}

	// Query the database to check if the data was inserted correctly
	fmt.Printf("... GetData after PutData: height: %d, index: %d\n", height, index)
	var blockData []byte
	row := da.DB.QueryRow("SELECT actualdata FROM data WHERE key1 = ? AND key2 = ?", height, index)
	fmt.Printf("... GetData: SELECT actualdata FROM data WHERE key1 = %d AND key2 = %d\n", height, index)

	err = row.Scan(&blockData)
	if err != nil {
		fmt.Printf("Error in GetData after PutData: %v\n", err)
		return err
	}

	// Compare the inserted data with the retrieved data
	fmt.Printf("... GetData after PutData: data: %v\n", blockData)
	if !bytes.Equal(data, blockData) {
		fmt.Printf("... Data mismatch after PutData: expected %v, got %v\n", data, blockData)
		return errors.New("data mismatch after PutData")
	}
	fmt.Printf("... Data inserted successfully\n")

	return nil
}
