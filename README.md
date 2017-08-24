Purpose
-----------------------

Oh no! Not another forward proxy! Don't we already have enough forward proxies with full-featured systems such as Squid, mitmproxy, tinyproxy, Apache Traffic Server, privoxy, and even nginx, let alone custom packages written in Python, Node, C#, and every other language under the sun?!

You're absolutely correct. We are drowning in choices for proxies. But like every solution, we have a very specific use case and none of the existing options that we found could be configured to work in the way we needed.

Proxy Basics
-----------------------
There are different kinds of proxies: forward proxies like Squid (often just called proxies) and reverse proxies like HAProxy. Further, you can break the forward proxies grouping down into two basic camps: implicit and explicit. An implicit proxy will silently proxy all content to and from the client application while an explicit proxy must be specifically configured within the client application to shuttle the request. Different corporate firewall policies may have one or the other and sometimes both. For example, large organizations may have an implicit proxy which decrypts SSL/TLS traffic by having a root certificate installed on all client machines. When was the last time you looked carefully the trusted root certificate store on your primary workstation?

Explicit forward proxies can work in a few different ways, but one of the most basic is by establishing an `HTTP CONNECT` tunnel. Configured this way the browser/client connects to the proxy using HTTP but instead of performing a `GET`, `POST`, etc. it performs a `CONNECT`, e.g. `HTTP/1.1 CONNECT some-remote-site.com:443`. Once the proxy inspects that header (which is about the only thing it has), it will authorize the client saying `HTTP/1.1 200 OK`. At that point, the client will start streaming TCP packets which are routed through to the remote/upstream server specified in the `CONNECT` verb.

Here's another pretty cool part about explicit forward proxies. The client can connect to the proxy via plaintext HTTP and perform the `CONNECT` upgrade and then stream TLS packets through the plaintext connection with perfect confidence that the connection to the specified server remains secure. The only information leak would be the site that you're connecting to which would often be visible anyway to anyway analyzing the traffic because you'd be connecting to a specific IP address. In other words, if you connect to a proxy over plaintext, but you're connecting to an HTTPS resource, you can be confident that your connection is secure—so long as you trust the certificate authority issuing the certificate (which is the foundational premise of public key infrastructure).

Why `cproxy`?
-------------------------

Now we arrive at the main raison d'être for `cproxy`: proxying TLS connections with the `CONNECT` verb without losing the client IP address. How do you maintain the client IP information while at the same time proxying TLS traffic? If you can't manipulate/modify the traffic without breaking TLS, how do you add the standard `X-Forwarded-For` HTTP header to the request? Simple answer: you can't.

Okay, what if we decode the TLS traffic and append the header? This is the solution offered by Squid, mitmproxy, privoxy, along with numerous other proxies. If we decrypt the traffic that also means we have to re-encrypt the traffic. It's slight more expensive for the CPU but also has one other really nasty side effect: we have to fake the remote/upstream server certificate. This is typically done by creating a self-signed root certificate and having the proxy use that certificate to dynamically generate domain-specific on the fly for each request. So why doesn't that work? It could, but you'd have to install that self-signed root certificate on every client that wants to connect through the proxy. If that's not an option, we need another solution.

So how do we add the client IP address to a request we don't modify? We do it at the TCP layer using HAProxy's PROXY Protocol. This is the missing piece from the myriad of other proxying solutions. They haven't yet implemented **outbound** connections using the PROXY Protocol standard. Sure, many of them have recently implemented and can **receive** traffic and understand that protocol to maintain the client IP, but sending according to that protocol is another matter altogether.

**`cproxy` is an explicit, forward `HTTP CONNECT` proxy which can connect to remote systems using the HAProxy PROXY Protocol v1 standard.**

How to Use:
---------------------
While `cproxy` can and ultimately should be an executable, it's primary design is that of a library upon which to extend its forward proxy capabilities. It is nothing more than a `http.Handler` ready to be wrapped and hosted in a process.

To create a `http.Handler` that's ready for traffic configure it this way:
```
handler := cproxy.Configure().Build()
```

The above snippet configures `cproxy` as a basic `HTTP CONNECT` proxy with no specific abilities. If you would like to append the PROXY Protocol header to oubound traffic (which is the reason why it was written), configure it like so:
```
handler := cproxy.Configure().WithProxyProtocol().Build()
```

The handler created is a vanilla `http.Handler` that can be attached to a `http.Server` instance so that it can now receive and proxy traffic.

If you would like to filter traffic based upon custom conditions, you can create an implementation of the `Filter` interface to fully inspect the HTTP request to determine if, for example, the appropriate headers are present and sufficient to allow the `CONNECT` process to proceed. Additionally, you could inspect the destination host, and/or determine if a username/password/token is present as part of the request. You could even limit traffic to a known set of client IPs. All of that behavior can be created by implementing the `Filter` interface and configuring the handler as follows:
```
filter := &MyCustomFilter{}
handler := cproxy.Configure().WithFilter(filter).Build()

http.ListenAndServe(":8080", handler) // now listen for traffic
```

A full, working example ready for compilation and deployment to production environments can be found in the `examples/basic` directory.


Naming `cproxy`
----------------
Why is it called `cproxy`? It's short for `CONNECT` proxy. The name is meant to express the idea of an explicit, forward `HTTP CONNECT` proxy.