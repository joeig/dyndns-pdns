# Dynamic DNS Collector for PowerDNS

Collects IPv4/IPv6 addresses of network devices (routers, firewalls etc.) and writes the corresponding PowerDNS resource records.

[![Build Status](https://travis-ci.org/joeig/dyndns-pdns.svg?branch=master)](https://travis-ci.org/joeig/dyndns-pdns)

## Update IP address of a certain host

In order to update the IP address of a certain host, you can choose between to ingest modes:

- Use the value provided by a HTTP GET parameter (IPv4 and/or IPv6)
- Use the value provided by the TCP remote address field (either IPv4 or IPv6, depending on the client's capabilities)

### HTTP GET parameter

~~~ bash
http "https://dyn.example.com/v1/host/<device name >/sync?key=<key>&ipv4=<IPv4 address>&ipv6=<IPv6 address>"
~~~

You have to provide at least one IP address family.

### TCP remote address

~~~ bash
http "https://dyn.example.com/v1/host/<device name>/sync?key=<key>"
~~~

This option takes the IP address (IPv4 or IPv6) used by the client during the TCP handshake.
