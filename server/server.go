package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Revolyssup/screenlimit/db"
	sysevents "github.com/Revolyssup/screenlimit/db/events"
	"github.com/Revolyssup/screenlimit/sysstats"
)

type PassRequest struct {
	Role string `json:"role"`
	Pass string `json:"pass"`
}

func auth(w *http.ResponseWriter, r *http.Request, pass string) bool {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	// fmt.Println("fucking entered")
	// cookie, err := r.Cookie("pass")
	// if err != nil {
	// user, rpass, ok := r.BasicAuth()
	// if !ok || user != db.ADMIN {
	// 	fmt.Println("not ok", user, rpass)
	// 	// http.Error(*w, "Please pass the username and password in the authorization header", http.StatusUnauthorized)
	// 	return false
	// }
	// fmt.Println("ok", user, rpass)

	// 	if pass != rpass {
	// 		// http.Error(*w, "Incorrect password", http.StatusUnauthorized)
	// 		return false
	// 	}
	// r.AddCookie(&http.Cookie{
	// 	Name:  "pass",
	// 	Value: pass,
	// })
	// c, _ := r.Cookie("pass")
	// fmt.Println("c", c)
	return true

	// fmt.Println("here")
	// if cookie.Value != pass {
	// 	return false
	// }
	// return true
}

func Run(port string, store *db.RoleStore, events *db.EventStore, pass string, apps *[]string, s *sysstats.StatCollector) {

	http.HandleFunc("/signin", func(w http.ResponseWriter, r *http.Request) {
		auth(&w, r, store.GetAdminPassword())
	})
	http.HandleFunc("/api/app", func(w http.ResponseWriter, r *http.Request) {
		if !auth(&w, r, store.GetAdminPassword()) {
			return
		}
		if r.Method == http.MethodGet {
			apps := s.GetApp()
			b, err := json.Marshal(apps)
			if err != nil {
				http.Error(w, "Bad request", http.StatusInternalServerError)
				return
			}
			w.Write(b)
			return
		}
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			var m = make(map[string]interface{})
			err = json.Unmarshal(body, &m)
			if err != nil {
				http.Error(w, "Bad request", http.StatusBadRequest)
				return
			}
			if app, ok := m["app"].(string); ok {
				s.AddApp(app)
				w.Write([]byte("App added: " + app))
				return
			}
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	})
	http.HandleFunc("/api/roles", func(w http.ResponseWriter, r *http.Request) {
		if !auth(&w, r, store.GetAdminPassword()) {
			return
		}
		if r.Method == http.MethodGet {
			r, err := store.GetRoles()
			if err != nil {
				fmt.Fprint(w, "err "+err.Error())
				return

			}
			roles := make(map[string]interface{}, 0)
			roles["data"] = r
			b, _ := json.Marshal(roles)
			// role := r.URL.Query().Get("role")
			// r, err := store.GetRole(role)
			// if err != nil {
			// 	fmt.Fprint(w, "err "+err.Error())
			// 	return

			// }
			// b, _ := json.Marshal(r)
			w.Write(b)
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
			store.SetRole(p.Role, p.Pass)
		}
	})

	http.HandleFunc("/api/events", func(w http.ResponseWriter, r *http.Request) {
		if !auth(&w, r, store.GetAdminPassword()) {
			return
		}
		fmt.Println(w.Header().Get("Access-Control-Allow-Origin"))
		if r.Method == http.MethodGet {
			ps := r.URL.Query().Get("page_size")
			off := r.URL.Query().Get("offset")
			t := r.URL.Query().Get("type")
			psi, offi := getpageoffset(ps, off)
			evs, err := events.Get(int(psi), int(offi), sysevents.Type(t))
			if err != nil {
				fmt.Println("Error while getting events: ", err.Error())
				fmt.Fprintf(w, err.Error())
				return
			}
			res := make(map[string]interface{}, 0)
			res["data"] = evs
			jsonevents, err := json.Marshal(res)
			if err != nil {
				fmt.Println("Error while parsing events: ", err.Error())
				fmt.Fprintf(w, err.Error())
				return
			}
			fmt.Println("successfully sent ", string(jsonevents))
			fmt.Fprintf(w, string(jsonevents))
			return
		}
		if r.Method == http.MethodPost {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				fmt.Fprint(w, "err "+err.Error())
				return
			}
			var pr = make(map[string]string)
			err = json.Unmarshal(body, &pr)
			if err != nil {
				fmt.Fprint(w, "err "+err.Error())
				return
			}
			s.AddApp(pr["app"])
			s.AddActionApp(pr["app"], pr["action"])
			fmt.Fprintf(w, string("Now will "+pr["action"]+" on starting "+pr["app"]))
			return
		}
	})
	// http.HandleFunc("/home", func(w http.ResponseWriter, r *http.Request) {
	// 	if !auth(&w, r, pass) {
	// 		fmt.Println("here brpooooo")
	// 		http.Redirect(w, r, "/login", http.StatusMovedPermanently)
	// 		return
	// 	}
	// 	http.ServeFile(w, r, "./server/build")
	// })
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
