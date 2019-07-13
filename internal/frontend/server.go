package frontend

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/illfalcon/parser/internal/crawler"
	"github.com/illfalcon/parser/internal/parser"

	"github.com/illfalcon/parser/internal/db"

	"log"
	"net/http"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

type Landing struct {
	RowID int
	URL   string
	//Hash string
	Name string
}

type Event struct {
	RowID       int
	URL         string
	Hash        string
	Article     string
	Probability string
	Start       string
	End         string
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

var tmpl = template.Must(template.ParseGlob("internal/frontend/view/*.tmpl"))

func Index(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Not Found")
	} else {
		http.Redirect(w, r, "/events", 301)
	}
}

func Calendar(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)
	if r.URL.Path != "/" {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Not Found")
	} else {
		http.Redirect(w, r, "/events", 301)
	}
}

func Events(w http.ResponseWriter, r *http.Request) {
	conn := db.Conn()
	approved := r.URL.Query().Get("approved")
	if approved == "yes" {
		approved = "WHERE approved = 1"
	} else if approved == "no" {
		approved = "WHERE approved = 0"
	} else if approved == "null" {
		approved = "WHERE approved is null"
	} else if approved == "all" {
		approved = ""
	} else {
		http.Redirect(w, r, "/events?approved=null", 301)
	}
	stmt, err := conn.Query("SELECT rowid, url, hash, article, probability, event_start, event_end FROM events " + approved + " ORDER BY probability DESC;")
	if err != nil {
		panic(err.Error())
	}
	event := Event{}
	res := []Event{}
	for stmt.Next() {
		var rowid int
		var probability float64
		var url, hash, article string
		var eventStartNullable, eventEndNullable sql.NullString
		err = stmt.Scan(&rowid, &url, &hash, &article, &probability, &eventStartNullable, &eventEndNullable)
		if err != nil {
			panic(err.Error())
		}
		if !eventStartNullable.Valid {
			event.Start = time.Now().AddDate(0, 0, 1).Format("2006-01-02") + " 12:00"
		} else {
			event.Start = eventStartNullable.String
		}
		if !eventEndNullable.Valid {
			event.End = time.Now().AddDate(0, 0, 1).Format("2006-01-02") + " 16:00"
		} else {
			event.End = eventEndNullable.String
		}

		event.RowID = rowid
		event.Probability = fmt.Sprintf("%.0f", probability*100) + "%"
		event.URL = url
		event.Hash = hash
		event.Article = article
		res = append(res, event)
	}
	tmpl.ExecuteTemplate(w, "Events", res)
	defer conn.Close()
}

func EventsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		fmt.Println()
		col := r.FormValue("name")
		value := r.FormValue("value")
		id := r.FormValue("pk")
		conn := db.Conn()
		if col == "article" {
			stmt, err := conn.Prepare(`UPDATE events SET article = ? WHERE rowid = ?`)
			checkErr(err)
			_, err = stmt.Exec(value, id)
			checkErr(err)
		} else if col == "approval" {
			stmt, err := conn.Prepare(`UPDATE events SET approved = ? WHERE rowid = ?`)
			checkErr(err)
			_, err = stmt.Exec(value, id)
			checkErr(err)
		} else if col == "event_start" || col == "event_end" {
			parsed, err := time.Parse("2006-01-02 15:04", value)
			if err == nil {
				fmt.Println(parsed)
				stmt, err := conn.Prepare(`UPDATE events SET ` + col + ` = ? WHERE rowid = ?`)
				checkErr(err)
				_, err = stmt.Exec(value, id)
				checkErr(err)
			}
		} else if col == "url" {
			stmt, err := conn.Prepare(`UPDATE events SET url = ? WHERE rowid = ?`)
			checkErr(err)
			_, err = stmt.Exec(value, id)
			checkErr(err)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
		defer conn.Close()
		http.Redirect(w, r, r.Referer(), 301)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func Landings(w http.ResponseWriter, r *http.Request) {
	conn := db.Conn()
	selDB, err := conn.Query("SELECT rowid, * FROM landings ORDER BY rowid ASC")
	if err != nil {
		panic(err.Error())
	}
	landing := Landing{}
	res := []Landing{}
	for selDB.Next() {
		var rowid int
		var url, hash, name sql.NullString
		err = selDB.Scan(&rowid, &url, &hash, &name)
		if err != nil {
			panic(err.Error())
		}
		landing.RowID = rowid
		landing.URL = url.String
		//landing.Hash = hash
		landing.Name = name.String
		res = append(res, landing)
	}
	tmpl.ExecuteTemplate(w, "Landings", res)
	defer conn.Close()
}

func LandingsUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		col := r.FormValue("name")
		value := r.FormValue("value")
		id := r.FormValue("pk")
		if col == "name" || col == "url" {
			conn := db.Conn()
			stmt, err := conn.Prepare(`UPDATE landings SET ` + col + ` = ? WHERE rowid = ?`)
			checkErr(err)
			_, err = stmt.Exec(value, id)
			checkErr(err)
			defer conn.Close()
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func LandingsAdd(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		name := r.FormValue("name")
		url := r.FormValue("url")
		conn := db.Conn()
		stmt, err := conn.Prepare(`INSERT INTO landings(name,url) VALUES(?, ?);`)
		checkErr(err)
		_, err = stmt.Exec(name, url)
		checkErr(err)
		defer conn.Close()
		http.Redirect(w, r, "/landings", 301)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func LandingsDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		id := r.FormValue("id")
		conn := db.Conn()
		stmt, err := conn.Prepare(`DELETE FROM landings WHERE rowid=?`)
		checkErr(err)
		_, err = stmt.Exec(id)
		checkErr(err)
		defer conn.Close()
		http.Redirect(w, r, "/landings", 301)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	fmt.Print(r.URL.Path[1:])
	http.ServeFile(w, r, "internal/frontend/"+r.URL.Path[1:])
}

func Renew(w http.ResponseWriter, r *http.Request) {
	go func() {
		log.Println("Started")
		crawler.Crawl()
		parser.Parse()
		log.Println("Finished")
	}()
}

func Start() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/events", Events)
	http.HandleFunc("/events/update", EventsUpdate)
	http.HandleFunc("/landings", Landings)
	http.HandleFunc("/landings/update", LandingsUpdate)
	http.HandleFunc("/landings/add", LandingsAdd)
	http.HandleFunc("/landings/delete", LandingsDelete)
	http.HandleFunc("/calendar", Calendar)
	//http.HandleFunc("/favicon.ico", Favicon)
	http.HandleFunc("/public/", ServeFile)
	http.HandleFunc("/renew", Renew)
	//http.HandleFunc("/show", Show)
	//http.HandleFunc("/new", New)
	//http.HandleFunc("/edit", Edit)
	//http.HandleFunc("/insert", Insert)
	//http.HandleFunc("/update", Update)
	//http.HandleFunc("/delete", Delete)
	log.Println("Server started on: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
