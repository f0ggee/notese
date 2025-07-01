package cmd

import (
	"Project2/iteranal"
	"context"
	_ "crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gorilla/sessions"
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

func (e *RegisterHandler) Register(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}
	var err error
	var person *Person

	///контекст

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	slog.Info("fffff")

	if err = json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Info("Func register1:", err)
		return
	}

	if !strings.Contains(person.Email, "@") {
		http.Error(w, "Person name must contain @", http.StatusBadRequest)
		slog.Info("Func register2:", err)
		return
	}

	var existingPerson bool

	err = e.DB.QueryRowContext(ctx, "SELECT EXISTS (select 1 FROM person WHERE email=$1)", person.Email).Scan(&existingPerson)
	///проверка на валидность запроса

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

	var userid int64

	err = e.DB.QueryRowContext(ctx, "INSERT INTO person(name,email,password) VALUES ($1,$2,$3) RETURNING id", person.Name, person.Email, person.Password).Scan(&userid)

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
	slog.Info("Func register11112212")

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
		MaxAge:   3000,
		Secure:   false,
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

	mape := map[string]interface{}{
		"ID":   id,
		"name": person.Name,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mape); err != nil {
		slog.Error("Func register14:", err)
		return
	}

}
