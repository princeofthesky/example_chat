package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/princeofthesky/example_chat/token"
	"github.com/princeofthesky/example_chat/trace_log"
)

func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenInfo := token.GetTokenJWT(c)
		if tokenInfo == nil {
			trace_log.Logger.Info().Msg("error when parse jwt token")
		}
		c.Set("jwt_info", tokenInfo)
		c.Next()
	}
}
