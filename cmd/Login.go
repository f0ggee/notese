package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
)

type LoginHandler struct {
	DB *sql.DB
}

type HandlerRegister struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

var store = sessions.NewCookieStore([]byte("KEY"))

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
	var id int

	err = d.DB.QueryRowContext(ctx, `SELECT  id FROM person WHERE email = $1 AND password = $2`, p.Email, p.Password).Scan(&id)
	switch {
	case err == context.DeadlineExceeded:
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		slog.Info("func login1:timed out")
		return

	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		slog.Info("func login2:no rows", err)
		return

	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func login3:no rows", err)
		return

	default:
		slog.Info("func login4:ok")

	}

	cookie := iteranal.Generateid()
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func login5:err", err)
		return
	}

	///expires := time.Now().Add(time.Hour * 24)
	session, err := store.Get(r, "token1")
	if err != nil {

		slog.Error("Cookie don't set")
		http.Error(w, "cokie don't set", http.StatusUnauthorized)
		return
	}

	session.Values["cookie"] = cookie
	slog.Info(cookie)

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   3000,
		Secure:   false,
		HttpOnly: true,
	}

	if err := session.Save(r, w); err != nil {
		slog.Error("Cokie can't send", err)
		return

	}

	_, err = d.DB.ExecContext(ctx, "UPDATE person SET cookie = $1  WHERE id = $2", cookie, id)

	switch {
	case err == context.DeadlineExceeded:
		http.Error(w, http.StatusText(http.StatusRequestTimeout), http.StatusRequestTimeout)
		slog.Info("func login6:timed out")
		return
	case err == sql.ErrNoRows:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		slog.Info("func login7:no rows")
		return

	case err != nil:
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		slog.Info("func login8:no rows")
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
