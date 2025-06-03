package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Handler_login struct {
	DB *sql.DB
}

type Handler_register struct {
	email    string `json:"email"`
	password string `json:"password"`
}

func (handler *Handler_login) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return

	}

	var err error

	var p Handler_register
	ctx := iteranal.Contexte()

	err = json.NewDecoder(r.Body).Decode(&handler)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var exists bool

	err = handler.DB.QueryRowContext(ctx, "SELECT exists(SELECT 2 FROM person WHERE email=$1 and passowrd=$2)", p.email, p.password).Scan(&exists)
	switch {
	case err == context.DeadlineExceeded:
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		slog.Info("func register1:timed out")
		return
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		slog.Info("func register2:no rows")
		return

	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func register3:no rows")
		return

	case !exists:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		slog.Info("func register4:no rows")
		return

	default:
		slog.Info("func register4:ok")

	}

	cookie, err := iteranal.Generateid()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func register5:err", err)
		return
	}

	expires := time.Now().Add(time.Hour * 24)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    cookie,
		Path:     "/",
		Expires:  expires,
		HttpOnly: false,
		Secure:   true,
	})

	slog.Info("func register6:ok")

	_, err = handler.DB.ExecContext(ctx, "UPDATE person SET session_id = $1  WHERE email = $2", cookie, p.email)

	switch {
	case err == context.DeadlineExceeded:
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		slog.Info("func register6:timed out")
		return
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		slog.Info("func register7:no rows")
		return

	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func register8:no rows")
		return

	default:
		slog.Info("func register9:ok")

	}
	mape := map[string]interface{}{
		"cookie": cookie,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(mape)

}
