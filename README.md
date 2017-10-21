# nagios-plugin-check_dnsblgo
Check an IP address against a list of DNSBLs

This tool checks a given IP address against a list of DNSBLs (DNS Blocklists) and reports on status

* If an IP address is listed, it will return a Critical response
* If timeouts are received, a Warning response will be returned
* If any supplied arguments are invalid or incorrect, a Critical response is returned

There is no need to run this more than a few times per day

## Dependencies

### Running pre-built executable

No dependencies

### Building from source

* Go (https://golang.org/dl/)
* Build
```
go build check_dnsblgo.go
```


## Installation

* Copy `check_dnsblgo` to your Nagios plugin folder and set executable

* Copy the `sample-blocklists.txt` to a location on your system that the NRPE/Nagios user has read access to

* Add the following to your `nrpe.cfg` file

```
command[check_dnsblgo]=/usr/lib/nagios/plugins/check_dnsblgo --ip=1.2.3.4 --file=/path/to/blocklists.txt
```


## Future Work

None planned


## Changelog

* 2017-10-21 :: 0.1.0   :: First Release


## Site

https://github.com/alasdairkeyes/nagios-plugin-check_dnsblgo


## Author

* Alasdair Keyes - https://akeyes.co.uk/


## License

Released under GPL V3 - See included license file.
