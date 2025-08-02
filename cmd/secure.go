package cmd

import (
	"Project2/iteranal"
	"context"
	"crypto/tls"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	gomail "gopkg.in/mail.v2"
	"log/slog"
	"math/rand"
	"net/http"
	"strconv"
	_ "strconv"
	"time"
)

type SecureHandler struct {
	DB *sql.DB
}

type HandlerSecure struct {
	Email string `json:"email"`
}

func (h *SecureHandler) Secure(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		slog.Error("func secure 1 : dont allow this method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return

	}

	session, err := store.Get(r, "token1")
	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "Cookie dont set", http.StatusUnauthorized)

		return
	}

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	var p HandlerSecure

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		slog.Error("func secure 2 : dont allow this method")
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()

	var two bool
	err = h.DB.QueryRowContext(ctx, "SELECT two_factor_enabled FROM PERSON WHERE email = $1", p.Email).Scan(&two)

	if errors.Is(err, context.DeadlineExceeded) {
		slog.Error("func secure 3: context timed out", err)
		http.Error(w, "context timed out ", http.StatusBadRequest)
		return

	} else if errors.Is(err, sql.ErrNoRows) {
		slog.Error("Func secer3 :context timed out ", err)
		http.Error(w, "errr", http.StatusBadRequest)
		return
	}
	if !two {

		slog.Info("func secure 4:  Data is write in request and dont write")
	} else {

		slog.Info("secure is set")
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"twofaEnabled": true,
		})
		return

	}
	rand.Seed(time.Now().UnixNano())

	random := rand.Intn(1000) + 9000
	re := strconv.Itoa(random)

	slog.Info("ааааааа", random)

	session.Values["code"] = re
	session.Values["Email"] = p.Email

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3000,
		Secure:   false,
		HttpOnly: true,
	}

	session.Save(r, w)
	m := gomail.NewMessage()
	m.SetHeader("From", "frederickscubi@gmail.com")
	m.SetHeader("To", p.Email)
	m.SetHeader("Subject", "YOUR CODE ")
	m.SetBody("text/plain", re)
	d := gomail.NewDialer("smtp.gmail.com", 587, "frederickscubi@gmail.com", "kxgehpuwmvgjmtev")
	d.TLSConfig = &tls.Config{ServerName: "smtp.gmail.com"}
	if err := d.DialAndSend(m); err != nil {
		slog.Info("func set secure 6 ", "err", err)
		return

	}
	slog.Info("email sent")
	w.Header().Set("Content-Type", "application/json")

	resp := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	json.NewEncoder(w).Encode(resp)

}
