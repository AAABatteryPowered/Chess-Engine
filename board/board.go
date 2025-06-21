package board

import "fmt"

type BoardMethods interface {
	FromFen(string)
}

type Board struct {
	WKings   Bitboard
	WQueens  Bitboard
	WRooks   Bitboard
	WBishops Bitboard
	WKnights Bitboard
	WPawns   Bitboard

	BKings   Bitboard
	BQueens  Bitboard
	BRooks   Bitboard
	BBishops Bitboard
	BKnights Bitboard
	BPawns   Bitboard
}

type Bitboard uint64

func (b *Bitboard) Set(pos int) {
	*b |= 1 << pos
}

func (b *Bitboard) Clear(pos int) {
	*b &^= 1 << pos
}

func (b *Bitboard) Toggle(pos int) {
	*b ^= 1 << pos
}

func (b Bitboard) IsSet(pos int) bool {
	return (b>>pos)&1 == 1
}

func (b *Board) FromFen(s string) {
	posPointer := 0
	for _, ch := range s {
		if ch >= '0' && ch <= '9' {
			num := int(ch - '0')
			posPointer += num
		}
		switch ch {
		case 'K':
			b.WKings.Set(posPointer)
			posPointer += 1
		case 'Q':
			b.WQueens.Set(posPointer)
			posPointer += 1
		case 'R':
			b.WRooks.Set(posPointer)
			posPointer += 1
		case 'B':
			b.WBishops.Set(posPointer)
			posPointer += 1
		case 'N':
			b.WKnights.Set(posPointer)
			posPointer += 1
		case 'P':
			b.WPawns.Set(posPointer)
			posPointer += 1

		case 'k':
			b.BKings.Set(posPointer)
			posPointer += 1
		case 'q':
			b.BQueens.Set(posPointer)
			posPointer += 1
		case 'r':
			b.BRooks.Set(posPointer)
			posPointer += 1
		case 'b':
			b.BBishops.Set(posPointer)
			posPointer += 1
		case 'n':
			b.BKnights.Set(posPointer)
			posPointer += 1
		case 'p':
			b.BPawns.Set(posPointer)
			posPointer += 1
		}
		if ch == 'K' {

		}
	}
}

func (b *Board) AllBitboards() []Bitboard {
	return []Bitboard{
		b.WKings,
		b.WQueens,
		b.WRooks,
		b.WBishops,
		b.WKnights,
		b.WPawns,

		b.BKings,
		b.BQueens,
		b.BRooks,
		b.BBishops,
		b.BKnights,
		b.BPawns,
	}
}

func (b *Board) DebugPrint() {
	var finalstr string
	newlinesadded := 0
	allbbs := b.AllBitboards()
	for pos := 0; pos < 64; pos++ {
		piecefound := false
		for i, bb := range allbbs {
			if bb.IsSet(pos) {
				switch i {
				case 0:
					finalstr += "K "
				case 1:
					finalstr += "Q "
				case 2:
					finalstr += "R "
				case 3:
					finalstr += "B "
				case 4:
					finalstr += "N "
				case 5:
					finalstr += "P "
				case 6:
					finalstr += "k "
				case 7:
					finalstr += "q "
				case 8:
					finalstr += "r "
				case 9:
					finalstr += "b "
				case 10:
					finalstr += "n "
				case 11:
					finalstr += "p "
				}
				//fmt.Println(len(finalstr))

				piecefound = true
				break
			}
		}

		if !piecefound {
			finalstr += ". "
		}

		if (len(finalstr)-newlinesadded)%16 == 0 {
			fmt.Println(len(finalstr))
			newlinesadded += 1
			finalstr += "\n"
		}
	}
	fmt.Println(finalstr)
}
