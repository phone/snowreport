package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"flag"
)
import _ "net/http/pprof"

var (
	db   *sql.DB
	err  error
	port *string
	www  *bool
)

func collect() {
	var (
		locations []*Location
	)
	locations, err = GetLocations(db)
	for _, l := range locations {
		fcs, err := l.GetForecasts()
		if err != nil {
			log.Println(err)
			continue
		}
		for _, fc := range fcs {
			fc.Upsert(db)
		}
	}
}

func srv(port string) {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func init() {
	db, err = GetDb()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	www = flag.Bool("www", false, "run webserver")
	port = flag.String("p", "8080", "server listen port")
	flag.Parse()
	if !*www {
		collect()
	} else {
		srv(*port)
	}
}
