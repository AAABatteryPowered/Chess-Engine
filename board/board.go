package board

import (
	"fmt"
	"math/bits"
	"slices"
	"strconv"
	"strings"
)

const (
	FileA  Bitboard = 0x0101010101010101
	FileH  Bitboard = 0x8080808080808080
	FileAB Bitboard = 0x303030303030303
	FileGH Bitboard = 0xC0C0C0C0C0C0C0C0

	WPawnStartRank Bitboard = 0x00FF000000000000
	BPawnStartRank Bitboard = 0x000000000000FF00
)

type BoardMethods interface {
	FromFen(string)
	GenMoves() []Move
	SetTurn(bool)
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

	WCastleQ bool
	WCastleK bool
	BCastleQ bool
	BCastleK bool

	EnPassantTarget int

	HalfMoves int
	FullMoves int

	FilledSquares Bitboard
	Turn          bool
}

type preCompTables struct {
	King   [64]Bitboard
	Knight [64]Bitboard
	Rook   [64]Bitboard
}

type Move struct {
	From      int
	To        int
	Castle    int
	EnPassant bool
	Promotion int
}

type Bitboard uint64

func algebraicToSquare(notation string) int {
	if len(notation) != 2 {
		return -1
	}

	file := int(notation[0] - 'a')
	rank := int(notation[1] - '1')

	return rank*8 + file
}

func (b *Bitboard) Set(pos int) {
	*b |= 1 << pos
}

func (b *Bitboard) Clear(pos int) {
	*b &^= 1 << pos
}

func (b *Bitboard) Toggle(pos int) {
	*b ^= 1 << pos
}

func (b *Board) Pieces() (Bitboard, Bitboard) {
	if !b.Turn {
		return b.BKings | b.BQueens | b.BRooks | b.BBishops | b.BKnights | b.BPawns, b.WKings | b.WQueens | b.WRooks | b.WBishops | b.WKnights | b.WPawns
	} else {
		return b.WKings | b.WQueens | b.WRooks | b.WBishops | b.WKnights | b.WPawns, b.BKings | b.BQueens | b.BRooks | b.BBishops | b.BKnights | b.BPawns
	}
}

func (b Bitboard) IsSet(pos int) bool {
	return (b>>pos)&1 == 1
}

func (b *Board) SetTurn(t bool) {
	b.Turn = t
}

func notationToSquare(notation string) int {
	if len(notation) != 2 {
		return -1
	}

	file := notation[0] - 'a'
	rank := notation[1] - '1'

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return -1
	}

	return int(rank)*8 + int(file)
}

func (b *Board) FromFen(s string) {
	posPointer := 0
	var subdata string
outerLoop:
	for xx, ch := range s {
		if ch >= '0' && ch <= '9' {
			num := int(ch - '0')
			posPointer += num
		}
		switch ch {
		case 'K':
			b.WKings.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'Q':
			b.WQueens.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'R':
			b.WRooks.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'B':
			b.WBishops.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'N':
			b.WKnights.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'P':
			b.WPawns.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'k':
			b.BKings.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'q':
			b.BQueens.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'r':
			b.BRooks.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'b':
			b.BBishops.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'n':
			b.BKnights.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case 'p':
			b.BPawns.Set(posPointer)
			b.FilledSquares.Set(posPointer)
			posPointer += 1
		case ' ':
			subdata = s[xx+1 : len(s)]
			break outerLoop
		}
	}
	//w QKqk e3 41 20

	for i, field := range strings.Fields(subdata) {
		if field != "-" {
			for _, runee := range field {
				switch i + 1 {
				case 1:
					if field == "w" {
						b.Turn = true
					} else {
						b.Turn = false
					}
				case 2:
					if runee == 'Q' {
						b.WCastleQ = true
					}
					if runee == 'K' {
						b.WCastleK = true
					}
					if runee == 'q' {
						b.BCastleQ = true
					}
					if runee == 'k' {
						b.BCastleK = true
					}
				case 3:
					b.EnPassantTarget = notationToSquare(field)
				case 4:
					integer, err := strconv.Atoi(field)
					if err != nil {
						fmt.Println(err)
					}
					b.HalfMoves = integer
				case 5:
					integer, err := strconv.Atoi(field)
					if err != nil {
						fmt.Println(err)
					}
					b.FullMoves = integer
				}
			}
		}
	}

	fmt.Println(b.Turn, b.WCastleQ, b.WCastleK, b.EnPassantTarget, b.HalfMoves, b.FullMoves)
}

