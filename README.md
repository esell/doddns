# doddns

## Purpose:

The purpose of this application is to dynamically update DNS A records, specifically subdomains, that are managed by Digital Ocean. The use case for this is if you have something 
running at home or somewhere where the IP changes often but you want to be able to hit it via DNS. Basically something like [DynDNS](http://dyn.com/remote-access/) or [Afraid.org](http://freedns.afraid.org/) 
but instead using Digital Ocean.


## Prerequisites:

* The domasubdomain) mxist in the Digital Ocean DNS system
* You need to be able to make outbound HTTP/HTTPS requests


## Usage:

`go get github.com/esell/doddns`

`./doddns -s SUBDOMAIN -d DOMAIN.COM -k DO_API_KEY`

At this point you could set the app to run via cron or whatever.


## Gotchas:
* Currently the default TTL that Digital Ocean gives you is 1800. Of course this is less than ideal for dynamic updates. The current Digital Ocean API does not provide a way to set the TTL 
  on update or addition so you will need to manually set this in their portal. The nice thing is that you only have to set it once, future updates will not reset the TTL you have defined.
