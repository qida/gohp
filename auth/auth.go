package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Auth struct {
	User  interface{}
	Xxtea *XxTea
	tm    time.Time
	salt  int
}

const (
	f分隔符 = ":"
)

func NewAuth(key string) *Auth {
	return &Auth{Xxtea: NewXxTea(key)}
}

//加密
func (a *Auth) Encrypt(uid int, s interface{}) string {
	a.salt = rand.Intn(1000)
	a.tm = time.Now()
	a.User = s
	src, _ := json.Marshal(a)
	encodeString := base64.URLEncoding.EncodeToString(a.Xxtea.Encrypt(false, src))
	return fmt.Sprintf("%d%s%s", uid, f分隔符, encodeString)
}

//解密
func (a *Auth) Decrypt(uid_auth string) (userAuth Auth, err error) {
	if strings.Contains(uid_auth, f分隔符) {
		uid_auth = strings.SplitN(uid_auth, f分隔符, 2)[1]
	}
	var decodeBytes []byte
	if decodeBytes, err = base64.URLEncoding.DecodeString(uid_auth); err == nil {
		err = json.Unmarshal(a.Xxtea.Decrypt(false, decodeBytes), &userAuth)
	}
	return
}
