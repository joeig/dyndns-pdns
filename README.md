# Dynamic DNS Collector for PowerDNS

Collects IPv4/IPv6 addresses of network devices (routers, firewalls etc.) and writes the corresponding PowerDNS resource records.

[![Build Status](https://travis-ci.org/joeig/dyndns-pdns.svg?branch=master)](https://travis-ci.org/joeig/dyndns-pdns)

## Usage

### Update IP address of a certain host

In order to update the IP address of a certain host, you can choose between to ingest modes:

- Use the value provided by a HTTP GET parameter (IPv4 and/or IPv6)
- Use the value provided by the TCP remote address field (either IPv4 or IPv6, depending on the client's capabilities)

This tool does not support the DDNS protocol (RFC2136), which is supported by PowerDNS out of the box.

#### HTTP GET parameter

Suitable for all common network devices.

~~~ bash
http "https://dyn-ingest.example.com/v1/host/<device name>/sync?key=<key>&ipv4=<IPv4 address>&ipv6=<IPv6 address>"
~~~

You have to provide at least one IP address family.

#### TCP remote address

If your router doesn't know its own egress IP address (might be the most promising solution for people that have to work behind NAT gateways or proxies).

~~~ bash
http "https://dyn-ingest.example.com/v1/host/<device name>/sync?key=<key>"
~~~

This option takes the IP address (IPv4 or IPv6) used by the client during the TCP handshake.

## Examples

### Fritzbox (FritzOS)

Your Fritzbox configuration might look like this:

| Key | Value |
| --- | ----- |
| Provider | Custom |
| Update URL | https://dyn-ingest.example.com/v1/host/\<username\>/sync?key=\<pass\>&ipv4=\<ipaddr\>&ipv6=\<ip6addr\> |
| Domain name | fb-home.dyn.example.com |
| Username | fb-home |
| Password | ... |

You have to copy the update URL as it is, including all the placeholders and \<\> brackets. They will be substituted by FritzOS internally with the corresponding values.
