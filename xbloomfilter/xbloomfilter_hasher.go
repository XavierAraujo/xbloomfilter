package xbloomfilter

type BloomFilterHasher interface {
	/*
	 * Should generate 'nHashes' distinct hashes for a given object.
	 * The multiple generated hashes should always be the same
	 * for the same object.
	 */
	GenerateHashes(object any, nHashes int) ([]uint64, error)
}
