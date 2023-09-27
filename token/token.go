package token

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type TokenJWTInfo struct {
	Uid uint64 `json:"uid,omitempty"`
	Exp int64  `json:"exp,omitempty"`
}

var jWTCache = expirable.NewLRU[string, *TokenJWTInfo](30000, nil, -1)
var JWT_SECRET_KEY []byte

func InitJWTSECRETKEY() {
	JWT_SECRET_KEY = []byte(os.Getenv("JWT_SECRET_KEY"))
}
func GenerateToken(user_id uint) (string, error) {

	token_lifespan, err := strconv.Atoi(os.Getenv("TOKEN_HOUR_LIFESPAN"))

	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{}
	claims["uid"] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(token_lifespan)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))

}

func TokenValid(c *gin.Context) error {
	tokenString := ExtractToken(c)
	// fmt.Println("tokenString TokenValid", tokenString)
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func TokenStringisValid(tokenString string) error {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return err
	}
	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}

func ExtractTokenJWT(tokenString string) *TokenJWTInfo {
	if len(tokenString) == 0 {
		return nil
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return JWT_SECRET_KEY, nil
	})
	if err != nil {
		return nil
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil
	}
	var tokenInfo TokenJWTInfo
	tokenInfo.Uid, err = strconv.ParseUint(fmt.Sprintf("%v", claims["uid"]), 10, 64)
	if err != nil {
		return nil
	}
	exp, err := strconv.ParseFloat(fmt.Sprintf("%v", claims["exp"]), 64)
	if err != nil {
		return nil
	}
	tokenInfo.Exp = int64(exp)
	return &tokenInfo
}

func GetTokenJWT(c *gin.Context) *TokenJWTInfo {

	tokenString := ExtractToken(c)
	if len(tokenString) == 0 {
		return nil
	}
	tokenInfo, _ := jWTCache.Get(tokenString)
	if tokenInfo != nil {
		return tokenInfo
	}
	tokenInfo = ExtractTokenJWT(tokenString)
	if tokenInfo != nil {
		jWTCache.Add(tokenString, tokenInfo)
	}
	return tokenInfo
}

func ExtractTokenID(c *gin.Context) (uint, error) {

	tokenString := ExtractToken(c)
	// fmt.Println("tokenString ExtractTokenID", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		uid, err := strconv.ParseUint(fmt.Sprintf("%.0f", claims["uid"]), 10, 32)
		if err != nil {
			uid, err = strconv.ParseUint(fmt.Sprintf("%s", claims["uid"]), 10, 32)
			if err != nil {
				return 0, err
			}
		}
		return uint(uid), nil
	}
	return 0, nil
}

func ExtractTokenExpire(c *gin.Context) string {

	tokenString := ExtractToken(c)
	// fmt.Println("tokenString ExtractTokenExpire", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return ""
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {

		var tm time.Time
		switch iat := claims["exp"].(type) {
		case float64:
			tm = time.Unix(int64(iat), 0)
		case json.Number:
			v, _ := iat.Int64()
			tm = time.Unix(v, 0)
		}

		exp := fmt.Sprintf("%s", tm)

		return exp
	}
	return ""
}

func ExtractTokenExpireByInt(c *gin.Context) int64 {

	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET_KEY")), nil
	})
	if err != nil {
		return 0
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		var tm = int64(0)
		switch iat := claims["exp"].(type) {
		case float64:
			tm = int64(iat)
		case json.Number:
			v, _ := iat.Int64()
			tm = v
		}

		return tm
	}
	return 0
}
