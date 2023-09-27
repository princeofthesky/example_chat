package transport

import (
	"github.com/gin-gonic/gin"
	"github.com/princeofthesky/example_chat/token"
	"github.com/princeofthesky/example_chat/trace_log"
)

func (hdl *skyHandler) OpenConnection(c *gin.Context) {

	username := c.Param("userId")
	userConnectionId := c.Param("userConnectionId")
	trace_log.Logger.Info().Str("username", username).Msg("open connection ")
	jwtTokenObject, ok := c.Get("jwt_info")
	var tokenInfo *token.TokenJWTInfo
	if ok && jwtTokenObject != nil {
		tokenInfo, ok = jwtTokenObject.(*token.TokenJWTInfo)
	}
	if tokenInfo == nil {
		tokenInfo = &token.TokenJWTInfo{
			Uid: 0,
			Exp: 0,
		}
	}
	hdl.repo.CreateConnection(c.Writer, c.Request, username, userConnectionId, tokenInfo)
}
