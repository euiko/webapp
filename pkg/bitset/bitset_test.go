package bitset

import "testing"

func TestBitSet(t *testing.T) {
	const (
		testBitIndex0 int64 = iota
		testBitIndex1
		testBitIndex2
		testBitIndex3
		testBitIndex4
		testBitIndex75  int64 = 75
		testBitIndex76  int64 = 76
		testBitIndex512 int64 = 512
	)

	allIndexes := []int64{
		testBitIndex0,
		testBitIndex1,
		testBitIndex2,
		testBitIndex3,
		testBitIndex4,
		testBitIndex75,
		testBitIndex76,
		testBitIndex512,
	}

	deletedIndexes := []int64{
		testBitIndex3,
		testBitIndex76,
	}

	remainingIndexes := []int64{
		testBitIndex0,
		testBitIndex1,
		testBitIndex2,
		testBitIndex4,
		testBitIndex75,
		testBitIndex512,
	}

	var bitset BitSet
	bitset.Add(allIndexes...)

	if !checkBitSet(bitset, allIndexes...) {
		t.Error("BitSet should have all indexes")
	}

	bitset.Delete(deletedIndexes...)

	if !checkBitSet(bitset, remainingIndexes...) {
		t.Error("BitSet should have all indexes except deleted indexes")
	}

	if checkBitSet(bitset, testBitIndex3) {
		t.Error("BitSet should not testBitIndex3")
	}

	if !checkBitSet(bitset, testBitIndex75) {
		t.Error("BitSet should not testBitIndex3")
	}

	if checkBitSet(bitset, testBitIndex76) {
		t.Error("BitSet should not testBitIndex3")
	}
}

func checkBitSet(b BitSet, indexes ...int64) bool {
	for _, index := range indexes {
		if !b.Has(index) {
			return false
		}
	}

	return true
}
