package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type IpAddress struct {
	ip string `json:"ip"`
}

func main() {
	var dir string
	resp, _ := http.Get("https://api.ipify.org?format=json")
	body, _ := ioutil.ReadAll(resp.Body)
	//var address IpAddress
	fmt.Println("ip now " + string(body))

	prt := os.Getenv("PRT")
	if _, err := strconv.Atoi(prt); err == nil {
		fmt.Printf("%q is the Port Set in Env Variables ", prt)
	} else {
		fmt.Printf("setting default port 8080\n")
		prt = "8080"
	}
	fmt.Println("running port ", prt, "\n")

	flag.StringVar(&dir, "dir", ".", "the directory to serve files from. Defaults to the current dir")
	flag.Parse()
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir(dir))))

	srv := &http.Server{
		Handler: r,
		Addr:    ":" + prt,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
