package jwt

import (
	"encoding/json"
	"errors"
	"github.com/mendsley/gojwk"
	"io/ioutil"
	"net/http"
)

type Key struct {
	E string `json:"e,omitempty"` //  ep. AQAB
	N string `json:"n,omitempty"` // RSA key 内容

	Crv string `json:"crv,omitempty"` // EC key 组成
	X   string `json:"x,omitempty"`   // EC key 组成
	Y   string `json:"y,omitempty"`   // EC key 组成

	Kty string `json:"kty"`           // ket type ep. RSA/EC
	Alg string `json:"alg"`           // algorithm  ep. RS256
	Use string `json:"use,omitempty"` // sig -> signature(签名) enc -> encryption(加密)
	Kid string `json:"kid"`           // key id
}

// Keys 只支持 PublicKey
type Keys struct {
	Keys []Key `json:"keys"`
}

// ParsePathToKey 将 wellknown 中的 jwk 转换成 public key
// 支持 EC/RSA
func ParsePathToKey(path string, kid string) (interface{}, Key, error) {
	resp, err := http.Get(path)
	if err != nil {
		return nil, Key{}, errors.New("cannot get resp, maybe it does not exist")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, Key{}, errors.New("error while reading resp.Body")
	}
	var keys = &Keys{}
	err = json.Unmarshal(body, keys)
	if err != nil {
		return nil, Key{}, err
	}
	for _, key := range keys.Keys {
		if key.Kid == kid {
			jwk, _ := json.Marshal(key)
			k, err := gojwk.Unmarshal(jwk)
			if err != nil {
				return nil, Key{}, err
			}
			publicKey, err := k.DecodePublicKey()
			return publicKey, key, err
		}
	}
	return nil, Key{}, errors.New("cannot find this key")
}
