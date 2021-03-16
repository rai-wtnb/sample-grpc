package game

import "fmt"

type Board struct {
	Cells [][]Color
}

// 盤面作成するファクトリ
func NewBoard() *Board {
	b := &Board{
		Cells: make([][]Color, 10),
	}

	for i := 0; i < 10; i++ {
		b.Cells[i] = make([]Color, 10)
	}

	// 盤面両端を壁に
	for i := 0; i < 10; i++ {
		b.Cells[0][i] = Wall
	}
	for i := 0; i < 10; i++ {
		b.Cells[i][0] = Wall
		b.Cells[i][9] = Wall
	}
	for i := 0; i < 10; i++ {
		b.Cells[9][i] = Wall
	}

	// 初期石配置
	b.Cells[4][4] = White
	b.Cells[5][5] = White
	b.Cells[4][5] = Black
	b.Cells[5][4] = Black

	return b
}

func (b *Board) PutStone(x, y int32, c Color) error {
	if !b.CanPutStone(x, y, c) {
		return fmt.Errorf("Cannot put stone x=%v, y=%v color=%v", x, y, ColorToStr(c))
	}

	// 石を配置
	b.Cells[x][y] = c

	// 置いた石の縦/横/斜めの各方向で反転できる石を全て反転
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			if b.CountTurnableStonesByDirection(x, y, c, int32(dx), int32(dy)) > 0 {
				b.TurnStonesByDirection(x, y, c, int32(dx), int32(dy))
			}
		}
	}

	return nil
}

func (b *Board) CanPutStone(x, y int32, c Color) bool {
	// すでに置かれている場合false
	if b.Cells[x][y] != Empty {
		return false
	}

	// 反転できる石があればtrue
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			if b.CountTurnableStonesByDirection(x, y, int32(dx), int32(dy), c) > 0 {
				return true
			}
		}
	}

	return false
}

func (b *Board) CountTurnableStonesByDirection(x, y, dx, dy int32, c Color) int {
	cnt := 0
	nx := x + dx
	ny := y + dy
	for {
		nc := b.Cells[nx][ny]
		if nc != OpponentColor(c) {
			break
		}

		cnt++

		nx += dx
		ny += dy
	}

	// その方向の相手の石が0より大きく、その先に自石がある時数を返す
	if cnt > 0 && b.Cells[nx][ny] == c {
		return cnt
	}

	return 0
}

// ある方向の石をひっくり返す。ひっくり返しても良い場合だけ呼ぶ。
func (b *Board) TurnStonesByDirection(x, y, dx, dy int32, c Color) {
	nx := x + dx
	ny := y + dy

	for {
		nc := b.Cells[nx][ny]

		if nc != OpponentColor(c) {
			break
		}

		b.Cells[nx][ny] = c

		nx += dx
		ny += dy
	}
}

// 盤面内である色の石を置けるマスの数を数える
func (b *Board) AvailableCellCount(c Color) int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.CanPutStone(int32(i), int32(j), c) {
				cnt++
			}
		}
	}

	return cnt
}

// 引数の色の石を数える
func (b *Board) Score(c Color) int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.Cells[i][j] != c {
				continue
			}
			cnt++
		}
	}

	return cnt
}

// 残りのマス(Empty)をカウント
func (b *Board) Rest() int {
	cnt := 0

	for i := 1; i < 9; i++ {
		for j := 1; j < 9; j++ {
			if b.Cells[i][j] == Empty {
				cnt++
			}
		}
	}

	return cnt
}
