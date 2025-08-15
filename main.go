package main

import (
	"Project2/cmd"
	"Project2/iteranal"
	_ "embed"
	"github.com/gorilla/mux"
	_ "github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"log"
	"log/slog"
	_ "log/slog"
	"net/http"
	_ "time"
)

// TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
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

	router.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/twofa-success.html")

	})

	router.HandleFunc("/study", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/Study.html")
	})
	router.HandleFunc("/two_check", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/two.check.html")
	})
	router.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/set.html")
	})

	router.HandleFunc("/main", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/main.html")
	})

	router.HandleFunc("/setting", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/setting.html")
	})
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла", err)

	}

	router.HandleFunc("/chose", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "fronted/chosetheme.html")
	})

	db, err := cmd.Connect()
	if db == nil {
		slog.Info("database connection failed")
		return
	}

	loginhandler := &cmd.LoginHandler{DB: db}
	registerhandler := &cmd.RegisterHandler{DB: db}
	profilehandler := &cmd.ProfileHandler{DB: db}
	addnoteshandler := &cmd.AddNotesHandler{DB: db}
	securehandler := &cmd.SecureHandler{DB: db}
	securesethandler := &cmd.SecureSethandler{DB: db}

	loging := iteranal.Mildwary
	handler := loging(router)

	router.HandleFunc("/addnotes/api", addnoteshandler.ServeHTTP).Methods("POST")

	router.HandleFunc("/setting/api", cmd.Settinge).Methods("GET")

	router.HandleFunc("/register/api", registerhandler.Register).Methods("POST")
	router.HandleFunc("/securest/api", securesethandler.Secure_set).Methods("POST")
	router.HandleFunc("/main/api", cmd.Maein).Methods("GET")
	router.HandleFunc("/chose/api", cmd.Theme).Methods("POST")

	router.HandleFunc("/secure/api", securehandler.Secure).Methods("POST")
	router.HandleFunc("/profile/api", profilehandler.ServeHTTP).Methods("GET")
	router.HandleFunc("/login/api", loginhandler.Login).Methods("POST")

	log.Fatal(http.ListenAndServe("localhost:8080", handler))
}
