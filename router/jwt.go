package router

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"github.com/raozhaofeng/beego/utils"
	"net/http"
	"strings"
	"time"
)

// Generate 生成Token
func (c *Token) Generate(rds redis.Conn, tokenKey string, adminId, userId int64) string {
	tokenParams := TokenManager.GetTokenParams(rds, tokenKey)

	nowTime := time.Now()
	claims := &Claims{
		UserId:  userId,
		AdminId: adminId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: nowTime.Add(tokenParams.Expire * time.Second).Unix(), //	过期时间
			IssuedAt:  nowTime.Unix(),                                       //	签发时间
		},
	}

	//	生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(tokenParams.Key))
	if err != nil {
		panic(err)
	}

	c.SetTokenValue(rds, adminId, userId, tokenStr)
	return tokenStr
}

// Verify 验证Token
func (c *Token) Verify(rds redis.Conn, r *http.Request) *Claims {
	claims := &Claims{}

	tokenStr, tokenKey := c.GetHeaderTokenAndTokenKey(r)
	if tokenStr == "" || tokenKey == "" {
		return nil
	}

	// 验证JWT TokenStr
	tokenParams := c.GetTokenParams(rds, tokenKey)
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		token.Method = jwt.SigningMethodHS256
		return []byte(tokenParams.Key), nil
	})

	if err != nil || token == nil || !token.Valid {
		return nil
	}

	//	判断是否唯一Token
	if tokenParams.Only && c.GetTokenValue(rds, claims.AdminId, claims.UserId) != tokenStr {
		return nil
	}

	//	判断是否白名单
	userRealIP := utils.GetUserRealIP(r)
	if tokenParams.Whitelist != "" && utils.SliceStringIndexOf(userRealIP, strings.Split(tokenParams.Whitelist, ",")) == -1 {
		return nil
	}

	//	判断是否黑名单
	if tokenParams.Blacklist != "" && utils.SliceStringIndexOf(userRealIP, strings.Split(tokenParams.Blacklist, ",")) > -1 {
		return nil
	}

	return claims
}
