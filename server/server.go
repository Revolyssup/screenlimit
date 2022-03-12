package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Revolyssup/screenlimit/db"
)

type PassRequest struct {
	Pass string `json:"pass"`
}

func Run(port string, store *db.Store) {
	http.HandleFunc("/api/pass", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Fprint(w, "password is "+store.GetPassword())
			return
		}
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprint(w, "err "+err.Error())
				return
			}
			var p PassRequest
			err = json.Unmarshal(body, &p)
			if err != nil {
				fmt.Fprint(w, "err "+err.Error())
				return
			}
			store.SetPassword(p.Pass)
		}
	})
	http.ListenAndServe(":"+port, nil)
}
