package iteranal

import (
	"context"
	"strconv"

	"encoding/hex"
	"log/slog"

	"math/rand"
	rand2 "math/rand"
	"net/http"
	"strings"
	"time"
)

func GetClientIP(r *http.Request) string {
	// Try getting IP from X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain a list of IPs, take the first one
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}
	// If X-Forwarded-For is not available, get the remote address
	return r.RemoteAddr
}

func Mildwary(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//fmt.Printf("ip:%s\n", GetClientIP(r))
		//fmt.Printf("Method %s\n", r.Method)
		//fmt.Printf("URL %s\n", r.URL)

		next.ServeHTTP(w, r)

	})
}

func Contexte() (ctx context.Context, cancel context.CancelFunc) {

	return context.WithTimeout(context.Background(), 2*time.Second)

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

func Num() string {

	rand2.Seed(time.Now().UnixNano())

	random := rand.Intn(1000) + 9000
	r := strconv.Itoa(random)

	return r

}
