gostatsd
========

An implementation of [Etsy's][etsy] [statsd][statsd] in Go.

The project provides both a server called "gostatsd" which works much like
Etsy's version, but also provides a library for developing customized servers.

Backends are pluggable and only need to support the [backend interface](backend/backend.go).  


Building the server
-------------------
From the `gostatsd/` directory run `go build`. The binary will be built in place
and called `gostatsd`


Running the server
------------------
`gostatsd --help` gives a complete description of available options and their
defaults.

Configuring the backends
------------------------
Backends are configured using `toml`, `json` or `yaml` configuration file passed through
the `--config-path` flag, see [example/config.toml](example/config.toml).


Sending metrics
---------------
The server listens for UDP packets on the address given by the `--metrics-addr` flag,
aggregates them, then sends them to the backend servers given by the `--backends`
flag (comma separated list of backend names).

Currently supported backends are:

* graphite
* stdout

The format of each metric is:

    <bucket name>:<value>|<type>\n

* `<bucket name>` is a string like `abc.def.g`, just like a graphite bucket name
* `<value>` is a string representation of a floating point number
* `<type>` is one of `c`, `g`, or `ms` for "counter", "gauge", and "timer"
respectively.

A single packet can contain multiple metrics, each ending with a newline.

Optionally, `gostatsd` supports sample rates and tags (unused):

* `<bucket name>:<value>|c|@<sample rate>\n` where `sample rate` is a float between 0 and 1
* `<bucket name>:<value>|c|@<sample rate>|#<tags>\n` where `tags` is a comma separated list of tags
* or `<bucket name>:<value>|<type>|#<tags>\n` where `tags` is a comma separated list of tags

Tags format is: `simple` or `key:value`.


A simple way to test your installation or send metrics from a script is to use
`echo` and the [netcat][netcat] utility `nc`:

    echo 'abc.def.g:10|c' | nc -w1 -u localhost 8125

Monitoring
----------
Currently you can get some basic idea of the status of the server by visiting the
address given by the `-c` option with your web browser.

Contributing
------------
Contribute more backends by sending pull requests.


Using the library
-----------------
In your source code:

    import "github.com/jtblin/gostatsd/statsd"

Documentation can be found via `go doc github.com/jtblin/gostatsd/statsd` or at
http://godoc.org/github.com/jtblin/gostatsd/statsd

[etsy]: http://www.etsy.com
[statsd]: http://www.github.com/etsy/statsd
[netcat]: http://netcat.sourceforge.net/