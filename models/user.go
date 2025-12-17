package models

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/benjamin-larsen/NoctesChat-WebSocket/util"
)

type User struct {
	ID        uint64 `json:"id,string"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
}

type UserToken struct {
	UserId uint64
	Token [32]byte
}

var nullToken = UserToken{}

var ErrInvalidFormat = errors.New("UserToken: invalid format (must have 2 segments)")
var ErrInvalidTokenSize = errors.New("UserToken: invalid token size (must be 32 bytes)")

func UserTokenFromString(token string) (UserToken, error) {
	segments := strings.SplitN(token, ":", 2)

	if len(segments) != 2 {
		return nullToken, ErrInvalidFormat
	}

	rawUserId, err := base64.RawURLEncoding.DecodeString(segments[0])

	if err != nil {
		return nullToken, fmt.Errorf("UserToken: %w", err)
	}

	userId, err := strconv.ParseUint(string(rawUserId), 10, 64)

	if err != nil {
		return nullToken, fmt.Errorf("UserToken: %w", err)
	}

	plainToken, err := base64.RawURLEncoding.DecodeString(segments[1])

	if err != nil {
		return nullToken, fmt.Errorf("UserToken: %w", err)
	}

	if len(plainToken) != 32 {
		return nullToken, ErrInvalidTokenSize
	}

	userToken := UserToken{
		UserId: userId,
		Token: sha256.Sum256(plainToken),
	}

	// Wipe token bytes for security, so it does not remain in-memory
	util.WipeBytes(plainToken)

	return userToken, nil
}