func (b *Board) Copy() *Board {
	newBoard := &Board{}
	*newBoard = *b

	// newBoard.pieces = make([]Piece, len(b.pieces))
	// copy(newBoard.pieces, b.pieces)

	return newBoard
}

func (b *Board) PieceAt(square int) int {
	allbb := b.AllBitboards()
	for index, board := range allbb {
		if board.IsSet(square) {
			return index
		}
	}
	return -1
}

func (b Bitboard) ToSquares() []int {
	var squares []int
	bb := uint64(b)
	for bb != 0 {
		square := bits.TrailingZeros64(bb)
		squares = append(squares, square)
		bb &= bb - 1
	}
	return squares
}

func (b *Board) AllBitboards() []*Bitboard {
	return []*Bitboard{
		&b.WKings,
		&b.WQueens,
		&b.WRooks,
		&b.WBishops,
		&b.WKnights,
		&b.WPawns,

		&b.BKings,
		&b.BQueens,
		&b.BRooks,
		&b.BBishops,
		&b.BKnights,
		&b.BPawns,
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
				piecefound = true
				break
			}
		}

		if !piecefound {
			finalstr += ". "
		}

		if (len(finalstr)-newlinesadded)%16 == 0 {
			newlinesadded += 1
			finalstr += "\n"
		}
	}
	fmt.Println(finalstr)
}

func (b Bitboard) DebugPrint() {
	for rank := 0; rank < 8; rank++ {
		for file := 0; file < 8; file++ {
			square := rank*8 + file
			if (b>>square)&1 == 1 {
				fmt.Print("1 ")
			} else {
				fmt.Print("0 ")
			}
		}
		fmt.Println()
	}
}

/*func (b *Board) FollowRay(turn bool, piecetype int, direction int, recurringmoves *[]Move) *[]Move {
	var raymoves []Move = *recurringmoves
	if recurringmoves == nil {
		raymoves = make([]Move, 1)
	}

	if turn {
		if piecetype == 1 {
			if direction == 1 {
				for (b.WRooks << 8 &^ b.FilledSquares) > 0 {
					b.WRooks = b.WRooks << 8 &^ b.FilledSquares
					append(raymoves)
					b.FollowRay(turn, piecetype, direction, &raymoves)
				}
			}
		}
	}

	return &raymoves
}*/

func (b *Board) IsSquareAttacked(square int) bool {
	nextturnboard := *b
	nextturnboard.SetTurn(!b.Turn)
	allMoves := nextturnboard.GenMoves()

	for _, move := range allMoves {
		if move.To == square {
			return true
		}
	}

	return false
}

