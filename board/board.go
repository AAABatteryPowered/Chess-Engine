package board

import (
	"bot/moves"
	"fmt"
	"math/bits"
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

type Undo struct {
	from, to      int8
	movingPiece   int8
	capturedPiece int8
	promotion     uint8
	enPassantOld  int8
	wCastleKOld   bool
	wCastleQOld   bool
	bCastleKOld   bool
	bCastleQOld   bool
	turnOld       bool
}

type BoardMethods interface {
	FromFen(string)
	GenMoves() []moves.Move
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

	EnPassantTarget int8

	AllBitboards [12]*Bitboard
	UndoStack    [64]Undo
	UndoCount    int

	HalfMoves int
	FullMoves int

	FilledSquares Bitboard
	Turn          bool
}

type preCompTables struct {
	King   [64]Bitboard
	Knight [64]Bitboard
	Pawn   [64]Bitboard
}

/*
type Move struct {
	From      int
	To        int
	Castle    int
	EnPassant bool
	Promotion int
}*/

type Bitboard uint64

func algebraicToSquare(notation string) int8 {
	if len(notation) != 2 {
		return -1
	}

	file := int(notation[0] - 'a')
	rank := int(notation[1] - '1')

	return int8(rank*8 + file)
}

func (b *Bitboard) Set(pos int8) {
	*b |= 1 << pos
}

func (b *Bitboard) Clear(pos int8) {
	*b &^= 1 << pos
}

func (b *Bitboard) Toggle(pos int8) {
	*b ^= 1 << pos
}

func (b *Board) Pieces() (Bitboard, Bitboard) {
	if !b.Turn {
		return b.BKings | b.BQueens | b.BRooks | b.BBishops | b.BKnights | b.BPawns, b.WKings | b.WQueens | b.WRooks | b.WBishops | b.WKnights | b.WPawns
	} else {
		return b.WKings | b.WQueens | b.WRooks | b.WBishops | b.WKnights | b.WPawns, b.BKings | b.BQueens | b.BRooks | b.BBishops | b.BKnights | b.BPawns
	}
}

func (b Bitboard) IsSet(pos int8) bool {
	return (b>>pos)&1 == 1
}

func (b *Board) SetTurn(t bool) {
	b.Turn = t
}

func notationToSquare(notation string) int8 {
	if len(notation) != 2 {
		return -1
	}

	file := notation[0] - 'a'
	rank := notation[1] - '1'

	if file < 0 || file > 7 || rank < 0 || rank > 7 {
		return -1
	}

	return int8(rank)*8 + int8(file)
}

func (b *Board) FromFen(s string) {
	b.AllBitboards = [12]*Bitboard{
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

	var posPointer int8 = 0
	var subdata string
outerLoop:
	for xx, ch := range s {
		switch ch {
		case '1', '2', '3', '4', '5', '6', '7', '8':
			num := int8(ch - '0')
			posPointer += num
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
}

func (b *Board) Copy() *Board {
	newBoard := &Board{}
	*newBoard = *b

	// newBoard.pieces = make([]Piece, len(b.pieces))
	// copy(newBoard.pieces, b.pieces)

	return newBoard
}

func (b *Board) PieceAt(square int8) int8 {
	if !b.FilledSquares.IsSet(square) {
		return -1
	}

	if b.WPawns.IsSet(square) {
		return 5
	}
	if b.BPawns.IsSet(square) {
		return 11
	}

	if b.WKnights.IsSet(square) {
		return 4
	}
	if b.BKnights.IsSet(square) {
		return 10
	}

	if b.WBishops.IsSet(square) {
		return 3
	}
	if b.BBishops.IsSet(square) {
		return 9
	}

	if b.WRooks.IsSet(square) {
		return 2
	}
	if b.BRooks.IsSet(square) {
		return 8
	}

	if b.WQueens.IsSet(square) {
		return 1
	}
	if b.BQueens.IsSet(square) {
		return 7
	}

	if b.WKings.IsSet(square) {
		return 0
	}
	if b.BKings.IsSet(square) {
		return 6
	}

	return -1
}

func (b Bitboard) ToSquares() []int {
	bb := uint64(b)
	count := bits.OnesCount64(bb)

	squares := make([]int, 0, count)

	for bb != 0 {
		sq := bits.TrailingZeros64(bb)
		squares = append(squares, sq)
		bb &= bb - 1
	}

	return squares
}

func (b *Board) DebugPrint() {
	var finalstr string
	newlinesadded := 0
	allbbs := &b.AllBitboards
	for pos := 0; pos < 64; pos++ {
		piecefound := false
		for i, bb := range allbbs {
			if bb.IsSet(int8(pos)) {
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
	for rank := 7; rank >= 0; rank-- {
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
	occ := b.FilledSquares
	squareBit := Bitboard(1 << square)

	if b.Turn {
		// White to move → check if square is attacked by black

		// Pawn attacks (black pawns move down, attack diagonally down)
		if (((b.BPawns<<7)&^FileH)|((b.BPawns<<9)&^FileA))&squareBit != 0 {
			return true
		}

		// Knight attacks
		if precomped.Knight[square]&b.BKnights != 0 {
			return true
		}

		// King attacks
		if precomped.King[square]&b.BKings != 0 {
			return true
		}

		// Rook/Queen attacks using magic bitboards
		rookMagic := rookMagics[square]
		rookHash := uint64(occ&rookMagic.Mask) * rookMagic.Magic
		rookIndex := (rookHash >> rookMagic.Shift) + uint64(rookMagic.Offset)
		if rookAttacks[rookIndex]&(b.BRooks|b.BQueens) != 0 {
			return true
		}

		// Bishop/Queen attacks using magic bitboards
		bishopMagic := bishopMagics[square]
		bishopHash := uint64(occ&bishopMagic.Mask) * bishopMagic.Magic
		bishopIndex := (bishopHash >> bishopMagic.Shift) + uint64(bishopMagic.Offset)
		if bishopAttacks[bishopIndex]&(b.BBishops|b.BQueens) != 0 {
			return true
		}
	} else {
		// Black to move → check if square is attacked by white

		// Pawn attacks (white pawns move up, attack diagonally up)
		if (((b.WPawns>>7)&^FileA)|((b.WPawns>>9)&^FileH))&squareBit != 0 {
			return true
		}

		// Knight attacks
		if precomped.Knight[square]&b.WKnights != 0 {
			return true
		}

		// King attacks
		if precomped.King[square]&b.WKings != 0 {
			return true
		}

		// Rook/Queen attacks using magic bitboards
		rookMagic := rookMagics[square]
		rookHash := uint64(occ&rookMagic.Mask) * rookMagic.Magic
		rookIndex := (rookHash >> rookMagic.Shift) + uint64(rookMagic.Offset)
		if rookAttacks[rookIndex]&(b.WRooks|b.WQueens) != 0 {
			return true
		}

		// Bishop/Queen attacks using magic bitboards
		bishopMagic := bishopMagics[square]
		bishopHash := uint64(occ&bishopMagic.Mask) * bishopMagic.Magic
		bishopIndex := (bishopHash >> bishopMagic.Shift) + uint64(bishopMagic.Offset)
		if bishopAttacks[bishopIndex]&(b.WBishops|b.WQueens) != 0 {
			return true
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

	if b.Turn {
		kingSq = bits.TrailingZeros64(uint64(b.WKings))

		if (((b.BPawns<<7)&^FileH)|((b.BPawns<<9)&^FileA))&(1<<kingSq) != 0 {
			return true
		}

		if precomped.Knight[kingSq]&b.BKnights != 0 {
			return true
		}

		if precomped.King[kingSq]&b.BKings != 0 {
			return true
		}

		rookMagic := rookMagics[kingSq]
		rookHash := uint64(occ&rookMagic.Mask) * rookMagic.Magic
		rookIndex := (rookHash >> rookMagic.Shift) + uint64(rookMagic.Offset)
		rookAtk := rookAttacks[rookIndex]
		if rookAtk&(b.BRooks|b.BQueens) != 0 {
			return true
		}

		bishopMagic := bishopMagics[kingSq]
		bishopHash := uint64(occ&bishopMagic.Mask) * bishopMagic.Magic
		bishopIndex := (bishopHash >> bishopMagic.Shift) + uint64(bishopMagic.Offset)
		bishopAtk := bishopAttacks[bishopIndex]
		if bishopAtk&(b.BBishops|b.BQueens) != 0 {
			return true
		}
	} else {
		kingSq = bits.TrailingZeros64(uint64(b.BKings))

		if (((b.WPawns>>7)&^FileA)|((b.WPawns>>9)&^FileH))&(1<<kingSq) != 0 {
			return true
		}

		if precomped.Knight[kingSq]&b.WKnights != 0 {
			return true
		}

		if precomped.King[kingSq]&b.WKings != 0 {
			return true
		}

		rookMagic := rookMagics[kingSq]
		rookHash := uint64(occ&rookMagic.Mask) * rookMagic.Magic
		rookIndex := (rookHash >> rookMagic.Shift) + uint64(rookMagic.Offset)
		rookAtk := rookAttacks[rookIndex]
		if rookAtk&(b.WRooks|b.WQueens) != 0 {
			return true
		}

		bishopMagic := bishopMagics[kingSq]
		bishopHash := uint64(occ&bishopMagic.Mask) * bishopMagic.Magic
		bishopIndex := (bishopHash >> bishopMagic.Shift) + uint64(bishopMagic.Offset)
		bishopAtk := bishopAttacks[bishopIndex]
		if bishopAtk&(b.WBishops|b.WQueens) != 0 {
			return true
		}
	}

	return false
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

func (b *Board) Moves() moves.MoveList {
	ourmoves := b.GenMoves()
	var filteredmoves moves.MoveList = moves.NewMoveList()
	/*incheck := b.IsKingAttacked()
	if !incheck {
		return ourmoves
	}*/
	originalturn := b.Turn
	for i := 0; i < ourmoves.Count; i++ {
		ourmove := ourmoves.Moves[i]
		b.PlayMove(ourmove)
		b.Turn = originalturn
		stillincheck := b.IsKingAttacked()
		if !stillincheck {
			filteredmoves.Add(ourmove)
		}
		b.UndoMove(ourmove)
	}
	return filteredmoves
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

func (b *Board) PlayMove(move moves.Move) {
	movingpiece := b.PieceAt(move.From())
	if movingpiece == -1 {
		fmt.Println(move)
	}
	targetpiece := b.PieceAt(move.To())

	u := &b.UndoStack[b.UndoCount]
	u.from = move.From()
	u.to = move.To()
	u.movingPiece = movingpiece
	u.capturedPiece = int8(targetpiece)
	u.promotion = move.PromotionPiece()
	u.enPassantOld = b.EnPassantTarget
	u.wCastleKOld = b.WCastleK
	u.wCastleQOld = b.WCastleQ
	u.bCastleKOld = b.BCastleK
	u.bCastleQOld = b.BCastleQ
	u.turnOld = b.Turn
	b.UndoCount++

	b.FilledSquares.Clear(move.From())
	b.EnPassantTarget = -1
	castling := move.IsCastling()
	allbb := &b.AllBitboards
	if castling {
		switch move.To() {
		case 56:
			//no clear move.froms cuz we do that 6 lines above
			b.FilledSquares.Clear(move.To())
			b.FilledSquares.Set(58)
			b.FilledSquares.Set(59)
			allbb[0].Clear(move.From()) // white king
			allbb[0].Set(58)
			allbb[2].Clear(move.To()) // white rook
			allbb[2].Set(59)
		case 63:
			b.FilledSquares.Clear(move.To())
			b.FilledSquares.Set(61)
			b.FilledSquares.Set(62)
			allbb[0].Clear(move.From()) // white king
			allbb[0].Set(62)
			allbb[2].Clear(move.To()) // white rook
			allbb[2].Set(61)
		case 0:
			b.FilledSquares.Clear(move.To())
			b.FilledSquares.Set(2)
			b.FilledSquares.Set(3)
			allbb[6].Clear(move.From())
			allbb[6].Set(2)
			allbb[8].Clear(move.To())
			allbb[8].Set(3)
		case 7:
			b.FilledSquares.Clear(move.To())
			b.FilledSquares.Set(5)
			b.FilledSquares.Set(6)
			allbb[6].Clear(move.From())
			allbb[6].Set(6)
			allbb[8].Clear(move.To())
			allbb[8].Set(5)
		}
	} else if move.IsPromotion() {
		b.FilledSquares.Set(move.To())
		allbb[movingpiece].Clear(move.From())
		if targetpiece != -1 {
			allbb[targetpiece].Clear(move.To())
		}
		if b.Turn {
			allbb[move.PromotionPiece()].Set(move.To())
		} else {
			allbb[move.PromotionPiece()+6].Set(move.To())
		}
	} else if move.IsEnPassant() && movingpiece != -1 {
		allbb[movingpiece].Clear(move.From())
		b.FilledSquares.Set(move.To())
		if targetpiece != -1 {
			allbb[targetpiece].Clear(move.To())
		}
		allbb[movingpiece].Set(move.To())

		if b.Turn {
			b.FilledSquares.Clear(move.To() + 8)
			allbb[11].Clear(move.To() + 8)
		} else {
			b.FilledSquares.Clear(move.To() - 8)
			allbb[5].Clear(move.To() - 8)
		}
	} else if targetpiece > 5 {
		// piece is black
		if b.Turn {
			b.FilledSquares.Set(move.To())
			allbb[movingpiece].Clear(move.From())
			allbb[targetpiece].Clear(move.To())
			allbb[movingpiece].Set(move.To())
		}
	} else if targetpiece == -1 {
		b.FilledSquares.Set(move.To())
		allbb[movingpiece].Clear(move.From())
		allbb[movingpiece].Set(move.To())

		if b.Turn && movingpiece == 5 {
			if move.From()-move.To() == 16 {
				b.EnPassantTarget = move.From() - 8
			}
		} else if !b.Turn && movingpiece == 11 {
			if move.To()-move.From() == 16 {
				b.EnPassantTarget = move.From() + 8
			}
		}
	} else if targetpiece > -1 && targetpiece < 6 {
		// piece is white
		if !b.Turn {
			b.FilledSquares.Set(move.To())
			allbb[movingpiece].Clear(move.From())
			allbb[targetpiece].Clear(move.To())
			allbb[movingpiece].Set(move.To())
		}
	}

	if movingpiece == 0 { // white king moved
		b.WCastleK = false
		b.WCastleQ = false
	} else if movingpiece == 6 { // black king moved
		b.BCastleK = false
		b.BCastleQ = false
	}
	if movingpiece == 2 || targetpiece == 2 { //white rook move or takne
		if move.From() == 0 || move.To() == 0 {
			b.WCastleQ = false
		}
		if move.From() == 7 || move.To() == 7 {
			b.WCastleK = false
		}
	}
	if movingpiece == 8 || targetpiece == 8 { //black rook moved or taken
		if move.From() == 56 || move.To() == 56 { //it move from square 56 or someone captured square 56
			b.BCastleQ = false
		}
		if move.From() == 63 || move.To() == 63 {
			b.BCastleK = false
		}
	}

	b.Turn = !b.Turn
}

func (b *Board) UndoMove(move moves.Move) {

	b.UndoCount--
	u := &b.UndoStack[b.UndoCount]

	allbb := &b.AllBitboards

	b.Turn = u.turnOld

	b.WCastleK = u.wCastleKOld
	b.WCastleQ = u.wCastleQOld
	b.BCastleK = u.bCastleKOld
	b.BCastleQ = u.bCastleQOld

	b.EnPassantTarget = u.enPassantOld

	b.FilledSquares.Clear(u.to)

	allbb[u.movingPiece].Clear(u.to)

	if u.capturedPiece != -1 {
		allbb[u.capturedPiece].Set(u.to)
		b.FilledSquares.Set(u.to)
	}

	allbb[u.movingPiece].Set(u.from)
	b.FilledSquares.Set(u.from)

	if u.promotion != 0 {
		promotedIndex := u.promotion
		if !u.turnOld {
			promotedIndex += 6
		}

		allbb[promotedIndex].Clear(u.to)

		allbb[u.movingPiece].Set(u.from)
	}

	// handle en-passant capture
	if move.IsEnPassant() {
		if u.turnOld { // white moved
			allbb[11].Set(u.to + 8)
			b.FilledSquares.Set(u.to + 8)
		} else { // black moved
			allbb[5].Set(u.to - 8)
			b.FilledSquares.Set(u.to - 8)
		}
	}

	if !move.IsCastling() {
		return
	}
	switch move.To() {
	case 56:
		allbb[0].Clear(58)
		allbb[2].Clear(59)
		allbb[0].Set(60)
		allbb[2].Set(56)
		b.FilledSquares.Clear(58)
		b.FilledSquares.Clear(59)
		b.FilledSquares.Set(60)
		b.FilledSquares.Set(56)

	case 63:
		allbb[0].Clear(62)
		allbb[2].Clear(61)
		allbb[0].Set(60)
		allbb[2].Set(63)
		b.FilledSquares.Clear(62)
		b.FilledSquares.Clear(61)
		b.FilledSquares.Set(60)
		b.FilledSquares.Set(63)

	case 0:
		allbb[6].Clear(2)
		allbb[8].Clear(3)
		allbb[6].Set(4)
		allbb[8].Set(0)
		b.FilledSquares.Clear(2)
		b.FilledSquares.Clear(3)
		b.FilledSquares.Set(4)
		b.FilledSquares.Set(0)

	case 7:
		allbb[6].Clear(6)
		allbb[8].Clear(5)
		allbb[6].Set(4)
		allbb[8].Set(7)
		b.FilledSquares.Clear(6)
		b.FilledSquares.Clear(5)
		b.FilledSquares.Set(4)
		b.FilledSquares.Set(7)
	}
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

func TrailingZerosLoop(b Bitboard) []int8 {
	var squareslist []int8
	bb := uint64(b)
	for bb != 0 {
		square := bits.TrailingZeros64(bb)
		squareslist = append(squareslist, int8(square))

		bb &= bb - 1
	}
	return squareslist
}

func (b *Board) IsInCheck() bool {

	return false
}

type Magic struct {
	Mask   Bitboard
	Magic  uint64
	Shift  uint8
	Offset uint32
}

var rookMagics [64]Magic
var rookAttacks []Bitboard

var bishopMagics [64]Magic
var bishopAttacks []Bitboard

type RookLookupKey struct {
	StartSquare  int
	Blockerboard Bitboard
}

var RookMovementMasks [64]Bitboard
var BishopMovementMasks [64]Bitboard
var BishopOccupancyMasks [64]Bitboard
var RookLookupTable map[RookLookupKey]Bitboard

var precomputedRookMagics = [64]uint64{
	0x8080002484104000, 0x8040042000401002,
	0x14800A8010002000, 0x1200084020A41200,
	0x0200040902006010, 0x0200020024100108,
	0x490000840E000100, 0x0080008000482100,
	0xA004800020804004, 0x0000401000200240,
	0x1002004011820020, 0x0042801800100084,
	0x8001000528001100, 0x0002000824100E00,
	0x0113000A00090044, 0x4603000892124100,
	0x1480010021004580, 0x8082120020830040,
	0x0200808020001000, 0x040800801000808A,
	0x0424008008028004, 0x8001010008040002,
	0x2068808001000E00, 0x04000200040880E3,
	0x0040400080008120, 0x1700600440100040,
	0x4008700080600080, 0x9040500080080480,
	0x1800840180080080, 0x004E1C0080420080,
	0x0080140101000200, 0x09000112000040A4,
	0x2000C0048A800A20, 0x0810102000400048,
	0x0020008022801000, 0x0140810802801000,
	0x0428800800802401, 0x1401810200800400,
	0x10000A0844000110, 0x0060004882000401,
	0x0300800040088020, 0x0860004000208080,
	0x0000600100410010, 0x2004100421010009,
	0x0001000801050010, 0x02A1220004008080,
	0x1002024801440010, 0x000001C084220011,
	0x281580050028C100, 0x0000220080490200,
	0x0030052000801080, 0x0004120021415A00,
	0x0282050088001100, 0x0000804400460080,
	0x0100090802901C00, 0x0022018405014200,
	0x0001401422810202, 0x0000400010802101,
	0x1101001108200045, 0x02C4C42009001001,
	0x9012002010046842, 0x2481008A24000801,
	0x00000800C1021024, 0x0384006402408502,
}

var precomputedRookShifts = [64]uint8{
	52, 53, 53, 53, 53, 53, 53, 52, 53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53, 53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53, 53, 54, 54, 54, 54, 54, 54, 53,
	53, 54, 54, 54, 54, 54, 54, 53, 52, 53, 53, 53, 53, 53, 53, 52,
}

var precomputedBishopMagics = [64]uint64{
	0x4004000400041004, 0x4004012204030084,
	0x0008080918A0000A, 0x0124105201010010,
	0x0202021062010000, 0x00020A9004080040,
	0x0203031010069002, 0x1011020004200080,
	0x2082240800080001, 0x0104025042028100,
	0x0010100129410000, 0x0000740400802000,
	0x2000141420000020, 0x0000148220210040,
	0x0800804210D00800, 0x8200100C00060030,
	0x0000408800100088, 0x002040A204010200,
	0x82408408020810A0, 0x0018008222084000,
	0x00A4000080A00202, 0x00C100208080C000,
	0x5008830202012000, 0x2902002040020040,
	0x0000224030001000, 0x0408028A04100226,
	0x000412030C080410, 0xC012006042008200,
	0x2201840200802000, 0x0000414002011000,
	0x0435040500C20881, 0x0000810015040200,
	0x1011004000080010, 0x00C1311003081080,
	0x0010280801040020, 0x6003040100100901,
	0x20A0860200040090, 0x0001050202150810,
	0x0001084100008404, 0x0808C40004000804,
	0x0082003080200200, 0x0008884808010220,
	0x0110820801000A00, 0x000060E038000100,
	0x00100803040060C0, 0x6001200080802101,
	0x0050030200908421, 0x0210000420201128,
	0x0010440008004082, 0xA102120104120C00,
	0x08C04111C1100008, 0xA002000294040020,
	0x0004109012020000, 0x0004110250044014,
	0x0040342104012402, 0x0050220080410001,
	0x8200008010080020, 0x0400006582086044,
	0x2062100110980404, 0x0008080020218800,
	0x0000020940104444, 0x06002040026C0101,
	0x2022104602040430, 0x0108080800010002,
}

var precomputedBishopShifts = [64]uint8{
	51, 59, 59, 59, 59, 59, 59, 53, 53, 59, 59, 59, 59, 59, 59, 54,
	54, 59, 57, 57, 57, 57, 59, 54, 54, 59, 57, 55, 55, 57, 59, 54,
	54, 59, 57, 55, 55, 57, 59, 54, 54, 59, 57, 57, 57, 57, 59, 54,
	54, 59, 59, 59, 59, 59, 59, 53, 53, 59, 59, 59, 59, 59, 59, 51,
}

func InitMagicBitboards() {
	GenerateMovementMasks()

	var offset uint32 = 0

	for square := 0; square < 64; square++ {
		relevantBits := 64 - precomputedRookShifts[square]

		rookMagics[square] = Magic{
			Mask:   RookMovementMasks[square],
			Magic:  precomputedRookMagics[square],
			Shift:  precomputedRookShifts[square],
			Offset: offset,
		}

		offset += 1 << relevantBits
	}

	rookAttacks = make([]Bitboard, offset)

	for square := 0; square < 64; square++ {
		magic := rookMagics[square]
		occupancies := GenerateOccupancyMasks(magic.Mask)

		for _, occ := range occupancies {
			attacks := RecurringRookDepth(square, occ)
			hash := uint64(occ&magic.Mask) * magic.Magic
			index := (hash >> magic.Shift) + uint64(magic.Offset)
			rookAttacks[index] = attacks
		}
	}

	var bishopOffset uint32 = 0
	for square := 0; square < 64; square++ {
		relevantBits := 64 - precomputedBishopShifts[square]
		bishopMagics[square] = Magic{
			Mask:   BishopMovementMasks[square],
			Magic:  precomputedBishopMagics[square],
			Shift:  precomputedBishopShifts[square],
			Offset: bishopOffset,
		}
		bishopOffset += 1 << relevantBits
	}
	bishopAttacks = make([]Bitboard, bishopOffset)
	for square := 0; square < 64; square++ {
		magic := bishopMagics[square]
		occupancies := GenerateOccupancyMasks(magic.Mask)
		for _, occ := range occupancies {
			attacks := RecurringBishopDepth(square, occ)
			hash := uint64(occ&magic.Mask) * magic.Magic
			index := (hash >> magic.Shift) + uint64(magic.Offset)
			bishopAttacks[index] = attacks
		}
	}
}

func GenerateOccupancyMasks(mask Bitboard) []Bitboard {
	bits := []int{}

	for i := 0; i < 64; i++ {
		if mask&(1<<i) != 0 {
			bits = append(bits, i)
		}
	}

	lenbits := len(bits)
	numPatterns := 1 << lenbits
	OccupancyMasks := make([]Bitboard, numPatterns)
	for patternIndex := 0; patternIndex < numPatterns; patternIndex++ {
		for bitIndex := 0; bitIndex < lenbits; bitIndex++ {
			bit := (patternIndex >> bitIndex) & 1
			OccupancyMasks[patternIndex] |= Bitboard(bit << bits[bitIndex])
		}
	}

	return OccupancyMasks
}

func GenerateMovementMasks() {
	for square := range 64 {
		var mask Bitboard = 0 //1 << square

		for i := square + 8; i < 56; i += 8 {
			mask |= 1 << i
		}
		for i := square - 8; i >= 8; i -= 8 {
			mask |= 1 << i
		}
		startFile := square % 8
		for i := square + 1; i < square-startFile+8-1; i++ {
			mask |= 1 << i
		}
		for i := square - 1; i > square-startFile; i-- {
			mask |= 1 << i
		}

		RookMovementMasks[square] = mask

		mask = 0 //1 << square

		for i := square + 9; i < 56; i += 9 {

			if i%8 == 7 {
				break
			}
			mask |= 1 << i
		}

		for i := square + 7; i < 56; i += 7 {
			if i%8 == 0 {
				break
			}
			mask |= 1 << i
		}

		for i := square - 7; i >= 8; i -= 7 {

			if i%8 == 7 {
				break
			}
			mask |= 1 << i
		}

		for i := square - 9; i >= 8; i -= 9 {

			if i%8 == 0 {
				break
			}
			mask |= 1 << i
		}

		BishopMovementMasks[square] = mask
	}
}

func GenerateRookLookuptable() {
	GenerateMovementMasks()
	for i := 0; i < 64; i++ {
		movementmask := RookMovementMasks[i]
		occmasks := GenerateOccupancyMasks(movementmask)

		for _, occupancyboard := range occmasks {
			legalmoveboard := RecurringRookDepth(i, occupancyboard)
			lkey := RookLookupKey{
				i,
				occupancyboard,
			}
			RookLookupTable[lkey] = legalmoveboard
		}
	}
}

func isValidSquare(target, start, dir int) bool {
	// Check if target is on the board
	if target < 0 || target >= 64 {
		return false
	}

	// Prevent wrapping around board edges
	startRank := start / 8
	startFile := start % 8
	targetRank := target / 8
	targetFile := target % 8

	switch dir {
	case 1: // East
		return targetRank == startRank && targetFile > startFile
	case -1: // West
		return targetRank == startRank && targetFile < startFile
	case 8: // South
		return targetFile == startFile
	case -8: // North
		return targetFile == startFile
	default:
		return false
	}
}

func RecurringRookDepth(square int, blockers Bitboard) Bitboard {
	var moveBitboard Bitboard = 0

	for _, dir := range []int{-8, 8, -1, 1} { // n s w e
		target := square
		for {
			target += dir
			if !isValidSquare(target, square, dir) {
				break
			}

			targetBit := uint64(1) << target
			moveBitboard |= Bitboard(targetBit)

			if blockers&Bitboard(targetBit) != 0 {
				break
			}
		}
	}
	return moveBitboard
}

func GenRookMoves(b *Board, pieces Bitboard) moves.MoveList {
	var friendly Bitboard

	if b.Turn {
		friendly = b.WPawns | b.WKnights | b.WBishops |
			b.WRooks | b.WQueens | b.WKings
	} else {
		friendly = b.BPawns | b.BKnights | b.BBishops |
			b.BRooks | b.BQueens | b.BKings
	}

	occupied := b.FilledSquares
	moveList := moves.NewMoveList()
	for pieces != 0 {
		var square int8 = int8(bits.TrailingZeros64(uint64(pieces)))
		pieces &= pieces - 1

		// Use the magic lookup table directly
		magic := rookMagics[square]
		hash := uint64(occupied&magic.Mask) * magic.Magic
		index := (hash >> magic.Shift) + uint64(magic.Offset)
		attacks := rookAttacks[index]

		legalMoves := attacks &^ friendly

		for legalMoves != 0 {
			var target int8 = int8(bits.TrailingZeros64(uint64(legalMoves)))
			legalMoves &= legalMoves - 1

			moveList.Add(moves.NewMove(square, target, moves.FlagNone))
		}
	}

	return moveList
}

func GenBishopMoves(b *Board, pieces Bitboard) moves.MoveList {
	var friendly Bitboard
	if b.Turn {
		friendly = b.WPawns | b.WKnights | b.WBishops |
			b.WRooks | b.WQueens | b.WKings
	} else {
		friendly = b.BPawns | b.BKnights | b.BBishops |
			b.BRooks | b.BQueens | b.BKings
	}

	occupied := b.FilledSquares
	moveList := moves.NewMoveList()

	for pieces != 0 {
		var square int8 = int8(bits.TrailingZeros64(uint64(pieces)))
		pieces &= pieces - 1

		magic := bishopMagics[square]
		hash := uint64(occupied&magic.Mask) * magic.Magic
		index := (hash >> magic.Shift) + uint64(magic.Offset)
		attacks := bishopAttacks[index]

		legalMoves := attacks &^ friendly

		for legalMoves != 0 {
			var target int8 = int8(bits.TrailingZeros64(uint64(legalMoves)))
			legalMoves &= legalMoves - 1
			moveList.Add(moves.NewMove(square, target, moves.FlagNone))
		}
	}

	return moveList
}
func GenQueenMoves(b *Board, pieces Bitboard) moves.MoveList {
	var friendly Bitboard
	if b.Turn {
		friendly = b.WPawns | b.WKnights | b.WBishops |
			b.WRooks | b.WQueens | b.WKings
	} else {
		friendly = b.BPawns | b.BKnights | b.BBishops |
			b.BRooks | b.BQueens | b.BKings
	}

	occupied := b.FilledSquares
	moveList := moves.NewMoveList()

	for pieces != 0 {
		var square int8 = int8(bits.TrailingZeros64(uint64(pieces)))
		pieces &= pieces - 1

		// Queen = Rook attacks + Bishop attacks
		rookMagic := rookMagics[square]
		rookHash := uint64(occupied&rookMagic.Mask) * rookMagic.Magic
		rookIndex := (rookHash >> rookMagic.Shift) + uint64(rookMagic.Offset)
		rookAtk := rookAttacks[rookIndex] // Changed variable name

		bishopMagic := bishopMagics[square]
		bishopHash := uint64(occupied&bishopMagic.Mask) * bishopMagic.Magic
		bishopIndex := (bishopHash >> bishopMagic.Shift) + uint64(bishopMagic.Offset)
		bishopAtk := bishopAttacks[bishopIndex] // Changed variable name

		attacks := rookAtk | bishopAtk
		legalMoves := attacks &^ friendly

		for legalMoves != 0 {
			var target int8 = int8(bits.TrailingZeros64(uint64(legalMoves)))
			legalMoves &= legalMoves - 1
			moveList.Add(moves.NewMove(square, target, moves.FlagNone))
		}
	}

	return moveList
}

func RecurringBishopDepth(square int, blockers Bitboard) Bitboard {
	var moveBitboard Bitboard = 0

	// NE (+9), NW (+7), SE (-7), SW (-9)
	for _, dir := range []int{9, 7, -7, -9} {
		target := square
		for {
			target += dir
			if target < 0 || target >= 64 {
				break
			}

			// Check for wrapping
			startFile := (target - dir) % 8
			targetFile := target % 8
			if abs(targetFile-startFile) != 1 {
				break
			}

			targetBit := uint64(1) << target
			moveBitboard |= Bitboard(targetBit)

			if blockers&Bitboard(targetBit) != 0 {
				break
			}
		}
	}
	return moveBitboard
}

func (b *Board) GenMoves() moves.MoveList {
	allMoves := moves.NewMoveList()
	ourpieces, opponentpieces := b.Pieces()

	if b.Turn {
		for i, bb := range []Bitboard{b.WKings, b.WQueens, b.WRooks, b.WBishops, b.WKnights, b.WPawns} {
			switch i {
			case 0:
				var square int8 = int8(bits.TrailingZeros64(uint64(bb)))
				if square < 64 {
					adjacentking := precomped.King[square] &^ ourpieces
					after := adjacentking.ToSquares()
					for _, v := range after {
						Move := moves.NewMove(square, int8(v), moves.FlagNone)
						//Move := Move{From: square, To: v}
						allMoves.Add(Move)
					}
				}
			case 1:
				moves := GenQueenMoves(b, bb)
				allMoves.Combine(&moves)
			case 2:
				moves := GenRookMoves(b, bb)
				allMoves.Combine(&moves)
				//var placeholder moves.MoveList
				//n := RecurringRookDepth(b, b.Turn, &placeholder)
				//allMoves.Combine(n)
			case 3:
				moves := GenBishopMoves(b, bb)
				allMoves.Combine(&moves)
			case 4:
				squares := TrailingZerosLoop(bb)
				for _, square := range squares {
					if square < 64 {
						adjacentknight := precomped.Knight[square] &^ ourpieces
						after := adjacentknight.ToSquares()
						for _, v := range after {
							Move := moves.NewMove(square, int8(v), moves.FlagNone)
							allMoves.Add(Move)
						}
					}
				}

			case 5:
				push1pawns := (bb >> 8) &^ b.FilledSquares
				after := push1pawns.ToSquares()
				for _, v := range after {
					var from int8 = int8(v + 8)
					var to int8 = int8(v)
					if to <= 7 {
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionKnight))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionBishop))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionRook))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionQueen))
					} else {
						Move := moves.NewMove(from, to, moves.FlagNone)
						allMoves.Add(Move)
					}
				}
				push2pawns := ((push1pawns & (WPawnStartRank >> 8)) >> 8) &^ b.FilledSquares
				after = push2pawns.ToSquares()
				for _, v := range after {
					Move := moves.NewMove(int8(v+16), int8(v), moves.FlagNone)
					allMoves.Add(Move)
				}

				leftcapture := ((bb &^ FileA) >> 9) & opponentpieces
				after = leftcapture.ToSquares()
				for _, v := range after {
					var from int8 = int8(v + 9)
					var to int8 = int8(v)
					if v <= 7 { // Promotion rank
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionKnight))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionBishop))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionRook))
						allMoves.Add(moves.NewMove(to, to, moves.FlagPromotionQueen))
					} else {
						Move := moves.NewMove(from, to, moves.FlagNone)
						allMoves.Add(Move)
					}
				}

				rightcapture := ((bb &^ FileH) >> 7) & opponentpieces
				after = rightcapture.ToSquares()
				for _, v := range after {
					var from int8 = int8(v + 7)
					var to int8 = int8(v)
					if v <= 7 { // Promotion rank
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionKnight))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionBishop))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionRook))
						allMoves.Add(moves.NewMove(to, to, moves.FlagPromotionQueen))
					} else {
						Move := moves.NewMove(from, to, moves.FlagNone)
						allMoves.Add(Move)
					}
				}

				if b.EnPassantTarget != -1 {
					leftEP := ((bb &^ FileA) >> 9)
					if leftEP.IsSet(b.EnPassantTarget) {
						Move := moves.NewMove(b.EnPassantTarget+9, b.EnPassantTarget, moves.FlagEnPassantCapture)
						allMoves.Add(Move)
					}

					rightEP := ((bb &^ FileH) >> 7)
					if rightEP.IsSet(b.EnPassantTarget) {
						Move := moves.NewMove(b.EnPassantTarget+7, b.EnPassantTarget, moves.FlagEnPassantCapture)
						allMoves.Add(Move)
					}
				}
			}
		}
		//castling
		if b.WCastleQ {
			if !(b.FilledSquares.IsSet(57) || b.FilledSquares.IsSet(58) || b.FilledSquares.IsSet(59)) && (b.WKings.IsSet(60) && b.WRooks.IsSet(56)) {
				if !(b.IsSquareAttacked(57) || b.IsSquareAttacked(58) || b.IsSquareAttacked(59) || b.IsSquareAttacked(60)) {
					Move := moves.NewMove(60, 56, moves.FlagCastling)
					allMoves.Add(Move)
				}
			}
		}
		if b.WCastleK {
			if !(b.FilledSquares.IsSet(61) || b.FilledSquares.IsSet(62)) && (b.WKings.IsSet(60) && b.WRooks.IsSet(63)) {
				if !(b.IsSquareAttacked(61) || b.IsSquareAttacked(62) || b.IsSquareAttacked(60)) {
					Move := moves.NewMove(60, 63, moves.FlagCastling)
					allMoves.Add(Move)
				}
			}
		}
	} else {
		for i, bb := range []Bitboard{b.BKings, b.BQueens, b.BRooks, b.BBishops, b.BKnights, b.BPawns} {
			switch i {
			case 0:
				var square int8 = int8(bits.TrailingZeros64(uint64(bb)))
				if square < 64 {
					adjacentking := precomped.King[square] &^ ourpieces
					after := adjacentking.ToSquares()
					for _, v := range after {
						Move := moves.NewMove(square, int8(v), moves.FlagNone)
						allMoves.Add(Move)
					}
				}
			case 1:
				moves := GenQueenMoves(b, bb)
				allMoves.Combine(&moves)
			case 2:
				moves := GenRookMoves(b, bb)
				allMoves.Combine(&moves)
			case 3:
				moves := GenBishopMoves(b, bb)
				allMoves.Combine(&moves)
			case 4:
				squares := TrailingZerosLoop(bb)
				for _, square := range squares {
					if square < 64 {
						adjacentknight := precomped.Knight[square] &^ ourpieces
						after := adjacentknight.ToSquares()
						for _, v := range after {
							Move := moves.NewMove(square, int8(v), moves.FlagNone)
							allMoves.Add(Move)
						}
					}
				}

			case 5:
				push1pawns := (bb << 8) &^ b.FilledSquares
				after := push1pawns.ToSquares()
				for _, v := range after {
					from := int8(v - 8)
					to := int8(v)
					if to >= 56 {
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionKnight))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionBishop))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionRook))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionQueen))
					} else {
						Move := moves.NewMove(from, to, moves.FlagNone)
						allMoves.Add(Move)
					}
				}
				push2pawns := ((push1pawns & (BPawnStartRank << 8)) << 8) &^ b.FilledSquares
				after = push2pawns.ToSquares()
				for _, v := range after {
					Move := moves.NewMove(int8(v-16), int8(v), moves.FlagNone)
					allMoves.Add(Move)
				}

				//the files may be the wrong way around like filea should be fileh and vice versa
				leftcapture := ((bb &^ FileH) << 9) & opponentpieces
				after = leftcapture.ToSquares()
				for _, v := range after {
					from := int8(v - 9)
					to := int8(v)
					if to >= 56 { // Promotion rank
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionKnight))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionBishop))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionRook))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionQueen))
					} else {
						Move := moves.NewMove(from, to, moves.FlagNone)
						allMoves.Add(Move)
					}
				}

				rightcapture := ((bb &^ FileA) << 7) & opponentpieces
				after = rightcapture.ToSquares()
				for _, v := range after {
					from := int8(v - 7)
					to := int8(v)
					if v >= 56 { // Promotion rank
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionKnight))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionBishop))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionRook))
						allMoves.Add(moves.NewMove(from, to, moves.FlagPromotionQueen))
					} else {
						Move := moves.NewMove(from, to, moves.FlagNone)
						allMoves.Add(Move)
					}
				}

				if b.EnPassantTarget != -1 {
					leftEP := ((bb &^ FileH) << 9)
					if leftEP.IsSet(b.EnPassantTarget) {

						Move := moves.NewMove(b.EnPassantTarget-9, b.EnPassantTarget, moves.FlagEnPassantCapture)
						allMoves.Add(Move)
					}

					rightEP := ((bb &^ FileA) << 7)
					if rightEP.IsSet(b.EnPassantTarget) {
						Move := moves.NewMove(b.EnPassantTarget-7, b.EnPassantTarget, moves.FlagEnPassantCapture)
						allMoves.Add(Move)
					}
				}
			}
		}
		//castling
		if b.BCastleQ {
			if !(b.FilledSquares.IsSet(1) || b.FilledSquares.IsSet(2) || b.FilledSquares.IsSet(3)) && (b.BKings.IsSet(4) && b.BRooks.IsSet(0)) {
				if !(b.IsSquareAttacked(1) || b.IsSquareAttacked(2) || b.IsSquareAttacked(3) || b.IsSquareAttacked(4)) {
					Move := moves.NewMove(4, 0, moves.FlagCastling)
					allMoves.Add(Move)
				}
			}
		}
		if b.BCastleK {
			if !(b.FilledSquares.IsSet(5) || b.FilledSquares.IsSet(6)) && (b.BKings.IsSet(4) && b.BRooks.IsSet(7)) {
				if !(b.IsSquareAttacked(5) || b.IsSquareAttacked(6) || b.IsSquareAttacked(4)) {
					Move := moves.NewMove(4, 7, moves.FlagCastling)
					allMoves.Add(Move)
				}
			}
		}
	}

	return allMoves
}
