package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/golang/glog"
	"github.com/k8s_native_camp_0/httpserver/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	_ "net/http/pprof"
)

func main() {
	flag.Set("v", "4")
	glog.V(2).Info("Starting http server...")
	//log.Println("start")
	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	//http.HandleFunc("/", rootHandler)
	mux.HandleFunc("/healthz", healthz)
	mux.HandleFunc("/sayhelloName", sayhelloName)
	mux.HandleFunc("/addHeader", addHeader)
	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    ":80",
		Handler: mux,
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: %s\n", err)
		}
	}()
	log.Print("Server Started")
	<-done
	log.Print("Server Stopped")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// extra handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Printf("Server Exited Properly")
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

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("entering root handler")
	glog.V(4).Info("entering root handler")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(10, 2000)
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	glog.V(4).Infof("Respond in %d ms", delay)
}
