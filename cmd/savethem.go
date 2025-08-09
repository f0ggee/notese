package main

import (
	"github.com/gorilla/sessions"
	"log/slog"
	"net/http"
)

type Themestruch struct {
	Themee string `json:"json:theme"`
}

var store = sessions.NewCookieStore([]byte("KEY"))

	defer r.Body.Close()

	return &e, err
}

func Theme(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Mehod dont allow", http.StatusBadRequest)
		slog.Info("H")
		return

		}

		session, err := store.Get(r, "token1")
		if err != nil {
			slog.Error("cookie don't send", err)
			http.Error(w, "cookie dont sen", http.StatusUnauthorized)
			return
		}
		slog.Info("Them", ea.Themee)

		session.Values["color"] = ea.Themee
		session.Save(r, w)

	}

}