func lineAttack(from, to int, occ Bitboard, dirs []int) bool {
	for _, dir := range dirs {
		sq := from
		for {
			sq += dir
			if sq < 0 || sq > 63 {
				break
			}
			// Prevent rank wrap for horizontal moves
			if (dir == 1 || dir == -1) && abs((sq%8)-((sq-dir)%8)) != 1 {
				break
			}
			if sq == to {
				return true
			}
			if Bitboard(1<<sq)&occ != 0 {
				break
			}
		}
	}
	return false
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (b *Board) IsKingAttacked() bool {
	var kingSq int
	var occ = b.FilledSquares
	var attackers Bitboard

	if b.Turn {
		// White to move → check if white king is attacked by black
		kingSq = bits.TrailingZeros64(uint64(b.WKings))

		// Pawn attacks (black → down)
		attackers |= (((b.BPawns >> 7) &^ FileH) | ((b.BPawns >> 9) &^ FileA)) & (1 << kingSq)

		// Knight attacks
		for _, sq := range b.BKnights.ToSquares() {
			if precomped.Knight[sq].IsSet(kingSq) {
				return true
			}
		}

		// King attacks
		for _, sq := range b.BKings.ToSquares() {
			if precomped.King[sq].IsSet(kingSq) {
				return true
			}
		}

		// Sliding pieces
		for _, sq := range b.BBishops.ToSquares() {
			if lineAttack(sq, kingSq, occ, []int{7, 9, -7, -9}) {
				return true
			}
		}
		for _, sq := range b.BRooks.ToSquares() {
			if lineAttack(sq, kingSq, occ, []int{8, -8, 1, -1}) {
				return true
			}
		}
		for _, sq := range b.BQueens.ToSquares() {
			if lineAttack(sq, kingSq, occ, []int{7, 9, -7, -9, 8, -8, 1, -1}) {
				return true
			}
		}
	} else {
		// Black to move → check if black king is attacked by white
		kingSq = bits.TrailingZeros64(uint64(b.BKings))

		// Pawn attacks (white → up)
		attackers |= (((b.WPawns << 7) &^ FileA) | ((b.WPawns << 9) &^ FileH)) & (1 << kingSq)

		// Knight attacks
		for _, sq := range b.WKnights.ToSquares() {
			if precomped.Knight[sq].IsSet(kingSq) {
				return true
			}
		}

		// King attacks
		for _, sq := range b.WKings.ToSquares() {
			if precomped.King[sq].IsSet(kingSq) {
				return true
			}
		}

		// Sliding pieces
		for _, sq := range b.WBishops.ToSquares() {
			if lineAttack(sq, kingSq, occ, []int{7, 9, -7, -9}) {
				return true
			}
		}
		for _, sq := range b.WRooks.ToSquares() {
			if lineAttack(sq, kingSq, occ, []int{8, -8, 1, -1}) {
				return true
			}
		}
		for _, sq := range b.WQueens.ToSquares() {
			if lineAttack(sq, kingSq, occ, []int{7, 9, -7, -9, 8, -8, 1, -1}) {
				return true
			}
		}
	}

	// If pawn or other attacks exist on the king’s square
	return attackers != 0
}

/*
func IsKingAttacked(b *Board) (bool, []Move) {
	nextturnboard := *b
	nextturnboard.SetTurn(!b.Turn)
	allMoves := nextturnboard.GenMoves()

	var kingboard Bitboard

	if b.Turn {
		kingboard = b.WKings
	} else {
		kingboard = b.BKings
	}

	for _, move := range allMoves {
		tosquares := kingboard.ToSquares()
		if move.To == tosquares[0] {
			return true, allMoves
		}
	}

	return false, allMoves
}*/

func (b *Board) Moves() []Move {
	incheck := b.IsKingAttacked()
	ourmoves := b.GenMoves()
	var filteredmoves []Move
	if incheck {
		for _, ourmove := range ourmoves {
			var copyboard *Board = &Board{}
			*copyboard = *b
			copyboard.PlayMove(ourmove)
			stillincheck := copyboard.IsKingAttacked()
			if !stillincheck {
				filteredmoves = append(filteredmoves, ourmove)
			}
		}
		return filteredmoves
	}
	return ourmoves
}

/*func (b *Board) Moves() []Move {
	ourmoves := b.GenMoves()
	var filteredmoves []Move

	for _, ourmove := range ourmoves {
		var copyboard *Board = &Board{}
		*copyboard = *b
		copyboard.PlayMove(ourmove)
		stillincheck, _ := IsKingAttacked(copyboard)
		if !stillincheck {
			filteredmoves = append(filteredmoves, ourmove)
		}
	}

	return filteredmoves
}*/

func (b *Board) PlayMove(move Move) {
	movingpiece := b.PieceAt(move.From)
	targetpiece := b.PieceAt(move.To)
	b.FilledSquares.Clear(move.From)
	castling := move.Castle
	allbb := b.AllBitboards()
	if castling != 0 {
		switch castling {
		case 1:
			//no clear move.froms cuz we do that 6 lines above
			b.FilledSquares.Clear(move.To)
			b.FilledSquares.Set(2)
			b.FilledSquares.Set(3)
			allbb[0].Clear(move.From) // white king
			allbb[0].Set(2)
			allbb[2].Clear(move.To) // white rook
			allbb[2].Set(3)
		case 2:
			b.FilledSquares.Clear(move.To)
			b.FilledSquares.Set(6)
			b.FilledSquares.Set(5)
			allbb[0].Clear(move.From) // white king
			allbb[0].Set(6)
			allbb[2].Clear(move.To) // white rook
			allbb[2].Set(5)
		case 3:
			b.FilledSquares.Clear(move.To)
			b.FilledSquares.Set(58)
			b.FilledSquares.Set(59)
			allbb[6].Clear(move.From)
			allbb[6].Set(58)
			allbb[8].Clear(move.To)
			allbb[8].Set(59)
		case 4:
			b.FilledSquares.Clear(move.To)
			b.FilledSquares.Set(61)
			b.FilledSquares.Set(62)
			allbb[6].Clear(move.From)
			allbb[6].Set(62)
			allbb[8].Clear(move.To)
			allbb[8].Set(61)
		}
	} else if move.Promotion != 0 {
		b.FilledSquares.Set(move.To)
		allbb[movingpiece].Clear(move.From)
		if targetpiece != -1 {
			allbb[targetpiece].Clear(move.To)
		}
		if b.Turn {
			allbb[move.Promotion].Set(move.To)
		} else {
			allbb[move.Promotion+6].Set(move.To)
		}
	} else if move.EnPassant {
		b.FilledSquares.Clear(move.To)
		allbb[movingpiece].Clear(move.From)
		allbb[targetpiece].Clear(move.To)
		allbb[movingpiece].Set(move.To)
		if b.Turn {
			allbb[b.PieceAt(move.To-8)].Clear(move.To - 8)
		} else {
			allbb[b.PieceAt(move.To+8)].Clear(move.To + 8)
		}
	} else if targetpiece > 5 {
		// piece is black
		if b.Turn {
			//b.FilledSquares.Clear(move.To)
			allbb[movingpiece].Clear(move.From)
			allbb[targetpiece].Clear(move.To)
			allbb[movingpiece].Set(move.To)
		}
	} else if targetpiece == -1 {
		b.FilledSquares.Set(move.To)
		allbb[movingpiece].Clear(move.From)
		allbb[movingpiece].Set(move.To)
	} else if targetpiece > -1 && targetpiece < 6 {
		// piece is white
		if !b.Turn {
			allbb[movingpiece].Clear(move.From)
			allbb[targetpiece].Clear(move.To)
			allbb[movingpiece].Set(move.To)
		}
	}
	b.Turn = !b.Turn
}

func RookDepth(startsquare int, depth int) []Bitboard {
	var b Bitboard = 0
	b.Set(startsquare)

	var bbs []Bitboard = make([]Bitboard, 0)

	north := b << (8 * depth) // &^ b.w
	south := b >> (8 * depth)
	east := (b << depth) &^ FileA
	west := (b >> depth) &^ FileH

	if depth <= 7-(startsquare%8) {
		bbs = append(bbs, east)
	} else {
		bbs = append(bbs, 0)
	}
	if depth <= (startsquare % 8) {
		bbs = append(bbs, west)
	} else {
		bbs = append(bbs, 0)
	}

	bbs = append(bbs, north)
	bbs = append(bbs, south)

	return bbs
}

func BishopDepth(startsquare int, depth int) ([]Bitboard, []bool) {
	var b Bitboard = 0
	b.Set(startsquare)

	var bbs []Bitboard = make([]Bitboard, 0)
	returningbools := make([]bool, 4)

	sw := b << (7 * depth) &^ FileH
	se := b << (9 * depth) &^ FileA
	nw := b >> (9 * depth) &^ FileH
	ne := b >> (7 * depth) &^ FileA

	bbs = append(bbs, nw)
	bbs = append(bbs, ne)
	bbs = append(bbs, sw)
	bbs = append(bbs, se)

	if nw == 0 {
		returningbools[0] = false
	} else {
		returningbools[0] = true
	}
	if ne == 0 {
		returningbools[1] = false
	} else {
		returningbools[1] = true
	}
	if sw == 0 {
		returningbools[2] = false
	} else {
		returningbools[2] = true
	}
	if se == 0 {
		returningbools[3] = false
	} else {
		returningbools[3] = true
	}

	return bbs, returningbools
}

func QueenDepth(startsquare int, depth int) ([]Bitboard, []bool) {
	var b Bitboard = 0
	b.Set(startsquare)

	var bbs []Bitboard = make([]Bitboard, 0)
	returningbools := make([]bool, 8)

	north := b << (8 * depth) // &^ b.w
	south := b >> (8 * depth)
	east := (b << depth) &^ FileA
	west := (b >> depth) &^ FileH
	nw := b << (7 * depth) &^ FileH
	ne := b << (9 * depth) &^ FileA
	sw := b >> (9 * depth) &^ FileH
	se := b >> (7 * depth) &^ FileA

	if depth <= 7-(startsquare%8) {
		bbs = append(bbs, east)
	} else {
		bbs = append(bbs, 0)
	}
	if depth <= (startsquare % 8) {
		bbs = append(bbs, west)
	} else {
		bbs = append(bbs, 0)
	}

	//e,w,nw,ne,sw,se,n,s

	bbs = append(bbs, nw)
	bbs = append(bbs, ne)
	bbs = append(bbs, sw)
	bbs = append(bbs, se)

	if nw == 0 {
		returningbools[2] = false
	} else {
		returningbools[2] = true
	}
	if ne == 0 {
		returningbools[3] = false
	} else {
		returningbools[3] = true
	}
	if sw == 0 {
		returningbools[4] = false
	} else {
		returningbools[4] = true
	}
	if se == 0 {
		returningbools[5] = false
	} else {
		returningbools[5] = true
	}

	returningbools[0] = true
	returningbools[1] = true
	returningbools[6] = true
	returningbools[7] = true

	bbs = append(bbs, north)
	bbs = append(bbs, south)

	return bbs, returningbools
}

func GeneratePrecomputedTables() *preCompTables {
	precomp := &preCompTables{}

	for i := 0; i < 64; i++ {
		b := Bitboard(1 << i)

		north := b << 8
		south := b >> 8
		east := (b << 1) &^ FileA
		west := (b >> 1) &^ FileH
		northEast := (b << 9) &^ FileA
		northWest := (b << 7) &^ FileH
		southEast := (b >> 7) &^ FileA
		southWest := (b >> 9) &^ FileH

		l1 := (b >> 1) & Bitboard(0x7f7f7f7f7f7f7f7f)
		l2 := (b >> 2) & Bitboard(0x3f3f3f3f3f3f3f3f)
		r1 := (b << 1) & Bitboard(0xfefefefefefefefe)
		r2 := (b << 2) & Bitboard(0xfcfcfcfcfcfcfcfc)
		h1 := l1 | r1
		h2 := l2 | r2

		precomp.Knight[i] = (h1 << 16) | (h1 >> 16) | (h2 << 8) | (h2 >> 8)
		precomp.King[i] = north | south | east | west | northEast | northWest | southEast | southWest
	}

	return precomp
}

var precomped *preCompTables = GeneratePrecomputedTables()

func TrailingZerosLoop(b Bitboard) []int {
	var squareslist []int
	bb := uint64(b)
	for bb != 0 {
		square := bits.TrailingZeros64(bb)
		squareslist = append(squareslist, square)

		bb &= bb - 1
	}
	return squareslist
}

func (b *Board) IsInCheck() bool {

	return false
}

func RecurringRookDepth(bb *Board, turn bool, moves *[]Move) *[]Move {
	var recurringrookmoves *[]Move = moves

	var squares []int
	if !turn {
		squares = TrailingZerosLoop(bb.BRooks)
	} else {
		squares = TrailingZerosLoop(bb.WRooks)
	}

	for _, startsquare := range squares {
		var savededgeblockeddirs []bool = make([]bool, 4)
		savededgeblockeddirs[0] = true
		savededgeblockeddirs[1] = true
		savededgeblockeddirs[2] = true
		savededgeblockeddirs[3] = true
		if startsquare < 64 {
			for n := range 7 {
				moveboards := RookDepth(startsquare, n+1)
				if len(moveboards) > 0 {
					for i := range 4 {
						if !savededgeblockeddirs[i] {
							continue
						}
						targetsquare := bits.TrailingZeros64(uint64(moveboards[i]))

						if targetsquare < 64 {
							pieceat := bb.PieceAt(targetsquare)
							if (pieceat > 5 && turn) || (pieceat < 6 && pieceat >= 0 && !turn) {
								//fmt.Println(fmt.Sprintf("%d wants to go to %d", startsquare, targetsquare))
								Move := Move{From: startsquare, To: targetsquare}
								*recurringrookmoves = append(*recurringrookmoves, Move)
								savededgeblockeddirs[i] = false
							} else if pieceat == -1 {
								Move := Move{From: startsquare, To: targetsquare}
								*recurringrookmoves = append(*recurringrookmoves, Move)
							} else {
								//fmt.Println(targetsquare)
								savededgeblockeddirs[i] = false
							}
						}
					}
				}
			}
		}
	}
	return recurringrookmoves
}

func RecurringBishopDepth(bb *Board, turn bool, moves *[]Move) *[]Move {
	var recurringbishopmoves *[]Move = moves

	var squares []int
	if !turn {
		squares = TrailingZerosLoop(bb.BBishops)
	} else {
		squares = TrailingZerosLoop(bb.WBishops)
	}

	for _, startsquare := range squares {
		var savededgeblockeddirs []bool
		if startsquare < 64 {
			for n := range 7 {
				moveboards, edgeblockeddirs := BishopDepth(startsquare, n+1)
				if savededgeblockeddirs == nil {
					savededgeblockeddirs = edgeblockeddirs
				}
				//fmt.Println(moveboards)
				for ss, v := range edgeblockeddirs {
					if !v {
						savededgeblockeddirs[ss] = false
					}
				}
				if len(moveboards) > 0 {
					for i := range 4 {
						if savededgeblockeddirs[i] == false {
							continue
						}
						targetsquare := bits.TrailingZeros64(uint64(moveboards[i]))
						if targetsquare < 64 {
							pieceat := bb.PieceAt(targetsquare)
							if (pieceat > 5 && turn) || (pieceat < 6 && pieceat >= 0 && !turn) {
								//fmt.Println(fmt.Sprintf("%d wants to go to %d", startsquare, targetsquare))
								Move := Move{From: startsquare, To: targetsquare}
								*recurringbishopmoves = append(*recurringbishopmoves, Move)
								savededgeblockeddirs[i] = false
							} else if pieceat == -1 {
								Move := Move{From: startsquare, To: targetsquare}
								*recurringbishopmoves = append(*recurringbishopmoves, Move)
							} else {
								//fmt.Println(targetsquare)
								savededgeblockeddirs[i] = false
							}
						}
					}
				}
			}
		}
	}
	return recurringbishopmoves
}

func RecurringQueenDepth(bb *Board, turn bool, moves *[]Move) *[]Move {
	var recurringqueenmoves *[]Move = moves

	var squares []int
	if !turn {
		squares = TrailingZerosLoop(bb.BQueens)
	} else {
		squares = TrailingZerosLoop(bb.WQueens)
	}

	for _, startsquare := range squares {
		var savededgeblockeddirs []bool
		if startsquare < 64 {
			for n := range 7 {
				moveboards, edgeblockeddirs := QueenDepth(startsquare, n+1)
				if savededgeblockeddirs == nil {
					savededgeblockeddirs = edgeblockeddirs
				}
				for ss, v := range edgeblockeddirs {
					if !v {
						savededgeblockeddirs[ss] = false
					}
				}

				if len(moveboards) > 0 {
					for i := range 8 {

						if savededgeblockeddirs[i] == false {
							continue
						}
						targetsquare := bits.TrailingZeros64(uint64(moveboards[i]))
						if targetsquare < 64 {
							pieceat := bb.PieceAt(targetsquare)
							if (pieceat > 5 && turn) || (pieceat < 6 && pieceat >= 0 && !turn) {
								//fmt.Println(fmt.Sprintf("%d wants to go to %d", startsquare, targetsquare))
								Move := Move{From: startsquare, To: targetsquare}
								*recurringqueenmoves = append(*recurringqueenmoves, Move)
								savededgeblockeddirs[i] = false
							} else if pieceat == -1 {
								Move := Move{From: startsquare, To: targetsquare}
								*recurringqueenmoves = append(*recurringqueenmoves, Move)
							} else {
								//fmt.Println(targetsquare)
								savededgeblockeddirs[i] = false
							}
						}
					}
				}
			}
		}
	}
	return recurringqueenmoves
}

func (b *Board) GenMoves() []Move {
	allMoves := []Move{}
	ourpieces, opponentpieces := b.Pieces()

	if b.Turn {
		if b.IsInCheck() {
			return nil
		}
		for i, bb := range []Bitboard{b.WKings, b.WQueens, b.WRooks, b.WBishops, b.WKnights, b.WPawns} {
			switch i {
			case 0:
				square := bits.TrailingZeros64(uint64(bb))
				if square < 64 {
					adjacentking := precomped.King[square] &^ ourpieces
					after := adjacentking.ToSquares()
					for _, v := range after {
						Move := Move{From: square, To: v}
						allMoves = append(allMoves, Move)
					}
				}
			case 1:
				var placeholder []Move = make([]Move, 0)
				n := RecurringQueenDepth(b, b.Turn, &placeholder)
				allMoves = slices.Concat(allMoves, *n)
			case 2:
				var placeholder []Move = make([]Move, 0)
				n := RecurringRookDepth(b, b.Turn, &placeholder)
				allMoves = slices.Concat(allMoves, *n)
			case 3:
				var placeholder []Move = make([]Move, 0)
				n := RecurringBishopDepth(b, b.Turn, &placeholder)
				allMoves = slices.Concat(allMoves, *n)
			case 4:
				squares := TrailingZerosLoop(bb)
				for _, square := range squares {
					if square < 64 {
						adjacentknight := precomped.Knight[square] &^ ourpieces
						after := adjacentknight.ToSquares()
						for _, v := range after {
							Move := Move{From: square, To: v}
							allMoves = append(allMoves, Move)
						}
					}
				}

			case 5:
				push1pawns := (bb >> 8) &^ b.FilledSquares
				after := push1pawns.ToSquares()
				for _, v := range after {
					from := v + 8
					to := v
					if to <= 7 {
						allMoves = append(allMoves,
							Move{From: from, To: to, Promotion: 1}, //queen
							Move{From: from, To: to, Promotion: 2}, //rook
							Move{From: from, To: to, Promotion: 3}, //bishop
							Move{From: from, To: to, Promotion: 4}, //knight
						)
					} else {
						Move := Move{From: v + 8, To: v}
						allMoves = append(allMoves, Move)
					}
				}
				push2pawns := ((push1pawns & (WPawnStartRank >> 8)) >> 8) &^ b.FilledSquares
				after = push2pawns.ToSquares()
				for _, v := range after {
					Move := Move{From: v + 16, To: v}
					allMoves = append(allMoves, Move)
				}

				leftcapture := ((bb &^ FileH) >> 9) & opponentpieces
				after = leftcapture.ToSquares()
				for _, v := range after {
					if v <= 7 { // Promotion rank
						allMoves = append(allMoves,
							Move{From: v + 9, To: v, Promotion: 1}, //queen
							Move{From: v + 9, To: v, Promotion: 2}, //rook
							Move{From: v + 9, To: v, Promotion: 3}, //bishop
							Move{From: v + 9, To: v, Promotion: 4}, //knight
						)
					} else {
						Move := Move{From: v + 9, To: v}
						allMoves = append(allMoves, Move)
					}
				}

				// Right capture with promotions
				rightcapture := ((bb &^ FileA) >> 7) & opponentpieces
				after = rightcapture.ToSquares()
				for _, v := range after {
					if v <= 7 { // Promotion rank
						allMoves = append(allMoves,
							Move{From: v + 7, To: v, Promotion: 1}, //queen
							Move{From: v + 7, To: v, Promotion: 2}, //rook
							Move{From: v + 7, To: v, Promotion: 3}, //bishop
							Move{From: v + 7, To: v, Promotion: 4}, //knight
						)
					} else {
						Move := Move{From: v + 7, To: v}
						allMoves = append(allMoves, Move)
					}
				}
			}
		}
		//castling
		if b.WCastleQ {
			if !(b.FilledSquares.IsSet(1) || b.FilledSquares.IsSet(2) || b.FilledSquares.IsSet(3)) {
				if !(b.IsSquareAttacked(1) || b.IsSquareAttacked(2) || b.IsSquareAttacked(3)) {
					move := Move{From: 4, To: 0, Castle: 1}
					allMoves = append(allMoves, move)
				}
			}
		}
		if b.WCastleK {
			if !(b.FilledSquares.IsSet(5) || b.FilledSquares.IsSet(6)) {
				if !(b.IsSquareAttacked(5) || b.IsSquareAttacked(6)) {
					move := Move{From: 4, To: 7, Castle: 2}
					allMoves = append(allMoves, move)
				}
			}
		}
	} else {
		if b.IsInCheck() {
			return nil
		}
		for i, bb := range []Bitboard{b.BKings, b.BQueens, b.BRooks, b.BBishops, b.BKnights, b.BPawns} {
			switch i {
			case 0:
				square := bits.TrailingZeros64(uint64(bb))
				if square < 64 {
					adjacentking := precomped.King[square] &^ ourpieces
					after := adjacentking.ToSquares()
					for _, v := range after {
						Move := Move{From: square, To: v}
						allMoves = append(allMoves, Move)
					}
				}
			case 1:
				var placeholder []Move = make([]Move, 0)
				n := RecurringQueenDepth(b, b.Turn, &placeholder)
				allMoves = slices.Concat(allMoves, *n)
			case 2:
				var placeholder []Move = make([]Move, 0)
				n := RecurringRookDepth(b, b.Turn, &placeholder)
				allMoves = slices.Concat(allMoves, *n)
			case 3:
				var placeholder []Move = make([]Move, 0)
				n := RecurringBishopDepth(b, b.Turn, &placeholder)
				allMoves = slices.Concat(allMoves, *n)
			case 4:
				squares := TrailingZerosLoop(bb)
				for _, square := range squares {
					if square < 64 {
						adjacentknight := precomped.Knight[square] &^ ourpieces
						after := adjacentknight.ToSquares()
						for _, v := range after {
							Move := Move{From: square, To: v}
							allMoves = append(allMoves, Move)
						}
					}
				}

			case 5:
				push1pawns := (bb << 8) &^ b.FilledSquares
				after := push1pawns.ToSquares()
				for _, v := range after {
					from := v - 8
					to := v
					if to >= 56 {
						allMoves = append(allMoves,
							Move{From: from, To: to, Promotion: 1}, //queen
							Move{From: from, To: to, Promotion: 2}, //rook
							Move{From: from, To: to, Promotion: 3}, //bishop
							Move{From: from, To: to, Promotion: 4}, //knight
						)
					} else {
						Move := Move{From: v - 8, To: v}
						allMoves = append(allMoves, Move)
					}
				}
				push2pawns := ((push1pawns & (BPawnStartRank << 8)) << 8) &^ b.FilledSquares
				after = push2pawns.ToSquares()
				for _, v := range after {
					Move := Move{From: v - 16, To: v}
					allMoves = append(allMoves, Move)
				}

				leftcapture := ((bb &^ FileH) << 9) & opponentpieces
				after = leftcapture.ToSquares()
				for _, v := range after {
					if v >= 56 { // Promotion rank
						allMoves = append(allMoves,
							Move{From: v - 9, To: v, Promotion: 1}, //queen
							Move{From: v - 9, To: v, Promotion: 2}, //rook
							Move{From: v - 9, To: v, Promotion: 3}, //bishop
							Move{From: v - 9, To: v, Promotion: 4}, //knight
						)
					} else {
						Move := Move{From: v - 9, To: v}
						allMoves = append(allMoves, Move)
					}
				}

				// Right capture with promotions
				rightcapture := ((bb &^ FileA) << 7) & opponentpieces
				after = rightcapture.ToSquares()
				for _, v := range after {
					if v >= 56 { // Promotion rank
						allMoves = append(allMoves,
							Move{From: v - 7, To: v, Promotion: 1}, //queen
							Move{From: v - 7, To: v, Promotion: 2}, //rook
							Move{From: v - 7, To: v, Promotion: 3}, //bishop
							Move{From: v - 7, To: v, Promotion: 4}, //knight
						)
					} else {
						Move := Move{From: v - 7, To: v}
						allMoves = append(allMoves, Move)
					}
				}
			}
		}
		//castling
		if b.BCastleQ {
			if !(b.FilledSquares.IsSet(57) || b.FilledSquares.IsSet(58) || b.FilledSquares.IsSet(59)) {
				if !(b.IsSquareAttacked(57) || b.IsSquareAttacked(58) || b.IsSquareAttacked(59)) {
					move := Move{From: 60, To: 56, Castle: 3}
					allMoves = append(allMoves, move)
				}
			}
		}
		if b.BCastleK {
			if !(b.FilledSquares.IsSet(61) || b.FilledSquares.IsSet(62)) {
				if !(b.IsSquareAttacked(61) || b.IsSquareAttacked(62)) {
					move := Move{From: 60, To: 63, Castle: 4}
					allMoves = append(allMoves, move)
				}
			}
		}
	}

	return allMoves
}
