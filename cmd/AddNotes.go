package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"net/http"
	"time"
)

type AddNotesHandler struct {
	DB *pgxpool.Pool
}

type AddNotesResponse struct {
	Title   string    `json:"title"`
	Data    time.Time `json:"data"`
	Content string    `json:"content"`
}

func (e *AddNotesHandler) NewAddNotesCmd(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", http.StatusMethodNotAllowed)
		slog.Info("func addnotes1:")
		return
	}

	var err error
	var b AddNotesResponse

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "", http.StatusBadRequest)
		slog.Info("func addnotes2:", err)
		return
	}

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	cookie, err := r.Cookie("token")
	if err != nil {
		slog.Info("func addnotes3:", err)
		return
	}
	var id int

	err = e.DB.QueryRow(ctx, "INSERT INTO notes(title,created_at,content,author_cookie) VALUES ($1,$2,$3,$4) RETURNING id", b.Title, b.Data, b.Content, cookie.Value).Scan(&id)

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
