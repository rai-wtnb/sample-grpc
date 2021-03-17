package handler

import (
	"fmt"
	"reversi/game"
	"reversi/gen/pb"
	"sync"

	"reversi/build"
)

type GameHandler struct {
	sync.RWMutex
	games  map[int32]*game.Game
	client map[int32][]pb.GameService_PlayServer // 状態変更時にクライアントにストリーミングを返すように格納
}

func NewGameHandler() *GameHandler {
	return &GameHandler{
		games:  make(map[int32]*game.Game),
		client: make(map[int32][]pb.GameService_PlayServer),
	}
}

func (h *GameHandler) Play(stream pb.GameService_PlayServer) error {
	for {
		// clientからreqが来たら代入
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		roomID := req.GetRoom()
		player := build.Player(req.GetPlayer())

		// oneofで指定されたものをswitchで処理
		switch req.GetAction().(type) {
		case *pb.PlayRequest_Start: // ゲーム開始
			err := h.start(stream, roomID, player)
			if err != nil {
				return err
			}
		case *pb.PlayRequest_Move: // 石を置く
			action := req.GetMove()
			x := action.GetMove().GetX()
			y := action.GetMove().GetY()
			err := h.move(roomID, x, y, player)
			if err != nil {
				return err
			}
		}
	}
}

func (h *GameHandler) start(
	stream pb.GameService_PlayServer,
	roomID int32,
	me *game.Player,
) error {
	h.Lock()
	defer h.Unlock()

	// ゲーム情報がない場合作成
	g := h.games[roomID]
	if g == nil {
		g = game.NewGame(game.None)
		h.games[roomID] = g
		h.client[roomID] = make([]pb.GameService_PlayServer, 0, 2)
	}

	// 自分のclientを登録
	h.client[roomID] = append(h.client[roomID], stream)

	// 2人プレイヤーが揃ったのでゲーム開始
	if len(h.client[roomID]) == 2 {
		for _, s := range h.client[roomID] {
			err := s.Send(&pb.PlayResponse{
				Event: &pb.PlayResponse_Ready{
					Ready: &pb.PlayResponse_ReadyEvent{},
				},
			})
			if err != nil {
				return err
			}
		}
		fmt.Printf("game has started room_id=%v\n", roomID)
	} else {
		// 待機中であることをクライアントに通知
		err := stream.Send(&pb.PlayResponse{
			Event: &pb.PlayResponse_Waiting{
				Waiting: &pb.PlayResponse_WaitingEvent{},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *GameHandler) move(roomID, x, y int32, p *game.Player) error {
	h.Lock()
	defer h.Unlock()

	g := h.games[roomID]

	finished, err := g.Move(x, y, p.Color)
	if err != nil {
		return err
	}

	for _, s := range h.client[roomID] {
		// 石が置かれたことをクライアントに通知
		err := s.Send(&pb.PlayResponse{
			Event: &pb.PlayResponse_Move{
				Move: &pb.PlayResponse_MoveEvent{
					Player: build.PBPlayer(p),
					Move: &pb.Move{
						X: x,
						Y: y,
					},
					Board: build.PBBoard(g.Board),
				},
			},
		})
		if err != nil {
			return err
		}

		if finished { //ゲーム終了
			err := s.Send(
				&pb.PlayResponse{
					Event: &pb.PlayResponse_Finished{
						Finished: &pb.PlayResponse_FinishedEvent{
							Winner: build.PBColor(g.Winner()),
							Board:  build.PBBoard(g.Board),
						},
					},
				},
			)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
