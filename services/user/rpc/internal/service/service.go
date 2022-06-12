package service

import (
	"errors"
	"wonderful-hand-common/jwt"
	"wonderful-hand-user/rpc/internal/config"
	"wonderful-hand-user/rpc/user"
)

type UserService struct {
	user.UnimplementedUserServiceServer
	tokenGen *jwt.Generator
	tokenVer *jwt.Verifier
	cfg      config.Config
}

func MewUserService(cfg config.Config) *UserService {
	genOpts := jwt.GenOption{}
	verOpts := jwt.VerOption{}
	return &UserService{
		tokenGen: jwt.NewGenerator(genOpts.WithKey([]byte(cfg.Auth.AccessSecret), jwt.SigningMethodHS256)),
		tokenVer: jwt.NewVerifier(verOpts.WithDefaultKey([]byte(cfg.Auth.AccessSecret))),
		cfg:      cfg,
	}
}

var (
	passwordNotMatched = errors.New("wrong password")
	accountExist       = errors.New("account exist")
)

type tokenClaims struct {
	jwt.StandClaims
	IsAccess bool `json:"is_access"`
}
