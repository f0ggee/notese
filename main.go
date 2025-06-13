package main

import (
	"Project2/cmd"
	"github.com/joho/godotenv"
	"log/slog"

	_ "github.com/jackc/pgx/v5/pgxpool"
	"log"
	_ "log/slog"
	"net/http"
	_ "time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "fronted/index.html")
			return
		}
		http.ServeFile(w, r, "./fronted"+r.URL.Path)
	})
	http.HandleFunc("/add_notes", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/addnotes.html")
	})
	http.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/profile.html")
	})
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")

	}

	db, err := cmd.Connect()
	if db == nil {
		slog.Info("database connection failed")
		return
	}

	loginHandler := &cmd.LoginHandler{DB: db}
	profileHandler := &cmd.ProfileHandler{DB: db}
	registerHandler := &cmd.RegisterHandler{DB: db}
	addNotesHandler := &cmd.AddNotesHandler{DB: db}
	slog.Info("f")

	http.HandleFunc("/addnotes/api", addNotesHandler.NewAddNotesCmd)

	http.HandleFunc("/register/api", registerHandler.Register)
	http.HandleFunc("/profile/api", profileHandler.Profile)
	http.HandleFunc("/login/api", loginHandler.Login)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
