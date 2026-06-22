package board

// Dedicated lookup slices per square to store precalculated attack masks
var rookTable [64][]uint64
var bishopTable [64][]uint64

var rookMasks [64]uint64
var bishopMasks [64]uint64

var rookMaskBitCount [64]int
var bishopMaskBitCount [64]int

func init() {
	initializeSlidingMasks()
	populateLookupTables()
}

// GetRookAttacks compresses occupancy bits into a perfect sequential index with zero branches
func GetRookAttacks(sq Square, occupancy uint64) uint64 {
	occMasked := occupancy & rookMasks[sq]
	idx := occupancyToIndex(occMasked, rookMasks[sq])
	return rookTable[sq][idx]
}

// GetBishopAttacks compresses occupancy bits into a perfect sequential index with zero branches
func GetBishopAttacks(sq Square, occupancy uint64) uint64 {
	occMasked := occupancy & bishopMasks[sq]
	idx := occupancyToIndex(occMasked, bishopMasks[sq])
	return bishopTable[sq][idx]
}

// Compresses scattered occupancy bits down into a contiguous index integer
func occupancyToIndex(occupancy uint64, mask uint64) int {
	idx := 0
	bitIndex := 0
	for mask != 0 {
		lsb := mask & -mask // Isolate the lowest set bit of the mask
		mask ^= lsb         // Clear it
		
		if (occupancy & lsb) != 0 {
			idx |= (1 << bitIndex)
		}
		bitIndex++
	}
	return idx
}

// Maps a contiguous index back into the scattered bits of the target occupancy mask
func indexToOccupancy(index int, bits int, mask uint64) uint64 {
	var occupancy uint64
	for i := 0; i < bits; i++ {
		lsb := mask & -mask // Isolate the lowest set bit of the mask
		mask ^= lsb         // Clear it
		
		if (index & (1 << i)) != 0 {
			occupancy |= lsb
		}
	}
	return occupancy
}

func generateRookAttacksOnTheFly(sq Square, occupancy uint64) uint64 {
	var attacks uint64
	r, f := int(sq>>3), int(sq&7)
	for i := r + 1; i < 8; i++ { target := Square((i << 3) | f); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	for i := r - 1; i >= 0; i-- { target := Square((i << 3) | f); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	for i := f + 1; i < 8; i++ { target := Square((r << 3) | i); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	for i := f - 1; i >= 0; i-- { target := Square((r << 3) | i); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	return attacks
}

func generateBishopAttacksOnTheFly(sq Square, occupancy uint64) uint64 {
	var attacks uint64
	r, f := int(sq>>3), int(sq&7)
	for i, j := r+1, f+1; i < 8 && j < 8; i, j = i+1, j+1 { target := Square((i << 3) | j); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	for i, j := r+1, f-1; i < 8 && j >= 0; i, j = i+1, j-1 { target := Square((i << 3) | j); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	for i, j := r-1, f+1; i >= 0 && j < 8; i, j = i-1, j+1 { target := Square((i << 3) | j); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	for i, j := r-1, f-1; i >= 0 && j >= 0; i, j = i-1, j-1 { target := Square((i << 3) | j); attacks |= uint64(1) << target; if (occupancy & (uint64(1) << target)) != 0 { break } }
	return attacks
}

func initializeSlidingMasks() {
	for sq := 0; sq < 64; sq++ {
		r, f := sq>>3, sq&7
		var rMask, bMask uint64
		
		for i := r + 1; i < 7; i++ { rMask |= uint64(1) << Square((i << 3) | f) }
		for i := r - 1; i > 0; i-- { rMask |= uint64(1) << Square((i << 3) | f) }
		for i := f + 1; i < 7; i++ { rMask |= uint64(1) << Square((r << 3) | i) }
		for i := f - 1; i > 0; i-- { rMask |= uint64(1) << Square((r << 3) | i) }
		rookMasks[sq] = rMask
		rookMaskBitCount[sq] = countBits(rMask)

		for i, j := r+1, f+1; i < 7 && j < 7; i, j = i+1, j+1 { bMask |= uint64(1) << Square((i << 3) | j) }
		for i, j := r+1, f-1; i < 7 && j > 0; i, j = i+1, j-1 { bMask |= uint64(1) << Square((i << 3) | j) }
		for i, j := r-1, f+1; i > 0 && j < 7; i, j = i-1, j+1 { bMask |= uint64(1) << Square((i << 3) | j) }
		for i, j := r-1, f-1; i > 0 && j > 0; i, j = i-1, j-1 { bMask |= uint64(1) << Square((i << 3) | j) }
		bishopMasks[sq] = bMask
		bishopMaskBitCount[sq] = countBits(bMask)
	}
}

func countBits(n uint64) int {
	count := 0
	for n != 0 {
		count += int(n & 1)
		n >>= 1
	}
	return count
}

func populateLookupTables() {
	for sq := 0; sq < 64; sq++ {
		// Allocate size matching exactly 2^(number of mask bits) entries
		sizeR := 1 << rookMaskBitCount[sq]
		rookTable[sq] = make([]uint64, sizeR)
		for idx := 0; idx < sizeR; idx++ {
			subOccupancy := indexToOccupancy(idx, rookMaskBitCount[sq], rookMasks[sq])
			rookTable[sq][idx] = generateRookAttacksOnTheFly(Square(sq), subOccupancy)
		}

		sizeB := 1 << bishopMaskBitCount[sq]
		bishopTable[sq] = make([]uint64, sizeB)
		for idx := 0; idx < sizeB; idx++ {
			subOccupancy := indexToOccupancy(idx, bishopMaskBitCount[sq], bishopMasks[sq])
			bishopTable[sq][idx] = generateBishopAttacksOnTheFly(Square(sq), subOccupancy)
		}
	}
}