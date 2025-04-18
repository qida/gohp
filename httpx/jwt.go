package httpx

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserInfo interface{}
	jwt.RegisteredClaims
}

// 签名密钥
const sign_key = "hello jwt"

func randStr(strLen int) string {
	bytes := make([]byte, strLen)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:strLen]
}

// 生成 JWT
func JWTGenerate(user_info interface{}) (string, error) {
	claim := Claims{
		UserInfo: user_info,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者
			Subject:   "Tom",                                           // 签发对象
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      // 签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),   // 过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), // 最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  // 签发时间
			ID:        randStr(10),                                     // wt ID, 类似于盐值
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(sign_key))
	return token, err
}

// 验证 JWT
func JWTValidate(token_string string) (Claims, error) {
	token, err := jwt.ParseWithClaims(
		token_string,
		&Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(sign_key), nil // 返回签名密钥
		},
	)
	if err != nil {
		return Claims{}, err
	}

	if !token.Valid {
		return Claims{}, errors.New("claim invalid")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return Claims{}, errors.New("invalid claim type")
	}

	return *claims, nil
}
