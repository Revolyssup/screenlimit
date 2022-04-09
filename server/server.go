package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Revolyssup/screenlimit/db"
	"github.com/Revolyssup/screenlimit/db/actions"
)

type PassRequest struct {
	Role string `json:"role"`
	Pass string `json:"pass"`
}

func Run(port string, store *db.Roles, events *db.Events) {
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
			store.SetPassword(p.Role, p.Pass)
		}
	})

	http.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			ps := r.URL.Query().Get("page_size")
			off := r.URL.Query().Get("offset")
			t := r.URL.Query().Get("type")
			psi, offi := getpageoffset(ps, off)
			evs, err := events.Get(int(psi), int(offi), actions.Type(t))
			if err != nil {
				fmt.Println("Error while getting events: ", err.Error())
				fmt.Fprintf(w, err.Error())
				return
			}
			jsonevents, err := json.Marshal(evs)
			if err != nil {
				fmt.Println("Error while parsing events: ", err.Error())
				fmt.Fprintf(w, err.Error())
				return
			}
			fmt.Println("successfully sent ", string(jsonevents))
			fmt.Fprintf(w, string(jsonevents))
			return
		}
		// if r.Method == http.MethodPost {
		// 	event := db.Events{}
		// 	body, err := ioutil.ReadAll(r.Body)
		// 	if err != nil {
		// 		fmt.Fprint(w, "err "+err.Error())
		// 		return
		// 	}
		// 	err = json.Unmarshal(body, &event)
		// 	if err != nil {
		// 		fmt.Fprint(w, "err "+err.Error())
		// 		return
		// 	}
		// 	err = events.Add(event.Timestamp, event.Action, events.ActionType)
		// 	if err != nil {
		// 		fmt.Fprint(w, "err "+err.Error())
		// 		return
		// 	}
		// }
	})
	http.ListenAndServe(":"+port, nil)
}

func getpageoffset(ps string, off string) (psi int64, offi int64) {
	if ps == "" {
		psi = 10
	}
	if off == "" {
		offi = 0
	}
	psi, err := strconv.ParseInt(ps, 10, 64)
	if err != nil {
		psi = 10
	}
	offi, err = strconv.ParseInt(off, 10, 64)
	if err != nil {
		offi = 0
	}
	return
}
