package cmd

import (
	"Project2/iteranal"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
)

type ProfileHandler struct {
	DB *sql.DB
}

func Transfer_Passworde(s string, nonce string, data string) (string, error) {

	scrypt, err := hex.DecodeString(s)
	if err != nil {
		slog.Error("Err in add notes", err)
		return "", err
	}
	Nonce, err := hex.DecodeString(nonce)
	if err != nil {
		slog.Error("Error in ", err)
	}

	datae, err := hex.DecodeString(data)
	if err != nil {
		slog.Error("Cant convert", err)
		return "", err
	}
	block, err := aes.NewCipher(scrypt)
	if err != nil {
		slog.Error("Err", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("Errr")
		return "", err
	}

	cyphertext, err := aesGCM.Open(nil, Nonce, datae, nil)

	return string(cyphertext), nil
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

	uuidid, ok := session.Values["cookie"]
	if !ok {
		slog.Error("dont get")
		return

	}
	themem, oke := session.Values["th"].(string)
	if !oke {
		slog.Info("Don't set,theme ", themem)

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

	rows, err2 := e.DB.QueryContext(ctx, "SELECT hash_notes,nonce  FROM notes  WHERE author_cookie = $1  ", uuidid)
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
	var hash_notes string
	var nonce string

	salate, oke := session.Values["salta"].(string)
	if !oke {
		slog.Info("Func addNotes can't get", oke)
	}

	defer rows.Close()

	var note []string

	for rows.Next() {

		err43 := rows.Scan(&hash_notes, &nonce)

		if err43 != nil {
			slog.Error("Err func profile 4", err43)
			return
		}

		dyscrypt, err23 := Transfer_Passworde(salate, nonce, hash_notes)
		if err23 != nil {
			slog.Error("ERR", err)
			return
		}

		note = append(note, dyscrypt)

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
