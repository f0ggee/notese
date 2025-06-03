package main

import (
	"Project2/cmd"
	"context"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"time"
)

//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./fronted/index.html")
			return
		}
		http.ServeFile(w, r, "./fronted"+r.URL.Path)
	})

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")

	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	db, _ := cmd.Connect(ctx)
	defer db.Close()

	e := &cmd.Logincmd{DB: db}
	a := &cmd.Handler_login{DB: db}

	http.HandleFunc("/register/api", e.Register)
	http.HandleFunc("/login/api", a.Login)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
