package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"reflect"
	"strings"
)

const (
	JWKsRelativePath = ".well-known/jwks.json"
)

type Token = jwt.Token

type StandClaims = jwt.StandardClaims

type MapClaims = jwt.MapClaims

type Claims interface {
	Valid() error
	VerifyAudience(cmp string, req bool) bool
	VerifyExpiresAt(cmp int64, req bool) bool
	VerifyIssuedAt(cmp int64, req bool) bool
	VerifyIssuer(cmp string, req bool) bool
	VerifyNotBefore(cmp int64, req bool) bool
}

type Verifier struct {
	p                   *jwt.Parser
	defaultKey          interface{}
	skipClaimVerify     bool
	useJSONNumber       bool
	validSigningMethods []string
	useClaims           jwt.Claims // 验证时使用的 claims，可使用 struct 来限制 jwt 的格式
	verifyAud           string     // 验证 Aud	仅当 skipClaimVerify = false 才启用
	verifyIss           string     // 验证 Iss   仅当 skipClaimVerify = false 才启用
}

type VerOption struct {
	ops []verifierOption
}

type verifierOption func(verifier *Verifier) *Verifier

// WithAudience 验证 jwt 中的 aud
func (v *VerOption) WithAudience(aud string) *VerOption {
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.verifyAud = aud
		return verifier
	})
	return v
}

// WithDefaultKey 使用 AcceptNoneSigningOption 后不能使用这个
func (v *VerOption) WithDefaultKey(key interface{}) *VerOption {
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.defaultKey = key
		return verifier
	})
	return v
}

// WithIssuer 验证 jwt 中的 iss
func (v *VerOption) WithIssuer(iss string) *VerOption {
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.verifyIss = iss
		return verifier
	})
	return v
}

// WithMethods 允许 jwt 使用传入的签名方法
func (v *VerOption) WithMethods(methods ...jwt.SigningMethod) *VerOption {
	algs := make([]string, len(methods))
	for i, method := range methods {
		algs[i] = method.Alg()
	}
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.validSigningMethods = algs
		return verifier
	})
	return v
}

// WithStructClaim 使用 StructClaim
// 例子：
//	type MyClaim struct {
//		jwt.StandardClaims
//		Uid int64 `json:"uid,omitempty"`
//		Sid int64 `json:"sid,omitempty"`
// 	}
// 	opt := VerOption{}
//  jwt.NewVerifier(opt.WithStructClaim(MyClaim{}))
func (v *VerOption) WithStructClaim(claims jwt.Claims) *VerOption {
	if reflect.TypeOf(claims) == reflect.TypeOf(jwt.MapClaims{}) {
		return v
	}
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.useClaims = claims
		return verifier
	})
	return v
}

// WithSkipVerifyClaim 不验证 jwt Claim 的有效性和用户自定义 Verifier 对其他字段的要求
func (v *VerOption) WithSkipVerifyClaim() *VerOption {
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.skipClaimVerify = true
		return verifier
	})
	return v
}

func (v *VerOption) WithJSONumber() *VerOption {
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.useJSONNumber = true
		return verifier
	})
	return v
}

func (v *VerOption) WithMapClaim() *VerOption {
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.useClaims = jwt.MapClaims{}
		return verifier
	})
	return v
}

func (v *VerOption) WithNoneSigning() *VerOption {
	v.ops = append(v.ops, func(verifier *Verifier) *Verifier {
		verifier.defaultKey = jwt.UnsafeAllowNoneSignatureType
		return verifier
	})
	return v
}

func NewVerifier(opt *VerOption) *Verifier {
	v := &Verifier{}
	for _, option := range opt.ops {
		v = option(v)
	}
	initVerifier(v)
	return v
}

func initVerifier(verifier *Verifier) {
	verifier.p = &jwt.Parser{
		ValidMethods:         verifier.validSigningMethods,
		UseJSONNumber:        verifier.useJSONNumber,
		SkipClaimsValidation: verifier.skipClaimVerify,
	}

	if verifier.useClaims == nil {
		verifier.useClaims = jwt.MapClaims{}
	}
}

