package main

func main() {
	// 	// Decode command-line argument as hex string
	// 	data, _ := hex.DecodeString(os.Args[1])

	// 	// Parse the decoded data to retrieve height and index
	// 	buf := bytes.NewBuffer(data)
	// 	var height, index int64
	// 	err := binary.Read(buf, binary.BigEndian, &height)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	err = binary.Read(buf, binary.BigEndian, &index)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Create a new DAConfig instance
	// 	daConfig, err := NewDAConfig("e8e5f679bf7116cb")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	// Retrieve data from the DAConfig
	// 	blockData, err := daConfig.GetData(height, index)
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	fmt.Printf("Block data at height %v and index %v: %x\n", height, index, blockData)
	// }

	// type DAConfig struct {
	// 	DB          *badger.DB
	// 	NamespaceId []byte // Namespace ID as []byte
	// }

	// func NewDAConfig(namespaceIdStr string) (*DAConfig, error) {
	// 	var namespaceId []byte

	// 	if namespaceIdStr == "" {
	// 		return nil, fmt.Errorf("namespace ID cannot be empty")
	// 	}

	// 	n, err := hex.DecodeString(namespaceIdStr)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	namespaceId = n

	// 	// Open or create the BadgerDB database
	// 	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	// 	if err != nil {
	// 		return nil, err
	// 	}

	// 	return &DAConfig{
	// 		DB:          db,
	// 		NamespaceId: namespaceId,
	// 	}, nil
	// }

	// // Retrieve data based on height and index
	// func (da *DAConfig) GetData(height int64, index int64) ([]byte, error) {
	// 	var blockData []byte
	// 	err := da.DB.View(func(tx *badger.Txn) error {
	// 		key := constructKey(height, index)
	// 		item, err := tx.Get(key)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return item.Value(func(val []byte) error {
	// 			blockData = append(blockData, val...)
	// 			return nil
	// 		})
	// 	})
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	return blockData, nil
	// }

	// // PutData writes data to the BadgerDB database
	// func (da *DAConfig) PutData(height int64, index int64, data []byte) error {
	// 	// Define a function to handle the update logic
	// 	updateFunc := func(tx *badger.Txn) error {
	// 		// Store the combined data in the database
	// 		return tx.Set(constructKey(height, index), data)
	// 	}

	// 	// Open a transaction to update the database
	// 	return da.DB.Update(updateFunc)
	// }

	// // Construct a key using height and index
	//
	//	func constructKey(height int64, index int64) []byte {
	//		key := make([]byte, 16)
	//		binary.BigEndian.PutUint64(key[:8], uint64(height))
	//		binary.BigEndian.PutUint64(key[8:], uint64(index))
	//		return key
}
