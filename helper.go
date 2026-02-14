package main

import (
	"hash/fnv"
	"math/bits"
)

// calculateHash handles both string and int64
func calculateHash(s string) uint64 {
	var h = fnv.New64a()
	h.Write([]byte(s))

	var sum = h.Sum64()
	h.Reset()

	return sum
}

func positionOfLeftmostOne(stream uint64) int {
	var numberOfLeadingZeros = bits.LeadingZeros64(stream)
	return numberOfLeadingZeros + 1
}
