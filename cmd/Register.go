package cmd

import (
	"Project2/iteranal"
	"context"
	"crypto/rand"
	_ "crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/scrypt"
	"log/slog"
	"net/http"
	"strings"
)

type RegisterHandler struct {
	DB *sql.DB
}

type Person struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func chehkjson(r *http.Request) (*Person, error) {
	var err error
	var e Person

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return &e, err
}

func (e *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Only POST method is supported.", http.StatusMethodNotAllowed)
		return

	}
	var err error

	///контекст

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	t, err := chehkjson(r)
	if err != nil {
		slog.Error("Func login beta", err)
		return
	} else {
		slog.Info("all okay func register ")
	}

	if !strings.Contains(t.Email, "@") {
		http.Error(w, "Person name must contain @", http.StatusBadRequest)
		slog.Info("Func register2:", err)
		return
	}

	var existingPerson bool

	err = e.DB.QueryRowContext(ctx, "SELECT EXISTS (select 1 FROM person WHERE email=$1)", t.Email).Scan(&existingPerson)
	///проверка на валидность запрос
	if err == context.DeadlineExceeded {

		slog.Error("Func register2:", err)
		http.Error(w, "Timeout exceeded.", http.StatusRequestTimeout)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusNotFound)
		slog.Info("Func register3:", err)
		return
	} else if err != nil {
		slog.Info("Func register4:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if existingPerson {
		http.Error(w, "person already exists", http.StatusConflict)
		slog.Info("Func register5:", err)
		return
	}

	////

	f, err := iteranal.HashPassowrd(t.Password)
	if err != nil {
		slog.Error("Cant generate password", err)
		return
	}
	salt := make([]byte, 16)
	rand.Read(salt)
	D3, _ := scrypt.Key([]byte(f), salt, 1<<15, 8, 1, 32)
	D4 := hex.EncodeToString(D3)

	var userid int
	err = e.DB.QueryRowContext(ctx, "INSERT INTO person(name,email,password,scrypt_salt) VALUES ($1,$2,$3,$4) RETURNING id", t.Name, t.Email, f, D4).Scan(&userid)

	if errors.Is(err, context.DeadlineExceeded) {
		http.Error(w, err.Error(), http.StatusRequestTimeout)
		slog.Info("Func register6:", err)
		return
	} else if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, err.Error(), http.StatusNotFound)
		slog.Info("Func register7:", err)
		return
	} else if err != nil {
		slog.Info("Func register8:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := iteranal.Generateid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info("Func register9:", err)
		return
	}

	//expires := time.Now().Add(time.Hour * 24)
	session, err := store.Get(r, "token1")
	if err != nil {

		slog.Error("Cookie don't set")
		http.Error(w, "cokie don't set", http.StatusUnauthorized)
		return
	}

	session.Values["cookie"] = id

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   1000000000,
		Secure:   true,
		HttpOnly: true,
	}

	if err := session.Save(r, w); err != nil {
		slog.Error("Cokie can't send", err)
		return

	}

	_, err = e.DB.ExecContext(ctx, "UPDATE person Set cookie = $1 WHERE id = $2", id, userid)
	switch {
	case err == context.DeadlineExceeded:
		slog.Info("Func register10:", err)
		http.Error(w, "Timeout trying to register a person", http.StatusRequestTimeout)
		return

	case err == sql.ErrNoRows:
		slog.Info("Func register11:", err)
		http.Error(w, "Person already exists", http.StatusConflict)
		return

	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info("Func register12:", err)
		return

	default:
		slog.Info("Func register13: all okay :)")

	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(w); err != nil {
		slog.Error("Func register14:", err)
		return
	}

}
