package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	_ "math/rand"
	"net/http"
)

type SecureSethandler struct {
	DB *sql.DB
}

type Securest struct {
	Check string `json:"code"`
}

func (s *SecureSethandler) Secure_set(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		slog.Error("func secure_set: method don't allow")
		http.Error(w, "Method dont allow", http.StatusMethodNotAllowed)
		return

	}
	var f Securest
	ctx, cancel := iteranal.Contexte()
	defer cancel()
	var err error

	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {

		slog.Error("func secure set 1: error json parsing", err)
		http.Error(w, "Error parsing", http.StatusBadRequest)
		return

	}
	session, err := store.Get(r, "token1")
	if err != nil {
		slog.Error("func set secure 2 : ", err)
		return
	}

	code, ok := session.Values["code"]
	if !ok {
		slog.Info("Func secure set 3: ok", ok)
		slog.Info("code", code)
		return
	}

	val, ok := code.(string)
	if !ok {
		slog.Error("errrr", ok)
		return
	} else {
		slog.Info("Code", val)
	}

	mail, ok := session.Values["Email"]
	if !ok {
		slog.Error("func set_secure 5: don' t set")
		return
	}
	slog.Info("Email:", mail)

	slog.Info(f.Check)

	if val != f.Check {
		slog.Error("func secure set 6  dont code ", err)
		http.Error(w, "code don't", http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"status":  http.StatusOK,
		})

		return
	}

	_, err = s.DB.ExecContext(ctx, "update person SET two_factor_enabled = true WHERE email=$1", mail)

	if errors.Is(err, context.DeadlineExceeded) {
		slog.Error("func secure set 7 :  context timed out", err)
		http.Error(w, "context timed out", http.StatusConflict)
		return

	} else if err != nil {
		slog.Info("Func secure set 8 : err don't nill ", err)
		http.Error(w, "No nil", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"status":  200,
	})

}
