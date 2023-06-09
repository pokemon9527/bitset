package bitset

import "fmt"

type Bitset struct {
	b        []byte
	n        int
	endIndex int
	endMask  byte
}

// NewBitset ...
func NewBitset(n int) *Bitset {
	endIndex, endOffset := n>>3, n&7
	endMask := ^byte(255 >> byte(endOffset))
	if endOffset == 0 {
		endIndex = -1
	}
	return &Bitset{make([]byte, (n+7)>>3), n, endIndex, endMask}
}

// NewBitsetFromBytes Creates a new bitset from a given byte stream. Returns nil if the
// data is invalid in some way.
func NewBitsetFromBytes(n int, data []byte) *Bitset {
	bitset := NewBitset(n)
	if len(bitset.b) != len(data) {
		return nil
	}
	copy(bitset.b, data)
	if bitset.endIndex >= 0 && bitset.b[bitset.endIndex]&(^bitset.endMask) != 0 {
		return nil
	}
	return bitset
}

// Set ...
func (b *Bitset) Set(index int) {
	b.checkRange(index)
	b.b[index>>3] |= byte(128 >> byte(index&7))
}

// Clear ...
func (b *Bitset) Clear(index int) {
	b.checkRange(index)
	b.b[index>>3] &= ^byte(128 >> byte(index&7))
}

// IsSet ...
func (b *Bitset) IsSet(index int) bool {
	b.checkRange(index)
	return (b.b[index>>3] & byte(128>>byte(index&7))) != 0
}

// Len ...
func (b *Bitset) Len() int {
	return b.n
}

// InRange ...
func (b *Bitset) InRange(index int) bool {
	return 0 <= index && index < b.n
}

func (b *Bitset) checkRange(index int) {
	if !b.InRange(index) {
		panic(any(fmt.Sprintf("Index %d out of range 0..%d.", index, b.n)))
	}
}

func (b *Bitset) AndNot(b2 *Bitset) {
	if b.n != b2.n {
		panic(any(fmt.Sprintf("Unequal bitset sizes %d != %d", b.n, b2.n)))
	}
	for i := 0; i < len(b.b); i++ {
		b.b[i] = b.b[i] & ^b2.b[i]
	}
	b.clearEnd()
}

func (b *Bitset) clearEnd() {
	if b.endIndex >= 0 {
		b.b[b.endIndex] &= b.endMask
	}
}

// IsEndValid ...
func (b *Bitset) IsEndValid() bool {
	if b.endIndex >= 0 {
		return (b.b[b.endIndex] & b.endMask) == 0
	}
	return true
}

// FindNextSet ...
// TODO: Make this fast
func (b *Bitset) FindNextSet(index int) int {
	for i := index; i < b.n; i++ {
		if (b.b[i>>3] & byte(128>>byte(i&7))) != 0 {
			return i
		}
	}
	return -1
}

// FindNextClear ...
// TODO: Make this fast
func (b *Bitset) FindNextClear(index int) int {
	for i := index; i < b.n; i++ {
		if (b.b[i>>3] & byte(128>>byte(i&7))) == 0 {
			return i
		}
	}
	return -1
}

// Bytes ...
func (b *Bitset) Bytes() []byte {
	return b.b
}
