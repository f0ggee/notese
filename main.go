package main

import (
	"Project2/cmd"
	"Project2/iteranal"
	_ "embed"
	"github.com/joho/godotenv"
	"log/slog"

	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"log"
	_ "log/slog"
	"net/http"
	_ "time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "fronted/index.html")
			return
		}
		http.ServeFile(w, r, "./fronted"+r.URL.Path)
	})
	router.HandleFunc("/add_notes", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/addnotes.html")
	})
	router.HandleFunc("/profile", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/profile.html")
	})
	router.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/register.html")

	})

	router.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/main.html")
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

	loginhandler := &cmd.LoginHandler{DB: db}
	registerhandler := &cmd.RegisterHandler{DB: db}
	profilehandler := &cmd.ProfileHandler{DB: db}
	addnoteshandler := &cmd.AddNotesHandler{DB: db}

	loging := iteranal.Mildwary

	router.HandleFunc("/addnotes/api", addnoteshandler.ServeHTTP).Methods("POST")

	router.HandleFunc("/register/api", registerhandler.Register).Methods("POST")
	handler := loging(router)

	router.HandleFunc("/profile/api", profilehandler.ServeHTTP).Methods("GET")
	router.HandleFunc("/login/api", loginhandler.Login).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", handler))

}
