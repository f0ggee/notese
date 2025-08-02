package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func chekcpasswor(hash string, password string) bool {

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err == nil {
		return true
	}
	return false

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
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return

	}

	var err error
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

	err = d.DB.QueryRowContext(ctx, `SELECT  id,password  FROM person WHERE email = $1`, t.Email).Scan(&id, &password)
	slog.Info(password)
	ok := chekcpasswor(password, t.Password)

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
	mape := map[string]interface{}{
		"cookie": cookie,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(mape)

}
