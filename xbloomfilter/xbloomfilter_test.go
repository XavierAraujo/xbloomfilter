package xbloomfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestObject struct {
	ID string
}

func TestNonExistentItem(t *testing.T) {
	bloomFilter := NewBloomFilter[TestObject](1000, 0.01, MurmurHasher)
	object := TestObject{"1"}
	objectInFilter, err := bloomFilter.MightContain(object)
	assert.Nil(t, err)
	assert.False(t, objectInFilter)
}

func TestExistentItem(t *testing.T) {
	bloomFilter := NewBloomFilter[TestObject](1000, 0.01, MurmurHasher)
	object := TestObject{"1"}
	err := bloomFilter.Add(object)
	objectInFilter, err := bloomFilter.MightContain(object)
	assert.Nil(t, err)
	assert.True(t, objectInFilter)
}
