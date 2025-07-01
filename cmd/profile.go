package cmd

import (
	"Project2/iteranal"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type ProfileHandler struct {
	DB *sql.DB
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
		return
	}

	uuidid, ok := session.Values["cookie"]
	if !ok {
		slog.Error("dont get")
		return

	}

	var name string
	var id int64
	err1 := e.DB.QueryRowContext(ctx, "SELECT name,id FROM person WHERE cookie = $1 LIMIT 1 ", uuidid).Scan(&name, &id)

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
		log.Println("Func profile2:failed to connetct", err1)
		http.Error(w, "Cookie not found", http.StatusUnauthorized)
		return
	}

	//response := map[string]interface{}{
	//	"id":   userID,
	//	"name": NameFromBD,
	//}

	response := map[string]interface{}{
		"id":   id,
		"name": name,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}
