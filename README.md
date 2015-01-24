# snowreport
Use the noaa.gov weather API to pull forecasts from east coast ski destinations.

This is by no means the best setup for this program.
It is partially an experiment with a couple of tools.

### installing

Prerequisistes:
* the Go runtime
* A MySQL instance
* An etcd instance

Install golang prerequisites with `go get -t -v ./`

Setup the MySQL schema with [etc/schema.sql](https://github.com/phone/snowreport/blob/master/etc/schema.sql), into an empty database.

Etcd stores the mysql connection information. Keys required are:
* `/snowreport/mysql/user`
* `/snowreport/mysql/password`
* `/snowreport/mysql/host`
* `/snowreport/mysql/port`
* `/snowreport/mysql/db`

They do what you think. Add them with e.g. `curl -L http://localhost:4001/v2/keys/snowreport/mysql/user -XPUT -d value="snowreport"`

Exercise the program with `go test -v`