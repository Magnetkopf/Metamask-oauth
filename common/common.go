package common

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"sync"
)

type LoginRequest struct {
	Address   string
	Nonce     string
	Timestamp int64
}

func GenerateNonce(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

var (
	LoginRequests = make(map[string]LoginRequest)
	Mu            sync.Mutex
	Logger        *log.Logger
)
