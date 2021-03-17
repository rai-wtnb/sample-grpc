package handler

import (
	"context"
	"fmt"
	"reversi/game"
	"reversi/gen/pb"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"reversi/build"
)

type MatchingHandler struct {
	sync.RWMutex
	Rooms       map[int32]*game.Room
	maxPlayerID int32
}

func NewMatchingHandler() *MatchingHandler {
	return &MatchingHandler{
		Rooms: make(map[int32]*game.Room),
	}
}

func (h *MatchingHandler) JoinRoom(
	req *pb.JoinRoomRequest,
	stream pb.MatchingService_JoinRoomServer,
) error {
	ctx, cancel := context.WithTimeout(stream.Context(), 2*time.Minute)
	defer cancel()
	h.Lock()

	// プレイヤー新規作成
	h.maxPlayerID++
	me := &game.Player{
		ID: h.maxPlayerID,
	}

	// 空き部屋を探す
	for _, room := range h.Rooms {
		if room.Guest == nil {
			me.Color = game.White
			room.Guest = me
			stream.Send(&pb.JoinRoomResponse{
				Status: pb.JoinRoomResponse_MATCHED,
				Room:   build.PBRoom(room),
				Me:     build.PBPlayer(room.Guest),
			})
			h.Unlock()
			fmt.Printf("matchedroom_id=%v\n", room.ID)
			return nil
		}
	}

	// 空き部屋がない場合、ホストとして部屋作成
	me.Color = game.Black
	room := &game.Room{
		ID:   int32(len(h.Rooms)) + 1,
		Host: me,
	}
	h.Rooms[room.ID] = room
	h.Unlock()

	stream.Send(&pb.JoinRoomResponse{
		Status: pb.JoinRoomResponse_WAITING,
		Room:   build.PBRoom(room),
	})

	ch := make(chan int)

	go func(ch chan<- int) {
		for {
			h.RLock()
			guest := room.Guest
			h.RUnlock()

			if guest != nil {
				stream.Send(&pb.JoinRoomResponse{
					Status: pb.JoinRoomResponse_MATCHED,
					Room:   build.PBRoom(room),
					Me:     build.PBPlayer(room.Host),
				})

				ch <- 0
				break
			}
			time.Sleep(1 * time.Second)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ch)

	select {
	case <-ch:
	case <-ctx.Done():
		return status.Errorf(codes.DeadlineExceeded, "マッチング失敗...")
	}

	return nil
}
