package models

type UserToken struct {
	userId uint64
	token [32]byte
}