func (v *Verifier) verifyTokenStr(jwtStr string) (*jwt.Token, []string, error) {
	token, parts, err := v.p.ParseUnverified(jwtStr, v.useClaims)
	// Verify signing method is in the required set
	if v.validSigningMethods != nil {
		var signingMethodValid = false
		var alg = token.Method.Alg()
		for _, m := range v.validSigningMethods {
			if m == alg {
				signingMethodValid = true
				break
			}
		}
		if !signingMethodValid {
			// signing method is not in the listed set
			return token, parts, errors.New(fmt.Sprintf("signing method %v is invalid", alg))
		}
	}

	// Validate Claims
	if !v.skipClaimVerify {
		if err = token.Claims.Valid(); err != nil {
			return token, parts, err
		}
		claims, ok := token.Claims.(Claims)
		if !ok {
			return token, parts, errors.New("do not recommend you to design your own Claims, if you want, implement awesome_jwt.Claims instead")
		}
		if v.verifyAud != "" {
			if claims.VerifyAudience(v.verifyAud, true) == false {
				return token, parts, errors.New(fmt.Sprintf("aud verified failed, expect %v", v.verifyAud))
			}
		}

		if v.verifyIss != "" {
			if claims.VerifyIssuer(v.verifyIss, true) == false {
				return token, parts, errors.New(fmt.Sprintf("iss verified failed, expect %v", v.verifyIss))
			}
		}
	}

	return token, parts, nil
}

func (v *Verifier) Verify(jwtStr string, key interface{}) (error, *Token) {
	token, parts, err := v.verifyTokenStr(jwtStr)
	if err != nil {
		return err, nil
	}

	// 防止这里使用 jwt.UnsafeAllowNoneSignatureType
	if _, ok := key.(string); ok {
		return errors.New("key cannot be string"), nil
	}
	return token.Method.Verify(strings.Join(parts[0:2], "."), parts[2], key), token
}

// VerifyWithNonKey 不验证签名
func (v *Verifier) VerifyWithNonKey(jwtStr string) (error, *Token) {
	if v.defaultKey != jwt.UnsafeAllowNoneSignatureType {
		return errors.New("not none key"), nil
	}
	token, parts, err := v.verifyTokenStr(jwtStr)
	if err != nil {
		return err, nil
	}
	return token.Method.Verify(strings.Join(parts[0:2], "."), parts[2], v.defaultKey), token
}

func (v *Verifier) VerifyWithDefaultKey(jwtStr string) (error, *Token) {
	if v.defaultKey == nil {
		return errors.New("no default key"), nil
	}
	token, parts, err := v.verifyTokenStr(jwtStr)
	if err != nil {
		return err, nil
	}
	return token.Method.Verify(strings.Join(parts[0:2], "."), parts[2], v.defaultKey), token
}

func (v *Verifier) VerifyWithKid(jwtStr string) (error, *Token) {
	token, parts, err := v.verifyTokenStr(jwtStr)
	if err != nil {
		return err, nil
	}
	var kid interface{}
	var ok bool
	if kid, ok = token.Header["kid"]; !ok {
		return errors.New("cannot find kid in jwt header"), nil
	}
	keyId, ok := kid.(string)
	if !ok {
		return errors.New("kid cannot cast to string"), nil
	}
	var iss string
	if c, ok := token.Claims.(jwt.MapClaims); ok {
		iss, ok = c["iss"].(string)
		if !ok {
			return errors.New("iss cannot cast to string"), nil
		}
	} else {
		claims := token.Claims.(jwt.StandardClaims)
		iss = claims.Issuer
	}

	if iss == "" {
		return errors.New("iss cannot be empty"), nil
	}

	if !strings.HasSuffix(iss, "/") {
		iss = fmt.Sprintf("%s%s", iss, "/")
	}

	publicKey, _, err := ParsePathToKey(fmt.Sprintf("%s%s", iss, JWKsRelativePath), keyId)

	if err != nil {
		return err, nil
	}

	return token.Method.Verify(strings.Join(parts[0:2], "."), parts[2], publicKey), token
}
