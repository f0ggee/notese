package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Maein(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Mehod don't allow", http.StatusMethodNotAllowed)
		return
	}

	session, err := store.Get(r, "token1")
	if err != nil {
		slog.Info("Cookie doen't set")
		http.Error(w, "Cookie don't has", http.StatusUnauthorized)
		return
	}

	theme, ok := session.Values["th"].(string)
	if !ok {
		slog.Info("Theme don't set ")
		return
	}

	slog.Info("Themee2", theme)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"theme": theme,
	})
}
