package rollup

import (
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
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

	// Get the current directory
	currentDir, err := os.Getwd()
	if err != nil {
		return &DAConfig{}, err
	}
	// Construct the database file path using the current directory
	dbPath := filepath.Join(currentDir, "mydatabase.db")

	// Open the SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return &DAConfig{}, err
	}

	// Create the necessary table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS data (
		key1 INTEGER,
		key2 INTEGER,
		data BLOB,
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
	row := da.DB.QueryRow("SELECT data FROM data WHERE key1 = ? AND key2 = ?", height, index)
	fmt.Printf("... GetData: SELECT data FROM data WHERE key1 = %d AND key2 = %d\n", height, index)
	fmt.Printf("... GetData: row: %v\n", row)

	err := row.Scan(&blockData)
	if err != nil {
		fmt.Printf("Error in GetData: %v\n", err)
		return nil, err
	}
	return blockData, nil
}

// PutData writes data to the SQLite database
func (da *DAConfig) PutData(height int64, index int64, data []byte) error {
	fmt.Printf("... SetData: height: %d, index: %d\n", height, index)
	_, err := da.DB.Exec("INSERT OR REPLACE INTO data (key1, key2, data) VALUES (?, ?, ?)", height, index, data)
	return err
}
