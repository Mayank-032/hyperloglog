package main

import (
	"math"
	"sync"
)

type IHyperloglog interface {
	GetCardinality() int
	AddElem(str string)
	ComputeCardinality() float64
}

// hash capacity
const (
	BITSET_CAPACITY = 64
)

// BIAS CONSTANTS - depending on the number of registers
const (
	BIAS_CONSTANT_64BIT   = 0.709
	BIAS_CONSTANT_32BIT   = 0.697
	BIAS_CONSTANT_16BIT   = 0.673
	BIAS_CONSTANT_DEFAULT = 0.7213
)

type hyperloglog struct {
	cardinality float64
	registers   []int
	nBits       int
	mu          sync.RWMutex
}

func NewHyperLogLog(nBits int) *hyperloglog {
	if nBits < 4 || nBits > 16 {
		return nil
	}

	var registerSize = int(math.Pow(2, float64(nBits)))
	return &hyperloglog{
		nBits:     nBits,
		registers: make([]int, registerSize),
	}
}

func (hll *hyperloglog) GetCardinality() int {
	return int(math.Floor(hll.cardinality))
}

func (hll *hyperloglog) AddElem(str string) {
	// firstly compute the binary value of current string
	var hashedValue = calculateHash(str)
	var firstNBits = hashedValue >> (BITSET_CAPACITY - uint64(hll.nBits))
	var remainingBits = hashedValue << uint64(hll.nBits)

	var leftmostOneBitPosition = positionOfLeftmostOne(remainingBits)

	hll.mu.Lock()
	hll.registers[firstNBits] = max(hll.registers[firstNBits], leftmostOneBitPosition)
	hll.mu.Unlock()

	_ = hll.ComputeCardinality()
}

func (hll *hyperloglog) ComputeCardinality() float64 {
	var sum = 0.0

	hll.mu.RLock()
	for _, register := range hll.registers {
		sum += math.Pow(2, float64((-1)*register))
	}
	hll.mu.RUnlock()

	var totalRegisters = len(hll.registers)
	var harmonicMean = float64(totalRegisters) * (float64(totalRegisters) / sum)
	var result float64

	switch totalRegisters {
	case 16:
		result = BIAS_CONSTANT_16BIT * harmonicMean
	case 32:
		result = BIAS_CONSTANT_32BIT * harmonicMean
	case 64:
		result = BIAS_CONSTANT_64BIT * harmonicMean
	default:
		var biasConstant = BIAS_CONSTANT_DEFAULT / (1 + (1.079 / float64(totalRegisters)))
		result = biasConstant * harmonicMean
	}

	hll.cardinality = result
	return result
}
