httpcat
===============

A netcat like tool for analysing HTTP requests (from both server and client perspective).  My first experiment with Go.


		Usage:

		  Client mode
			 httpcat -client [options] http://uri-to-send-request-to.com

		  Server mode
			 httpcat -server [options]

		  Options (either mode)
			 -complete or -c : Display entire request/response instead of just the body.
			 -verbose or -v : Be verbose.

		  Options (client mode only)
			 -accept or -a [accept string] : Adds 'Accept' header to request.

		  Options (server mode only)
			 -body or -b [body message] : Body to respond with.  Response code will default to 200.
			 -port or -p [port] : Port to listen on.
			 -response or -r [response code] : Status code to respond with.  Defaults to 204.
			 -separator or -s [separator string] : Use the provided separator to separate messages.


