package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
)

var listenAddr string

func main() {
	flag.StringVar(&listenAddr, "addr", ":3000", "Listener spec")
	flag.Parse()
	fs := http.StripPrefix("/static", http.FileServer(http.Dir("static")))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			fs.ServeHTTP(w, r)
			return
		}

		d, err := ioutil.ReadFile("index.html")
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write(d)
	})

	log.Println("Listening on " + listenAddr + " ...")
	http.ListenAndServe(listenAddr, nil)
}
