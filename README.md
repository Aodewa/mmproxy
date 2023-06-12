# mmproxy

This is a Go reimplementation of [mmproxy](https://github.com/cloudflare/mmproxy), created to improve on mmproxy's runtime stability while providing potentially greater performance in terms of connection and packet throughput.

`mmproxy` is a standalone application that unwraps HAProxy's [PROXY protocol](http://www.haproxy.org/download/1.8/doc/proxy-protocol.txt) (also adopted by other projects such as NGINX) so that the network connection to the end server comes from client's - instead of proxy server's - IP address and port number.
Because they share basic mechanisms, [Cloudflare's blogpost on mmproxy](https://blog.cloudflare.com/mmproxy-creative-way-of-preserving-client-ips-in-spectrum/) serves as a great write-up on how `mmproxy` works under the hood.

## Install

```shell
sudo ./install.sh
```

## Requirements

`mmproxy` has to be ran:

- on the same server as the proxy target, as the communication happens over the loopback interface;

## Running

### Modify mmproxy config

The config file path is /etc/mmproxy.json, you should edit it before running mmproxy.
For example, your server is listening at port 2222, and mmproxy is listening at port 1111, you can find it as:
```
{
    "router": [
        {
            "from":"0.0.0.0:1111",
            "to":"127.0.0.1:2222"
        }
    ],
    "p":"tcp",
    "mark":0,
    "v":0,
    "listeners":1,
}
```

The config filed is:
```
  router []subroute
    from string
        Address the proxy listens on
    to string
        Address to which IPv4 traffic will be forwarded to, the IP must be "127.0.0.1"
  listeners int
    	Number of listener sockets that will be opened for the listen address (Linux 3.9+) (default 1)
  mark int
    	The mark that will be set on outbound packets (default 0)
  p string
    	Protocol that will be proxied: tcp, udp (default "tcp")
  v int (default 0)
    	0 - no logging of individual connections
    	1 - log errors occurring in individual connections
    	2 - log all state changes of individual connections
```

### Start mmproxy

```
sudo systemctl start mmproxy
```

If you want to set as autostart:
```
sudo systemctl enable mmproxy
```

### Stop mmproxy

```
sudo systemctl stop mmproxy
```