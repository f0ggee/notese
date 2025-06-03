package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"
)

type AddNotesCmd struct {
	db *sql.DB
}

type AddNotesResponse struct {
	title   string `json:"title"`
	data    string `json:"data"`
	content string `json:"content"`
}

func NewAddNotesCmd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		slog.Info("func addnotes1:")
		return
	}
	var a *AddNotesCmd
	var err error
	var b *AddNotesResponse

	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		slog.Info("func addnotes2:", err)
		return
	}

	ctx := iteranal.Contexte()

	cookie, err := r.Cookie("token")
	if err != nil {
		slog.Info("func addnotes3:", err)
		return
	}
	var id int

	err = a.db.QueryRowContext(ctx, "INSERT INTO notes(title,data,content,author_cookie) VALUES ($1,$2,$3) RETURNING id", b.title, b.data, b.content, cookie.Value).Scan(&id)

	if err == context.DeadlineExceeded || err == context.Canceled {
		http.Error(w, "", http.StatusRequestTimeout)
		slog.Info("func addnotes3:", err)
		return
	}

	if err == sql.ErrNoRows {
		http.Error(w, "", http.StatusNotFound)
		slog.Info("func addnotes4:", err)
		return
	}

	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		slog.Info("func addnotes5:", err)
		return
	}

	slog.Info("func addnotes5: all okay :)", id)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("You notew write"))
	json.NewEncoder(w).Encode(w)

}
