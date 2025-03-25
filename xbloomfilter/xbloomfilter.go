package xbloomfilter

import (
	"math"
	"unsafe"
)

type BloomFilterHasherType int

const (
	MurmurHasher BloomFilterHasherType = iota
)

type BloomFilter[T any] struct {
	bitsArray           []uint64          // array containing the bits of the bloom filter
	bitArrayElementSize int               // number of bits in each element of the array
	nBits               int               // number of bits required to achieve the desirable false positive ratio based on the number of expected items
	nHashes             int               // number of hashes required to achieve the desirable false positive ratio based on the number of expected items
	hasher              BloomFilterHasher // hash module to use
}

func NewBloomFilter[T any](expectedItems int, acceptableFalsePositiveRatio float64, bloomFilterHasherType BloomFilterHasherType) *BloomFilter[T] {
	requiredBits := calculateRequiredNumberOfBits(expectedItems, acceptableFalsePositiveRatio)
	nHashes := calculateRequiredNumberOfHashes(requiredBits, expectedItems)
	bitArrayElementSize := int(unsafe.Sizeof(uint64(0)))
	bitsArray := make([]uint64, calculateSizeOfBitArray(requiredBits, bitArrayElementSize))
	bloomFilter := BloomFilter[T]{
		bitsArray:           bitsArray,
		bitArrayElementSize: bitArrayElementSize,
		nBits:               requiredBits,
		nHashes:             nHashes,
		hasher:              getHasher(bloomFilterHasherType),
	}
	return &bloomFilter
}

func calculateRequiredNumberOfBits(expectedItems int, acceptableFalsePositiveRatio float64) int {
	return int(math.Ceil(-float64(expectedItems)*math.Log(acceptableFalsePositiveRatio)) / (math.Pow(math.Log(2), 2)))
}

func calculateRequiredNumberOfHashes(requiredBits int, expectedItems int) int {
	return int(math.Ceil(float64(requiredBits) / float64(expectedItems) * math.Log(2)))
}

func calculateSizeOfBitArray(requiredBits int, bitArrayElementSize int) int {
	return int(math.Ceil(float64(requiredBits / bitArrayElementSize)))
}

func (bloomFilter *BloomFilter[T]) Add(object T) error {
	hashes, err := bloomFilter.hasher.GenerateHashes(object, bloomFilter.nHashes)
	if err != nil {
		return err
	}

	for i := range bloomFilter.nHashes {
		bitToSet := hashes[i] % uint64(bloomFilter.nBits)
		bitArrayIndex, bitArrayIndexOffset := bloomFilter.getBitArrayIndexAndIndexOffset(bitToSet)
		bloomFilter.bitsArray[bitArrayIndex] = bloomFilter.bitsArray[bitArrayIndex] | (1 << bitArrayIndexOffset) // set the relevant bit
	}

	return nil
}

func (bloomFilter *BloomFilter[T]) MightContain(object T) (bool, error) {
	hashes, err := bloomFilter.hasher.GenerateHashes(object, bloomFilter.nHashes)
	if err != nil {
		return false, err
	}

	for i := range bloomFilter.nHashes {
		bitToCheck := hashes[i] % uint64(bloomFilter.nBits)
		bitArrayIndex, bitArrayIndexOffset := bloomFilter.getBitArrayIndexAndIndexOffset(bitToCheck)
		b := bloomFilter.bitsArray[bitArrayIndex] & (1 << bitArrayIndexOffset) // extract only the relevant bit
		if b != (1 << bitArrayIndexOffset) {                                   // check if the relevant bit is set
			return false, nil
		}
	}
	return true, nil
}

func (bloomFilter *BloomFilter[T]) getBitArrayIndexAndIndexOffset(bitToSet uint64) (uint64, uint64) {
	bitArrayIndex := bitToSet / uint64(bloomFilter.bitArrayElementSize)
	bitArrayIndexOffset := bitToSet % uint64(bloomFilter.bitArrayElementSize)
	return bitArrayIndex, bitArrayIndexOffset
}

func getHasher(bloomFilterHasherType BloomFilterHasherType) BloomFilterHasher {
	switch bloomFilterHasherType {
	case MurmurHasher:
		return NewMurmurHash()
	default:
		panic("Invalid BloomFilterHasherType")
	}
}
