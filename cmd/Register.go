package cmd

import (
	"Project2/iteranal"
	"context"
	_ "crypto/rand"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
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

	var err error
	var person *Person

	///контекст

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	if r.Method != "POST" {
		http.Error(w, "Only GET method is supported.", http.StatusMethodNotAllowed)
		return

	}

	if err = json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Info("Func register1:", err)
		return
	}

	var exits bool

	err = e.DB.QueryRowContext(ctx, "SELECT EXISTS(select 1 FROM person WHERE email = $1)", person.Email).Scan(&exits)
	///проверка на валидность запроса

	if err == context.DeadlineExceeded {
		slog.Info("Func register2:", err)
		http.Error(w, "Timed out", http.StatusRequestTimeout)
		return
	} else if err == sql.ErrNoRows {
		slog.Info("Func register3:", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	} else if err != nil {
		slog.Info("Func register4:", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	////

	if exits {
		slog.Info("Func register5:person already exists")
		http.Error(w, "Person already exists", http.StatusConflict)
		return

	}

	var userid int64

	err = e.DB.QueryRowContext(ctx, "INSERT INTO person(name,email,password) VALUES ($1,$2,$3) RETURNING id", person.Name, person.Email, person.Password).Scan(&userid)

	switch {
	case err == context.DeadlineExceeded:
		http.Error(w, "Timeout trying to register a person", http.StatusRequestTimeout)
		slog.Info("Func register6:", err)
		return

	case err == sql.ErrNoRows:
		slog.Info("Func register7:", err)
		http.Error(w, "Person already exists", http.StatusUnauthorized)
		return

	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info("Func register8:", err)
		return

	}
	id := iteranal.Generateid()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info("Func register9:", err)
		return
	}

	//expires := time.Now().Add(time.Hour * 24)

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "id",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: false,
		Secure:   false,
	})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Info("Func register10:", err)
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
	json.NewEncoder(w).Encode(mape)

}
