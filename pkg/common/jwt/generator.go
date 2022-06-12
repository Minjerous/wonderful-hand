package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"strings"
)

var (
	SigningMethodES256 = jwt.SigningMethodES256 // 圆锥曲线 256
	SigningMethodES384 = jwt.SigningMethodES384 // 圆锥曲线 384
	SigningMethodES512 = jwt.SigningMethodES512 // 圆锥曲线 512

	SigningMethodRS256 = jwt.SigningMethodRS256 // RSA 256
	SigningMethodRS384 = jwt.SigningMethodRS384 // RSA 384
	SigningMethodRS512 = jwt.SigningMethodRS512 // RSA 512

	SigningMethodHS256 = jwt.SigningMethodHS256
	SigningMethodHS384 = jwt.SigningMethodHS384
	SigningMethodHS512 = jwt.SigningMethodHS512

	SigningMethodPS256 = jwt.SigningMethodPS256
	SigningMethodPS384 = jwt.SigningMethodPS384
	SigningMethodPS512 = jwt.SigningMethodPS512
)

type Generator struct {
	t         *jwt.Token
	key       interface{}
	globalIss string
}

type GenOption struct {
	ops []generatorOption
}

type generatorOption func(generator *Generator) *Generator

// WithKey 为jwt生成器添加一个key
// 参数
//	key: 	密钥(对称加密密钥 []byte 或者 非对称加密私钥 *rsa.PrivateKey)
//	method: 加密方式
func (g *GenOption) WithKey(key interface{}, method jwt.SigningMethod) *GenOption {
	g.ops = append(g.ops, func(generator *Generator) *Generator {
		generator.key = key
		generator.t = jwt.New(method)
		return generator
	})
	return g
}

// WithKeyId 一般用于非对称加密，jwk 中不能缺少 alg
// kid 指向的 key 是公钥，key 参数是私钥
// 不要和 WithKey 一起用!
func (g *GenOption) WithKeyId(kid, iss string, key interface{}) *GenOption {
	// 检查验证key的 kid 是否存在
	if iss == "" {
		return g
	}
	if !strings.HasSuffix(iss, "/") {
		iss = fmt.Sprintf("%s%s", iss, "/")
	}
	path := fmt.Sprintf("%s%s", iss, JWKsRelativePath)
	_, k, err := ParsePathToKey(path, kid) // 检查链接是否可以连接解析
	if err != nil {
		return g
	}
	g.ops = append(g.ops, func(generator *Generator) *Generator {
		generator.key = key // 签名 key
		generator.t = jwt.New(jwt.GetSigningMethod(k.Alg))
		generator.t.Header["kid"] = kid
		generator.globalIss = iss
		return generator
	})
	return g
}

func NewGenerator(opt *GenOption) *Generator {
	g := &Generator{}
	for _, option := range opt.ops {
		g = option(g)
	}
	return g
}

func (g *Generator) Sign(claims jwt.Claims) (string, error) {
	g.t.Claims = claims
	if g.globalIss != "" {
		switch claims.(type) {
		case jwt.MapClaims:
			c := claims.(jwt.MapClaims)
			c["iss"] = g.globalIss
			g.t.Claims = c
		case jwt.StandardClaims:
			c := claims.(jwt.StandardClaims)
			c.Issuer = g.globalIss
			g.t.Claims = c
		default:
			return "", errors.New("cannot support your own claims")
		}
	}
	return g.t.SignedString(g.key)
}

func (g *Generator) GetSigning(claims jwt.Claims) (string, error) {
	g.t.Claims = claims
	if g.globalIss != "" {
		switch claims.(type) {
		case jwt.MapClaims:
			c := claims.(jwt.MapClaims)
			c["iss"] = g.globalIss
			g.t.Claims = c
		case jwt.StandardClaims:
			c := claims.(jwt.StandardClaims)
			c.Issuer = g.globalIss
			g.t.Claims = c
		default:
			return "", errors.New("cannot support your own claims")
		}
	}
	return g.t.SigningString()
}
