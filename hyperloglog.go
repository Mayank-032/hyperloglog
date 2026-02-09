package main

import "math"

type IHyperloglog interface {
	GetCardinality() float64
	AddElement(str string)
	ComputeCardinality() float64
}

const (
	BITSET_CAPACITY = 64
	BIAS_CONSTANT   = 0.79402
)

type hyperloglog struct {
	cardinality float64
	registers   []int
	nBits       int
	// register datastructure
}

func NewHyperloglog(nBits int) IHyperloglog {
	return &hyperloglog{
		nBits:     nBits,
		registers: make([]int, nBits),
	}
}

func (hll *hyperloglog) GetCardinality() float64 {
	return hll.cardinality
}

func (hll *hyperloglog) AddElement(str string) {
	// firstly compute the binary value of current string
	var stringBytes = computeBinary(str)

	var index = 0
	var powOfTwo = 1
	var firstNBits = stringBytes[0:hll.nBits]
	for i := 1; i < hll.nBits; i++ {
		powOfTwo = powOfTwo * 2
		index += int(firstNBits[i]) * powOfTwo
	}

	var remainingBits = stringBytes[hll.nBits:]
	var position = positionOfLeftmostOne(remainingBits)
	hll.registers[index] = max(hll.registers[index], position)
}

func (hll *hyperloglog) ComputeCardinality() float64 {
	var sum = 0.0
	for _, register := range hll.registers {
		sum += math.Pow(2, float64((-1)*register))
	}

	var totalRegisters = len(hll.registers)
	var harmonicMean = float64(totalRegisters) * (float64(totalRegisters) / sum)
	var result = BIAS_CONSTANT * harmonicMean

	hll.cardinality = result
	return result
}
