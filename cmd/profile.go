package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type ProfileHandler struct {
	DB *sql.DB
}

type Notes struct {
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"createdAt"`
	Content   string    `json:"content"`
}

func (e *ProfileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		log.Printf("func profile: Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	session, err := store.Get(r, "token1")
	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "Cookie dont sent", http.StatusUnauthorized)
		return
	}
	themem, ok := session.Values["color"]
	if !ok {
		slog.Info("Don't set ")
	}

	uuidid, ok := session.Values["cookie"]
	if !ok {
		slog.Error("dont get")
		return

	}

	var name string
	var id int64

	err1 := e.DB.QueryRowContext(ctx, "SELECT name,id  FROM person WHERE cookie = $1 LIMIT 1 ", uuidid).Scan(&name, &id)

	if err1 == context.DeadlineExceeded {
		slog.Info("Func profile1:deadline exceeded")
		http.Error(w, "Deadline exceeded", http.StatusRequestTimeout)
		return
	}

	if err1 == sql.ErrNoRows {
		slog.Info("Func profile1:failed to find person")
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}
	if err1 != nil {
		log.Println("Func profile2:failed to connect", err1)
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}

	rows, err2 := e.DB.QueryContext(ctx, "SELECT title,created_at,content  FROM notes  WHERE author_cookie = $1  ", uuidid)
	if err2 == context.DeadlineExceeded {
		slog.Info("Func profile1:deadline exceeded")
		http.Error(w, "Deadline exceeded", http.StatusRequestTimeout)
		return
	}

	if err2 == sql.ErrNoRows {
		slog.Info("Func profile1:failed to find person")
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}
	if err2 != nil {
		log.Println("Func profile2:failed to connetct", err1)
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}

	defer rows.Close()

	var note []Notes

	for rows.Next() {
		var n Notes

		err43 := rows.Scan(&n.Title, &n.CreatedAt, &n.Content)

		if err43 != nil {
			slog.Error("Err func profile 4", err43)
			return
		}

		note = append(note, n)

	}

	//response := map[string]interface{}{
	//	"id":   userID,
	//	"name": NameFromBD,
	//}

	response := map[string]interface{}{
		"name":  name,
		"theme": themem,

		"notes": note,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
