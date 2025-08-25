package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
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

func parse(r *http.Request) (*HandlerRegister, error) {
	var err error
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()

	sugar := logger.Sugar()

	var jsonparse HandlerRegister
	err = json.NewDecoder(r.Body).Decode(&jsonparse)
	if err != nil {
		sugar.Error(
			"Errr", err)

	}
	defer r.Body.Close()
	return &jsonparse, err
}

func (d *LoginHandler) Login(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return

	}

	session, err := store.Get(r, "token1")
	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return
	}
	logger := zap.Must(zap.NewProduction())
	defer logger.Sync()
	sugar := logger.Sugar()

	// Импорты для работы
	ctx, cancel := iteranal.Contexte()
	defer cancel()

	t, err := parse(r)
	if err != nil {
		sugar.Error(
			"Errr", err)

	}

	var id int
	var password string
	var scrypt_salt string

	err = d.DB.QueryRowContext(ctx, `SELECT  id,password,scrypt_salt  FROM person WHERE email = $1`, t.Email).Scan(&id, &password, &scrypt_salt)
	slog.Info(password)
	ok := iteranal.CheckPassword(t.Password, password)

	if !ok {
		slog.Info("Func login dont")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	} else {
		slog.Info("func login check password: all okay ")
	}

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

	session.Values["cookie"] = cookie
	session.Values["salta"] = scrypt_salt

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   100000,
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
		slog.Info("func login8:no rows", err)
		return

	default:
		slog.Info("func register9:ok")

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(w)

}
