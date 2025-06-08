package iteranal

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"
)

func Contexte() (ctx context.Context, cancel context.CancelFunc) {

	return context.WithTimeout(context.Background(), 5*time.Second)

}

func Generateid() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("Error generating id", err)
		return ""

	}
	return hex.EncodeToString(b)

}
