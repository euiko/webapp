package bitset

import (
	"fmt"
	"strconv"
	"strings"
)

type (
	BitSet struct {
		values []int64
	}
)

func (b *BitSet) ToString(delimiters ...string) string {
	delimiter := ","
	if len(delimiters) > 0 {
		delimiter = delimiters[0]
	}

	strValues := make([]string, len(b.values))
	for i, value := range b.values {
		strValues[i] = fmt.Sprintf("%d", value)
	}

	return strings.Join(strValues, delimiter)
}

func (b *BitSet) FromString(str string, delimiters ...string) error {
	delimiter := ","
	if len(delimiters) > 0 {
		delimiter = delimiters[0]
	}

	values := strings.Split(str, delimiter)
	for _, value := range values {
		v, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}

		b.Add(v)
	}

	return nil
}

func (b *BitSet) Add(indexes ...int64) {
	for _, index := range indexes {
		value := b.getValue(index)
		offset := indexToOffsetValue(index)
		b.setValue(index, value|offset)
	}
}

func (b *BitSet) Delete(indexes ...int64) {
	for _, index := range indexes {
		value := b.getValue(index)
		offset := indexToOffsetValue(index)
		b.setValue(index, value&^offset)
	}
}

func (b *BitSet) Has(index int64) bool {
	value := b.getValue(index)
	offset := indexToOffsetValue(index)
	return value&offset == offset
}

func (b *BitSet) getValue(index int64) int64 {
	valueIndex := indexToValueIndex(index)
	if valueIndex >= len(b.values) {
		return 0
	}

	return b.values[valueIndex]
}

func (b *BitSet) setValue(index int64, value int64) {
	valueIndex := indexToValueIndex(index)
	valuesSize := len(b.values)
	for i := valueIndex; i >= valuesSize; i-- {
		b.values = append(b.values, 0)
	}

	b.values[valueIndex] = value
}

func indexToValueIndex(index int64) int {
	return int(index / 64)
}

func indexToOffsetValue(index int64) int64 {
	index = index % 64
	return 1 << index
}
