package token

import (
	"github.com/golang-jwt/jwt/v4"
)

// Sign 使用 HS256 将 claims 签名为 JWT 字符串。
// claims 传 map[string]any 即可（等价于 jwt.MapClaims），由调用方自行决定过期时间等字段。
func Sign(secret string, claims map[string]any) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims(claims))
	return t.SignedString([]byte(secret))
}
