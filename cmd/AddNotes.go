package cmd

import (
	"Project2/iteranal"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
)

type AddNotesHandler struct {
	DB *sql.DB
}

type AddNotesResponse struct {
	Title   string `json:"title"`
	Data    string `json:"data"`
	Content string `json:"content"`
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

func Transfer_Password(slate, structe string) (string, string, error) {

	Dear, err := hex.DecodeString(slate)
	if err != nil {
		slog.Error("Err in add notes", err)
		return "", "", err
	}
	block, err := aes.NewCipher(Dear)
	if err != nil {
		slog.Error("Err", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("Errr")
		return "", "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err)
	}

	cyphertext := aesGCM.Seal(nil, nonce, []byte(structe), nil)

	cyphertexte := hex.EncodeToString(cyphertext)
	noncee := hex.EncodeToString(nonce)

	slog.Info("Шифровка", cyphertexte)

	return cyphertexte, noncee, nil
}

func (e *AddNotesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	//parsedtime, err := time.Parse("2006-01-02T15:04", b.Data)
	if err != nil {
		slog.Info("func addnotes2:", err)
		return
	}

	ctx, cancel := iteranal.Contexte()
	defer cancel()

	session, err := store.Get(r, "token1")
	if err != nil {
		slog.Error("cookie don't send", err)
		http.Error(w, "cookie dont sen", http.StatusUnauthorized)
		return
	}
	salate, oke := session.Values["salta"].(string)
	if !oke {
		slog.Info("Func addNotes can't get", oke)
	}
	str := fmt.Sprintf("%+v", b)

	Hash, nonce, err := Transfer_Password(salate, str)
	if err != nil {
		slog.Info("Can't parse", err)
		return
	}

	uuidid, ok := session.Values["cookie"]
	if !ok {
		slog.Error("dont get")
		return

	}

	_, _ = e.DB.ExecContext(ctx, "INSERT INTO notes(hash_notes,nonce,author_cookie) values ($1,$2,$3)", Hash, nonce, uuidid)

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

	slog.Info("func addnotes6: ALL OKAY :0 ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(w)

}
