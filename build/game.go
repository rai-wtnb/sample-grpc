package build

// pbパッケージの構造体を、別途用意したゲーム内で使用する構造体に変換しする処理を挟む
// → APIインターフェース側の仕様変更が基本ロジックに与える影響を最小にするため

import (
	"fmt"
	"reversi/game"
	"reversi/gen/pb"
)

func Room(r *pb.Room) *game.Room {
	return &game.Room{
		ID:    r.GetId(),
		Host:  Player(r.GetHost()),
		Guest: Player(r.GetGuest()),
	}
}

func Player(p *pb.Player) *game.Player {
	return &game.Player{
		ID:    p.GetId(),
		Color: Color(p.GetColor()),
	}
}

func Color(c pb.Color) game.Color {
	switch c {
	case pb.Color_BLACK:
		return game.Black
	case pb.Color_WHITE:
		return game.White
	case pb.Color_EMPTY:
		return game.Empty
	case pb.Color_WALL:
		return game.Wall
	}
	panic(fmt.Sprintf("unknown color: %v", c))
}
