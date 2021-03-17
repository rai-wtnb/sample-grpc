package game

import "fmt"

type Game struct {
	Board    *Board
	started  bool
	finished bool
	me       Color
}

// ファクトリ
func NewGame(me Color) *Game {
	return &Game{
		Board: NewBoard(),
		me:    me,
	}
}

// 手を打ち、盤面を出力
// 返り値としてゲーム終了の有無を返す
func (g *Game) Move(x, y int32, c Color) (bool, error) {
	if g.finished {
		return true, nil
	}

	if err := g.Board.PutStone(x, y, c); err != nil {
		return false, err
	}

	g.Display(g.me)

	if g.IsGameOver() {
		fmt.Println("game finished.")
		g.finished = true
		return true, nil
	}

	return false, nil
}

// 両者がおける場所がないか確認。両者ともなければtrueを返す
func (g *Game) IsGameOver() bool {
	if g.Board.AvailableCellCount(Black) > 0 {
		return false
	}
	if g.Board.AvailableCellCount(White) > 0 {
		return false
	}

	return true
}

func (g *Game) Winner() Color {
	black := g.Board.Score(Black)
	white := g.Board.Score(White)

	if black == white {
		return None
	} else if black > white {
		return Black
	}
	return White
}

func (g *Game) Display(me Color) {
	fmt.Println("")
	if me != None {
		fmt.Printf("You: %v\n", ColorToStr(me))
	}

	fmt.Print(" ｜")
	rs := []rune("ABCDEFGH")
	for i, r := range rs {
		fmt.Printf("%v", string(r))
		if i < len(rs)-1 {
			fmt.Print("｜")
		}
	}
	fmt.Print("\n")
	fmt.Println("ーーーーーーーーーーーーーー")

	for j := 1; j < 9; j++ {
		fmt.Printf("%d", j)
		fmt.Print("｜")
		for i := 1; i < 9; i++ {
			fmt.Print(ColorToStr(g.Board.Cells[i][j]))
			fmt.Print("｜")
		}
		fmt.Print("\n")
	}

	fmt.Println("ーーーーーーーーーーーーーー")

	fmt.Printf("Score: BLACK=%d, WHITE=%d REST=%d\n", g.Board.Score(Black), g.Board.Score(White), g.Board.Rest())

	fmt.Print("\n")
}
