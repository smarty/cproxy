To Run
-----------------------------

**NOTE: This example requires the Go runtime to be installed on your machine.**

`cd` to the the `examples/basic` directory and run the following in a terminal window:
```
$ go run main.go
```

Then, in another terminal window:
```
$ export https_proxy="http://127.0.0.1:8080/"
$ curl -v https://us-street.api.smartystreets.com/status # any secure URL
```

Now let's perform a direct connection to the server so we can observe that the TLS/SSL x509 certificate from the previous request and the following request are identical.
```
$ unset https_proxy
$ curl -v https://us-street.api.smartystreets.com/status
```

Production Use
------------------------------
```
$ go build -o cproxy
$ scp cproxy user@production-server.com:.
```
NOTE: While the above compiled binary is ready for deployment into production, this isn't the best strategy. The server has no filtering capabilities and is the equivalent of an "open relay". Without proper filtering, you've just opened up an anonymous proxy server that anyone can connect to and can connect to anything. It's a really bad idea. That's why `cproxy` is a library meant to be extended with custom `Filter` implementations to only allow trusted traffic through.