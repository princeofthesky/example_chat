package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/princeofthesky/example_chat/token"
	"math"
)

func (hdl *skyHandler) OpenConnectionPublic(c *gin.Context) {
	username := c.Param("id")
	tokenInfo := &token.TokenJWTInfo{
		Uid: 0,
		Exp: math.MaxInt64,
	}

	hdl.repo.CreateConnection(c.Writer, c.Request, username, "", tokenInfo)
}
