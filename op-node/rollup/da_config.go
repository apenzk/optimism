package rollup

import (
	"encoding/hex"
	"time"

	"github.com/celestiaorg/go-cnc"
)

type DAConfig struct {
	Rpc string
	// NamespaceId [8]byte
	Namespace cnc.Namespace
	Client    *cnc.Client
}

// error: cannot use daCfg.Namespace (variable of type rollup.Namespace) as cnc.Namespace value in argument to daCfg.Client.NamespacedData

func NewDAConfig(rpc string, namespaceId string) (*DAConfig, error) {
	var nid [8]byte

	n, err := hex.DecodeString(namespaceId)
	if err != nil {
		return &DAConfig{}, err
	}
	copy(nid[:], n)
	daClient, err := cnc.NewClient(rpc, cnc.WithTimeout(30*time.Second))
	if err != nil {
		return &DAConfig{}, err
	}

	namespace := cnc.Namespace{
		// ID: nid,
		ID: append([]byte(nil), nid[:]...),
	}

	return &DAConfig{
		Namespace: namespace,
		Rpc:       rpc,
		Client:    daClient,
	}, nil
}
