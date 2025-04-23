#

all: N4L N4L-db searchN4L http_server

N4L: N4L.go
	go build -o $@ $@.go

N4L-db: N4L-db.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

searchN4L: searchN4L.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

http_server: http_server.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

clean:
	rm -f N4L N4L-db
	rm -f *~ demo_pocs/*~

