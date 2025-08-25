package cmd

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
)

type S3write struct {
	Title   string `json:"title"`
	Data    string `json:"data"`
	Content string `json:"content"`
}

func readJson(r *http.Request) (*S3write, error) {

	var e S3write
	var err error

	err = json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		slog.Info("Func S3 write error", err)
		return nil, err
	}
	defer r.Body.Close()

	return &e, nil
}

func S3(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Error("Method don't allow from S3")
		http.Error(w, "Method don't allow", http.StatusMethodNotAllowed)
		return
	}
	slog.Info("dsadasdasd")

	t, err := readJson(r)
	if err != nil {
		slog.Error("Error with Json ")
		return
	}
	ose_getDir, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Cant read ")
		return
	}

	info, err := os.Stat(ose_getDir + "/Downloads/qoutes")
	if err != nil || !info.IsDir() {
		err = os.Mkdir(ose_getDir+"/Downloads/qoutes", 0777)
		if err != nil {
			slog.Error("Erorr create ", err)
			return
		}
	}

	info2, err := os.Stat(ose_getDir + "/Downloads/qoutes/" + t.Title)
	if err != nil || !info2.IsDir() {
		file, err := os.Create(ose_getDir + "/Downloads/qoutes/" + t.Title)
		if err != nil {
			slog.Error("Http error with write", err)

			return
		}
		defer file.Close()

		_, err = file.WriteString(t.Content)
		if err != nil {
			slog.Error("Error write file ", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(w)
		return

	}

	slog.Info("File exist already ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(w)
}
