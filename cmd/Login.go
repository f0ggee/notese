package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
)

type LoginHandler struct {
	DB *pgxpool.Pool
}

type HandlerRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (d *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return

	}

	var err error
	var p HandlerRegister

	// Импорты для работы
	ctx, cancel := iteranal.Contexte()
	defer cancel()

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	var exists bool

	err = d.DB.QueryRow(ctx, "SELECT exists(SELECT 2 FROM person WHERE email=$1 and password=$2)", p.Email, p.Password).Scan(&exists)
	switch {
	case err == context.DeadlineExceeded:
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		slog.Info("func login1:timed out")
		return

	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func login3:no rows", err)
		return

	case !exists:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		slog.Info("func login4:no rows", err)
		return

	default:
		slog.Info("func login4:ok")

	}

	cookie := iteranal.Generateid()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func register5:err", err)
		return
	}

	///expires := time.Now().Add(time.Hour * 24)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    cookie,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   true,
	})

	_, err = d.DB.Exec(ctx, "UPDATE person SET session_id = $1  WHERE email = $2", cookie, p.Email)

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
