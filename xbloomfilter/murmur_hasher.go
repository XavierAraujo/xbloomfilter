package xbloomfilter

import (
	"bytes"
	"encoding/gob"

	"github.com/spaolacci/murmur3"
)

type MurMurHash struct {
	hasher  murmur3.Hash128
}

func NewMurmurHash() *MurMurHash {
	return &MurMurHash{
		hasher:  murmur3.New128WithSeed(0),
	}
}

func (murmurHash MurMurHash) GenerateHashes(object any, nHashes int) ([]uint64, error) {
	hashes := make([]uint64, nHashes)
	h1, h2, err := hash(object, murmurHash.hasher)
	if err != nil {
		return nil, err
	}
	for i := range nHashes {
		// Instead of computing multiple hashes with distinct seeds
		// for the same object which can be terribly slow we do this
		// trick to generate multiple hash values from a hash function
		// call
		hashes[i] = (h1 + (uint64(i) * h2))
	}

	return hashes, nil
}

func hash(object any, hasher murmur3.Hash128) (uint64, uint64, error) {
	data, err := toBytes(object)
	if err != nil {
		return 0, 0, err
	}

	hasher.Write(data)
	h1, h2 := hasher.Sum128()
	hasher.Reset()

	return h1, h2, nil
}

func toBytes(object any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(object)
	return buf.Bytes(), err
}
