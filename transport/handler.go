package transport

import "github.com/princeofthesky/example_chat/repository"

type skyHandler struct {
	repo *repository.SocketLive
}

func NewSkyHandler(repo *repository.SocketLive) *skyHandler {
	return &skyHandler{
		repo: repo,
	}
}
