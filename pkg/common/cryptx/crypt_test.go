package cryptx

import (
	"bytes"
	"encoding/hex"
	"testing"
)

const (
	SALT = "114514"
)

func TestCompareRandom(t *testing.T) {
	type args struct {
		encrypt string
		input   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"$2a$10$Mi2u4xuAECbYFdhsUDYxgeuOj1zLJd4awYBz5tfrFbSfQW.1PEYIi", "1145141919810"}, true},
		{"2", args{"$2a$10$Mi2u4xuAECbYFdhsUDYxgeuOj1zLJd4awYBz5tfrFbSfQW.1PEYIi", "114514"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareRandom(tt.args.encrypt, tt.args.input); got != tt.want {
				t.Errorf("CompareRandom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareSalt(t *testing.T) {
	type args struct {
		encrypt string
		input   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"e042de6b56b275ad2d463fdd5e5adeaab407466d9894de0f09547a0c24aaa31d", "1919810"}, true},
		{"2", args{"e042de6b56b275ad2d463fdd5e5adeaab407466d9894de0f09547a0c24aaa31d", "114514"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareSalt(SALT, tt.args.encrypt, tt.args.input); got != tt.want {
				t.Errorf("CompareSalt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCompareSHA256(t *testing.T) {
	type args struct {
		encrypt string
		input   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"1", args{"384fde3636e6e01e0194d2976d8f26410af3e846e573379cb1a09e2f0752d8cc", "114514"}, true},
		{"2", args{"384fde3636e6e01e0194d2976d8f26410af3e846e573379cb1a09e2f0752d8cc", "1919810"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CompareSHA256(tt.args.encrypt, tt.args.input); got != tt.want {
				t.Errorf("CompareSHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkCompareRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareRandom("$2a$10$Mi2u4xuAECbYFdhsUDYxgeuOj1zLJd4awYBz5tfrFbSfQW.1PEYIi", "1145141919810")
		CompareRandom("$2a$10$Mi2u4xuAECbYFdhsUDYxgeuOj1zLJd4awYBz5tfrFbSfQW.1PEYIi", "114514")
	}
}

func BenchmarkCompareSalt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareSalt(SALT, "e042de6b56b275ad2d463fdd5e5adeaab407466d9894de0f09547a0c24aaa31d", "1919810")
		CompareSalt(SALT, "e042de6b56b275ad2d463fdd5e5adeaab407466d9894de0f09547a0c24aaa31d", "114514")
	}
}

func BenchmarkCompareSHA256(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CompareSHA256("384fde3636e6e01e0194d2976d8f26410af3e846e573379cb1a09e2f0752d8cc", "114514")
		CompareSHA256("384fde3636e6e01e0194d2976d8f26410af3e846e573379cb1a09e2f0752d8cc", "1919810")
	}
}

func BenchmarkCryptAES(b *testing.B) {
	key := bytes.Repeat([]byte{1}, 16)
	orig := []byte("1919810")
	encrypted, _ := hex.DecodeString("3d705a9ed120d0b2478758678d8cb85c")
	for i := 0; i < b.N; i++ {
		EncryptAES(key, orig)
		DecryptAES(key, encrypted)
	}
}

func BenchmarkCryptXOR(b *testing.B) {
	key := bytes.Repeat([]byte{1}, 16)
	orig := []byte("1919810")
	encrypted, _ := hex.DecodeString("30383038393031080808080808080808")
	for i := 0; i < b.N; i++ {
		EncryptXOR(key, orig)
		DecryptXOR(key, encrypted)
	}
}
