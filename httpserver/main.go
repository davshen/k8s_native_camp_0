package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	//"github.com/golang/glog"
)

func main() {
	flag.Set("v", "4")
	//glog.V(2).Info("Starting http server...")
	log.Println("start")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/sayhelloName", sayhelloName)
	http.HandleFunc("/addHeader", addHeader)
	fmt.Println("Runing...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		log.Fatal(err)
	}

}

func addHeader(w http.ResponseWriter, r *http.Request) {
	req_header := (map[string][]string)(r.Header)
	for k, arr := range req_header {
		w.Header().Add(k, arr[0])
	}
	fmt.Fprintf(w, "req_Header全部数据：", req_header)
	w.Header().Add("key", "dawei")
	fmt.Fprintf(w, "\nreponse_Header全部数据：", w.Header())

	var VERSION string
	VERSION = os.Getenv("VERSION")
	fmt.Println("VERSION", VERSION)

}

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form)
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello World!")
}

func healthz(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "ok!\n   ")
	w.WriteHeader(200)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	user := r.URL.Query().Get("user")
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
}
