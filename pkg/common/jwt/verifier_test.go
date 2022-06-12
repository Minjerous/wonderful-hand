package jwt

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
)

var defaultPublicKey = `
-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA4f5wg5l2hKsTeNem/V41
fGnJm6gOdrj8ym3rFkEU/wT8RDtnSgFEZOQpHEgQ7JL38xUfU0Y3g6aYw9QT0hJ7
mCpz9Er5qLaMXJwZxzHzAahlfA0icqabvJOMvQtzD6uQv6wPEyZtDTWiQi9AXwBp
HssPnpYGIn20ZZuNlX2BrClciHhCPUIIZOQn/MmqTD31jSyjoQoV7MhhMTATKJx2
XrHhR+1DcKJzQBSTAGnpYVaqpsARap+nwRipr3nUTuxyGohBTSmjJ2usSeQXHI3b
ODIRe1AuTyHceAbewn8b462yEWKARdpd9AjQW5SIVPfdsz5B6GlYQ5LdYKtznTuy
7wIDAQAB
-----END PUBLIC KEY-----
`

var verifyData = []struct {
	tokenStr string
	key      interface{}
	useKid   bool
	alg      string
	wantErr  bool
}{
	{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM2MDE2NDcxNjMxOTgsImlhdCI6MTY0NzE2MzE5OH0.RAeU-lYrUj85Uuyp-d-o0EeqP4QnFIaqFTaX1tAeVJY",
		[]byte("114514"),
		false,
		"HS256",
		false,
	},
	{
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjM2MDE2NDcxNjMxOTgsImlhdCI6MTY0NzE2MzE5OH0.RAeU-lYrUj85Uuyp-d-o0EeqP4QnFIaqFTaX1tAeVJY",
		[]byte("1145141919810"),
		false,
		"HS256",
		true,
	},
	{
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.e30.FvuwHdEjgGxPAyVUb-eqtiPl2gycU9WOHNzwpFKcpdN_QkXkBUxU3qFl3lLBaMzIuP_GjXLXcJZFhyQ2Ne3kfWuZSGLmob0Og8B4lAy7CA7iwpji2R3aUcwBwbJ41IJa__F8fMRz0dRDwhyrBKD-9y4TfV_-yZuzBZxq0UdjX6IdpzsdetphBSIZkPij5MY3thRwC-X_gXyIXi4-G2_CjRrV5lCGnPJrDbLqPCYqS71wK9NEsz_B8p5ENmwad8vZe4fEFR7XsqJrhPjbEVGeLpzSz0AOGp4G1iyvv1sdu4M3Y8KSSGYnZ8lXNGyi8QeUr374Y6XgJ5N5TVLWI2cMxg",
		nil,
		false,
		"RS256",
		false,
	},
	{
		"eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IjF6aXlIVG15M184MDRDOU1jUENHVERmYWJCNThBNENlZG9Wa3VweXdVeU0ifQ.eyJqdGkiOiIzWUJ5eWZ2TDB4b01QNXdwTXVsZ0wiLCJzdWIiOiI2MDE5NDI5NjgwMWRjN2JjMmExYjI3MzUiLCJpYXQiOjE2MTI0NDQ4NzEsImV4cCI6MTYxMzY1NDQ3MSwic2NvcGUiOiJvcGVuaWQgZW1haWwgbWVzc2FnZSIsImlzcyI6Imh0dHBzOi8vc3RlYW0tdGFsay5hdXRoaW5nLmNuL29pZGMiLCJhdWQiOiI2MDE5M2M2MTBmOTExN2U3Y2IwNDkxNTkifQ.cYyZ6buwAjp7DzrYQEhvz5rvUBhkv_s8xzuv2JHgzYx0jbqqsWrA_-gufLTFGmNkZkZwPnF6ktjvPHFT-1iJfWGRruOOMV9QKPhk0S5L2eedtbKJU6XIEkl3F9KbOFwYM53v3E7_VC8RBj5IKqEY0qd4mW36C9VbS695wZlvMYnmXhIopYsd5c83i39fLBF8vEBZE1Rq6tqTQTbHAasR2eUz1LnOqxNp2NNkV2dzlcNIksSDbEGjTNkWceeTWBRtFMi_o9EWaHExdm5574jQ-ei5zE4L7x-zfp9iAe8neuAgTsqXOa6RJswhyn53cW4DwWg_g26lHJZXQvv_RHZRlQ",
		nil,
		true,
		"RS256",
		false,
	},
}

func TestVerify(t *testing.T) {
	option := VerOption{}
	verifier := NewVerifier(option.WithSkipVerifyClaim()) // 不验证 Claims 的有效性 (即 是否过期|是否启用)
	for _, datum := range verifyData {
		// 没有使用 kid
		if datum.useKid == false {
			if datum.key == nil {
				datum.key, _ = jwt.ParseRSAPublicKeyFromPEM([]byte(defaultPublicKey))
			}
			err, _ := verifier.Verify(datum.tokenStr, datum.key)
			if (err != nil) != datum.wantErr {
				t.Errorf("error while verify %v", datum)
			}
			continue
		}
		// 使用 kid，无需指定 key
		err, _ := verifier.VerifyWithKid(datum.tokenStr)
		if (err != nil) != datum.wantErr {
			t.Errorf("error while verify %v", datum)
		}
	}
}

func BenchmarkVerify(b *testing.B) {
	option := VerOption{}
	verifier := NewVerifier(option.WithSkipVerifyClaim())
	for i := 0; i < b.N; i++ {
		err, _ := verifier.Verify(verifyData[0].tokenStr, verifyData[0].key)
		if err != nil {
			panic(err)
		}
	}
}
