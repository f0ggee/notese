package cmd

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
)

type Themestruch struct {
	Themee string `json:"theme"` // оставляем под фронт

}

func chechkjson(r *http.Request) (*Themestruch, error) {
	var e Themestruch
	var err error

	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return nil, err
	}

	defer r.Body.Close()

	return &e, err
}

func Theme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Mehod dont allow", http.StatusBadRequest)
		slog.Info("H")
		return

	}

	theme, err := chechkjson(r)
	if err != nil {
		slog.Error("Func savthem", err)
		return
	}

	session, err := store.Get(r, "token1")
	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return
	}
	slog.Info("Theme from chosethem", theme.Themee)

	session.Values["th"] = theme.Themee

	session.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   100000,
		Secure:   true,
		HttpOnly: true,
	}

	if err := session.Save(r, w); err != nil {
		slog.Error("Cokie can't send", err)
		return

	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"theme":  theme.Themee,
	})
}
