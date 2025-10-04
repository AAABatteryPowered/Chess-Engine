package evaluation

import (
	"bot/board"
	"fmt"
)

//var kingHeatMap [64]int = 

func Evaluate(b *board.Board) int {
	final := 0

	allboards := b.AllBitboards()
	//var piecevalues [12]int
	var piececounts [12]int
	for i, brd := range allboards {
		piececounts[i] = len(board.TrailingZerosLoop(brd))
	}

	final += piececounts[]
	fmt.Println(final)

	return final
}
