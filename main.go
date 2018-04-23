package main

import (
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"os"
	"net"
)

var (
	statusLineRegexp = regexp.MustCompile(`(?m)^(.*):\s+(.*)$`)
	fpmStatusURL     = ""
)

func main() {
	url := flag.String("status-url", "", "PHP-FPM status URL")
	addr := flag.String("unix-sock", "/dev/shm/php-fpm_exporter.sock", "unix sock for access")
	flag.Parse()

	if *url == "" {
		log.Fatal("The status-url flag is required.")
	} else {
		fpmStatusURL = *url
	}

	scrapeFailures := 0
	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		resp, err := http.Get(fpmStatusURL)

		if err != nil {
			log.Println(err)
			scrapeFailures = scrapeFailures+1
			x := strconv.Itoa(scrapeFailures)
			NewMetricsFromMatches([][]string{{"scrape failure:","scrape failure",x}}).WriteTo(w)
			return
		}

		if (resp.StatusCode != http.StatusOK){
			log.Println("php-fpm status code is not OK.")
			scrapeFailures = scrapeFailures+1
			x := strconv.Itoa(scrapeFailures)
			NewMetricsFromMatches([][]string{{"scrape failure:","scrape failure",x}}).WriteTo(w)
			return
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			scrapeFailures = scrapeFailures+1
			x := strconv.Itoa(scrapeFailures)
			NewMetricsFromMatches([][]string{{"scrape failure:","scrape failure",x}}).WriteTo(w)
			return
		}

		resp.Body.Close()

		x := strconv.Itoa(scrapeFailures)

		matches := statusLineRegexp.FindAllStringSubmatch(string(body), -1)
		matches = append(matches,[]string{"scrape failure:","scrape failure",x})

		NewMetricsFromMatches(matches).WriteTo(w)
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>php-fpm Exporter</title></head>
             <body>
             <h1>php-fpm Exporter</h1>
             <p><a href='` + "/metrics" + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	server := http.Server{
		Handler: mux, // http.DefaultServeMux,
	}
	os.Remove(*addr)

	listener, err := net.Listen("unix", *addr)
	if err != nil {
		panic(err)
	}
	server.Serve(listener)
}
