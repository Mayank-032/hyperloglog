package main

import (
	"fmt"
	"hash/fnv"
)

func calculateHash(s string) uint64 {
	h := fnv.New64a()

	h.Write([]byte(s))
	return h.Sum64()
}

func computeBinary(s string) string {
	var hash = calculateHash(s)

	return fmt.Sprintf("%064b", hash)
}

func positionOfLeftmostOne(stringBytes string) int {
	for index, byte := range stringBytes {
		if byte == 1 {
			return index
		}
	}

	return len(stringBytes)
}
