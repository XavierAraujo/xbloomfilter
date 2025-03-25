package xbloomfilter

import "testing"

type TestObject struct {
	ID string
}

func TestNonExistentItem(t *testing.T) {
	bloomFilter := NewBloomFilter[TestObject](1000, 0.01, MurmurHasher)
	object := TestObject{"1"}
	objectInFilter, err := bloomFilter.MightContain(object)
	if err != nil {
		t.Error("Error occurred calling ProbablyContains(): ", err)
	}
	if objectInFilter {
		t.Error("Object ", object, " contained in filter")
	}
}

func TestExistentItem(t *testing.T) {
	bloomFilter := NewBloomFilter[TestObject](1000, 0.01, MurmurHasher)
	object := TestObject{"1"}
	err := bloomFilter.Add(object)
	if err != nil {
		t.Error("Error occurred calling BloomFilter.Add(): ", err)
	}
	objectInFilter, err := bloomFilter.MightContain(object)
	if err != nil {
		t.Error("Error occurred calling BloomFilter.ProbablyContains(): ", err)
	}
	if !objectInFilter {
		t.Error("Object ", object, " contained in filter")
	}
}
