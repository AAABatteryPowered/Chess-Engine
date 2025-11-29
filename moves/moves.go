package moves

type Move uint16

type MoveList struct {
	Moves [218]Move
	Count int
}

/*
16 bits
Bits 0-5 = from
Bits 6-11 = to
Bit 12 = castle
Bit 13 = en passant
Bit 14/15 = promotion
*/

const (
	FlagNone             = iota << 12 // 0 << 12 = 0
	FlagEnPassantCapture              // 1 << 12 = 4096
	FlagCastling                      // 2 << 12 = 8192
	FlagPromotionQueen                // 3 << 12 = 12288
	FlagPromotionRook                 // 4 << 12 = 16384
	FlagPromotionBishop               // 5 << 12 = 20480
	FlagPromotionKnight               // 6 << 12 = 24576
)

func NewMove(from, to int, flags Move) Move {
	return Move(from) | Move(to)<<6 | flags
}

func NewMoveList() MoveList {
	return MoveList{
		Count: 0,
	}
}

func (m Move) From() int {
	return int(m & 0x3F) // bits 0-5
}

func (m Move) To() int {
	return int((m >> 6) & 0x3F) // bits 6-11
}

func (m Move) Flags() Move {
	return m & 0xF000 // bits 12-15
}

func (m Move) IsEnPassant() bool {
	return m.Flags() == FlagEnPassantCapture
}

func (m Move) IsCastling() bool {
	return m.Flags() == FlagCastling
}

func (m Move) IsPromotion() bool {
	flags := m.Flags()
	return flags >= FlagPromotionKnight && flags <= FlagPromotionQueen
}

func (m Move) PromotionPiece() int {
	//0=none, 1=knight, 2=bishop, 3=rook, 4=queen
	if !m.IsPromotion() {
		return 0
	}
	return int(m.Flags()>>12) - 2
}

func (ml *MoveList) Add(move Move) {
	ml.Moves[ml.Count] = move
	ml.Count++
}

func (ml *MoveList) Combine(other *MoveList) {
	copy(ml.Moves[ml.Count:], other.Moves[:other.Count])
	ml.Count += other.Count
}
