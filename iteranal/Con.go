package iteranal

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"time"
)

func Contexte() (ctx context.Context) {

	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)

	return ctx

}

func Generateid() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		slog.Error("Error generating id", err)
		return "", err

	}
	return hex.EncodeToString(b), nil

}
