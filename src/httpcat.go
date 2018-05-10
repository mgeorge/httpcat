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
	port      int
	status    int
	body      string
	entire    bool
	verbose   bool
	server    bool
	uri       string
	accept    string
	separator string
	cors      bool
)

func requestHandler(resp http.ResponseWriter, req *http.Request) {

	if verbose { fmt.Printf("Received a %s request to %s\n", req.Method, req.RequestURI) }

	// print entire request
	if entire {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump[:]))

	// print only request body
	} else {
		body, _ := ioutil.ReadAll(req.Body)
		fmt.Println(string(body[:]))
	}

	// cors
	if cors && req.Method == http.MethodOptions {
		if verbose { fmt.Println("CORS preflight") }

		if origin := req.Header.Get("Origin"); origin != "" {

			resp.Header().Set("Access-Control-Allow-Origin", origin);

			if headers := req.Header.Get("Access-Control-Request-Headers"); headers != "" {
				resp.Header().Set("Access-Control-Allow-Headers", headers);
			}

			if method := req.Header.Get("Access-Control-Request-Method"); method != "" {
				resp.Header().Set("Access-Control-Allow-Method", method);
			}

			resp.WriteHeader(200);

			return;

		} else {
			if verbose { fmt.Println("Preflight did not contain 'Origin' header!") }
		}
	}

	// send response
	if body != "" {
		if status == 204 {status = 200}  // if we have a body then we don't want to return the default 204 status
		resp.WriteHeader(status)
		fmt.Fprintf(resp, body)
	} else {
		resp.WriteHeader(status)
	}

	if separator != "" {
		fmt.Println(separator)
	}
}

func startServer() {
	http.HandleFunc("/", requestHandler)

	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		fmt.Fprintf(os.Stderr, "Could not start server.  Is port %d available?\n", port)
	}
}

func sendRequest(uri string) {

	client := http.Client{}

	request, _ := http.NewRequest("GET", uri, nil)

	// add Accept header if required by user
	if accept != "" {
		request.Header.Add("Accept", accept)
	}

	response, err := client.Do(request)

	if(err != nil) {
		fmt.Fprintln(os.Stderr, err)
	} else {

		// print entire response
		if entire {
			dump, _ := httputil.DumpResponse(response, true)
			fmt.Println(string(dump[:]))

		// print only the response body
		} else {
			body, _ := ioutil.ReadAll(response.Body)
			fmt.Println(string(body[:]))
		}
	}
}

func parseCommandLine() {
	flag.IntVar(&port, "port", 8080, "")
	flag.IntVar(&port, "p", 8080, "")

	flag.IntVar(&status, "response", 204, "")
	flag.IntVar(&status, "r", 204, "")

	flag.StringVar(&body, "body", "", "")
	flag.StringVar(&body, "b", "", "")

	flag.BoolVar(&entire, "entire", false, "")
	flag.BoolVar(&entire, "e", false, "")

	flag.StringVar(&separator, "separator", "", "")
	flag.StringVar(&separator, "s", "", "")

	flag.BoolVar(&verbose, "verbose", false, "")
	flag.BoolVar(&verbose, "v", false, "")

	flag.StringVar(&accept, "accept", "", "")
	flag.StringVar(&accept, "a", "", "")


	flag.BoolVar(&cors, "cors", false, "")
	flag.BoolVar(&cors, "c", false, "")

	serverMode := flag.Bool("server", false, "")
	clientMode := flag.Bool("client", false, "")

	flag.Usage = usage

	flag.Parse()

	if !(*serverMode || *clientMode) || (*serverMode && *clientMode) {
		usage()
		os.Exit(1)
	} else {
		server = *serverMode
	}

	if *clientMode {
		uri = flag.Arg(0)
		if uri == "" {
			usage()
			os.Exit(1)
		}
	}

}

var usage = func() {
	fmt.Fprintf(os.Stderr, "Usage:\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Client mode\n")
	fmt.Fprintf(os.Stderr, "    httpcat -client [options] http://uri-to-send-request-to.com\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "    Currently only supports GET requests.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Server mode\n")
	fmt.Fprintf(os.Stderr, "    httpcat -server [options]\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Options (either mode)\n")
	fmt.Fprintf(os.Stderr, "    -entire or -e : Display entire request/response instead of just the body.\n")
	fmt.Fprintf(os.Stderr, "    -verbose or -v : Be verbose.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Options (client mode only)\n")
	fmt.Fprintf(os.Stderr, "    -accept or -a [accept string] : Adds 'Accept' header to request.\n")
	fmt.Fprintf(os.Stderr, "\n")
	fmt.Fprintf(os.Stderr, "  Options (server mode only)\n")
	fmt.Fprintf(os.Stderr, "    -body or -b [body message] : Body to respond with.  Response code will default to 200.\n")
	fmt.Fprintf(os.Stderr, "    -port or -p [port] : Port to listen on.\n")
	fmt.Fprintf(os.Stderr, "    -response or -r [response code] : Status code to respond with.  Defaults to 204.\n")
	fmt.Fprintf(os.Stderr, "    -cors or -c : Enable Cross Origin Resource Sharing support.\n")
	fmt.Fprintf(os.Stderr, "    -separator or -s [separator string] : Use the provided separator to separate messages.\n")
}

func main() {

	parseCommandLine()

	if server {

		// server mode

		if verbose {
			fmt.Printf("Listening on port %d.\n", port)
			if entire { fmt.Println("Displaying entire request details.\n") } else { fmt.Println("Displaying bodies only.\n") }
		}

		startServer()

	} else {

		// client mode

		if verbose { fmt.Println("Sending GET request to " + uri) }
		sendRequest(uri)
	}
}
