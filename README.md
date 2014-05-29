httpcat
===============

A netcat like tool for analysing HTTP requests (from both server and client perspective).  My first experiment with Go.

    Usage:

      Client mode

        httpcat -get [options] http://uri-to-send-request-to.com


      Server mode

        httpcat -listen [options]


      Options (client or server mode)
        -complete or -c : Display entire request/response instead of just the body.
        -verbose or -v : Be verbose.

      Options (server mode only)
        -port or -p [port] : Port to listen on.
        -body or -b [body message] : Body to respond with.  Status will default to 200.
        -status or -s [status code] : Status code to respond with.  Defaults to 204.