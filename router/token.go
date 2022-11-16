package router

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gomodule/redigo/redis"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	TokenParamsRedisName                 = "_tokenParams"      //	Token参数缓存名称
	TokenValuesRedisName                 = "_tokenValues"      //	Token值缓存名称
	TokenAdminRolesRouterName            = "_adminRolesRouter" //	管理角色路由缓存名称
	ClaimsKey                 contextKey = 1                   //	claims类型
)

// TokenManager Token管理
var TokenManager *Token

type contextKey int

// Claims Token对象
type Claims struct {
	UserId             int64 //	用户ID
	AdminId            int64 //	管理员ID
	jwt.StandardClaims       //	jwt基础参数
}

// TokenParams Token 参数
type TokenParams struct {
	Key       string        `json:"key"`       //	密钥
	Only      bool          `json:"only"`      //	是否唯一
	Expire    time.Duration `json:"expire"`    //	过期时间
	Whitelist string        `json:"whitelist"` //	白名单
	Blacklist string        `json:"blacklist"` //	黑名单
}

type Token struct {
}

// InitializationToken 初始化Token
func InitializationToken(rds redis.Conn, tokenParamsList map[string]*TokenParams, adminRolesRouter map[int64][]string) {
	TokenManager = &Token{}
	//	初始化设置Token参数
	for tokenKey, tokenParams := range tokenParamsList {
		TokenManager.SetTokenParams(rds, tokenKey, tokenParams)
	}

	// 初始化设置管理路由列表
	for adminId, rolesRouter := range adminRolesRouter {
		TokenManager.SetTokenAdminRolesRouter(rds, adminId, rolesRouter)
	}
}

// GetTokenParams 获取Token参数
func (c *Token) GetTokenParams(rds redis.Conn, tokenKey string) *TokenParams {
	tokenParamsBytes, err := redis.Bytes(rds.Do("HGET", TokenParamsRedisName, tokenKey))
	if err != nil {
		return nil
	}
	tokenParams := new(TokenParams)
	_ = json.Unmarshal(tokenParamsBytes, &tokenParams)
	return tokenParams
}

// SetTokenParams 设置Token参数
func (c *Token) SetTokenParams(rds redis.Conn, tokenKey string, tokenParams *TokenParams) {
	tokenParamsBytes, _ := json.Marshal(tokenParams)
	_, _ = rds.Do("HSET", TokenParamsRedisName, tokenKey, tokenParamsBytes)
}

// GetTokenValue 获取Token值
func (c *Token) GetTokenValue(rds redis.Conn, adminId, userId int64) string {
	tokenStr, _ := redis.String(rds.Do("HGET", TokenValuesRedisName, c.GetTokenValueKey(adminId, userId)))
	return tokenStr
}

// SetTokenValue 设置Token值
func (c *Token) SetTokenValue(rds redis.Conn, adminId, userId int64, tokenStr string) {
	_, _ = rds.Do("HSET", TokenValuesRedisName, c.GetTokenValueKey(adminId, userId), tokenStr)
}

// GetTokenAdminRolesRouter 获取管理角色路由列表
func (c *Token) GetTokenAdminRolesRouter(rds redis.Conn, adminId int64) []string {
	adminRolesRouter, err := redis.String(rds.Do("HGET", TokenAdminRolesRouterName, adminId))
	if err != nil {
		return []string{}
	}
	return strings.Split(adminRolesRouter, ",")
}

// SetTokenAdminRolesRouter 设置管理角色路由列表
func (c *Token) SetTokenAdminRolesRouter(rds redis.Conn, adminId int64, rolesRouter []string) {
	_, _ = rds.Do("HSET", TokenAdminRolesRouterName, adminId, strings.Join(rolesRouter, ","))
}

// AuthRouter 验证路由
func (c *Token) AuthRouter(rds redis.Conn, adminId int64, router string) bool {
	adminRolesRouter := c.GetTokenAdminRolesRouter(rds, adminId)
	fmt.Println(adminRolesRouter)
	for _, adminRouter := range adminRolesRouter {
		if router == adminRouter {
			return true
		}
	}
	return false
}

// GetTokenValueKey 获取Token值key
func (c *Token) GetTokenValueKey(adminId, userId int64) string {
	adminIdStr := strconv.FormatInt(adminId, 10)
	userIdStr := strconv.FormatInt(userId, 10)
	return adminIdStr + "_" + userIdStr
}

// GetContextClaims 获取当前Claims
func (c *Token) GetContextClaims(r *http.Request) *Claims {
	return r.Context().Value(ClaimsKey).(*Claims)
}

// GetHeaderTokenAndTokenKey 获取头信息Token参数
func (c *Token) GetHeaderTokenAndTokenKey(r *http.Request) (string, string) {
	tokenStr := r.Header.Get("Token")
	tokenKey := r.Header.Get("Token-Key")

	urlTokenStr := r.URL.Query().Get("token")
	urlTokenKey := r.URL.Query().Get("key")
	if urlTokenStr != "" {
		tokenStr = urlTokenStr
	}
	if urlTokenKey != "" {
		tokenKey = urlTokenKey
	}
	return tokenStr, tokenKey
}
