/**
	Author:  Mark George <mark.george@otago.ac.nz>
	Warranty:  None.  Works for me.  If it doesn't work for you then you already have the source code...
	License: WTFPL v2 <http://www.wtfpl.net/txt/copying/>
*/

package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"flag"
	"strconv"
	"io/ioutil"
	"os"
)

var (
	port     int
	status   int
	body     string
	complete bool
	verbose  bool
	server   bool
	uri      string
)

func requestHandler(resp http.ResponseWriter, req *http.Request) {

	if verbose { fmt.Printf("Received a %s request to %s\n", req.Method, req.RequestURI) }

	// print complete request
	if complete {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump[:]))

	// print only request body
	} else {
		body, _ := ioutil.ReadAll(req.Body)
		fmt.Println(string(body[:]));
	}

	// send response
	if body != "" {
		if status == 204 {status = 200}  // if we have a body then we don't want to return the default 204 status
		resp.WriteHeader(status)
		fmt.Fprintf(resp, body)
	} else {
		resp.WriteHeader(status)
	}

}

func parseCommandLine() {
	flag.IntVar(&port, "port", 8080, "")
	flag.IntVar(&port, "p", 8080, "")

	flag.IntVar(&status, "status", 204, "")
	flag.IntVar(&status, "s", 204, "")

	flag.StringVar(&body, "body", "", "")
	flag.StringVar(&body, "b", "", "")

	flag.BoolVar(&complete, "complete", false, "")
	flag.BoolVar(&complete, "c", false, "")

	flag.BoolVar(&verbose, "verbose", false, "");
	flag.BoolVar(&verbose, "v", false, "");

	serverMode := flag.Bool("listen", false, "");
	clientMode := flag.Bool("get", false, "");

	flag.Usage = usage

	flag.Parse()

	if !(*serverMode || *clientMode) || (*serverMode && *clientMode) {
		usage();
		os.Exit(1)
	} else {
		server = *serverMode
	}

	if *clientMode {
		uri = flag.Arg(0)
		if uri == "" {
			usage();
			os.Exit(1)
		}
	}

}

func startServer() {
	http.HandleFunc("/", requestHandler)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		fmt.Printf("Could not start server.  Is port %d available?\n", port)
	}
}

func sendRequest(uri string) {
		response, err := http.Get(uri)

		if(err != nil) {
			fmt.Println(err);
		} else {
			dump, _ := httputil.DumpResponse(response, true)
			fmt.Println(string(dump[:]))
		}
}

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Client mode\n\n")
	fmt.Fprintf(os.Stderr, "    httpcat -get [options] http://uri-to-send-request-to.com\n")
	fmt.Fprintf(os.Stderr, "\n\n")
	fmt.Fprintf(os.Stderr, "  Server mode\n\n")
	fmt.Fprintf(os.Stderr, "    httpcat -listen [options]\n")
	fmt.Fprintf(os.Stderr, "\n\n")
	fmt.Fprintf(os.Stderr, "  Options (client or server mode)\n")
	fmt.Fprintf(os.Stderr, "    -complete or -c : Display entire request/response instead of just the body.\n")
	fmt.Fprintf(os.Stderr, "    -verbose or -v : Be verbose.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Options (server mode only)\n")
	fmt.Fprintf(os.Stderr, "    -port or -p [port] : Port to listen on.\n")
	fmt.Fprintf(os.Stderr, "    -body or -b [body message] : Body to respond with.  Status will default to 200.\n")
	fmt.Fprintf(os.Stderr, "    -status or -s [status code] : Status code to respond with.  Defaults to 204.\n")
}

func main() {

	parseCommandLine()

	if server {

		// server mode

		if verbose {
			fmt.Printf("Listening on port %d.\n", port)
			if complete { fmt.Println("Displaying complete requests.\n") } else { fmt.Println("Displaying bodies only.\n") }
		}

		startServer()

	} else {

		// client mode

		if verbose { fmt.Println("Sending GET request to " + uri); }
		sendRequest(uri)
	}
}
