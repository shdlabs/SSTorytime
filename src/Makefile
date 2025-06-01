#

OBJ=N4L N4L-db searchN4L http_server pathsolve notes graph_report API_EXAMPLE

all: $(OBJ)

N4L: N4L.go
	go build -o $@ $@.go

N4L-db: N4L-db.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

searchN4L: searchN4L.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

http_server: http_server.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

notes: notes.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

pathsolve: pathsolve.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

graph_report: graph_report.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

API_EXAMPLE: API_EXAMPLE.go ../pkg/SSTorytime/SSTorytime.go

	go build -o $@ $@.go

clean:
	rm -f $(OBJ)
	rm -f *~ demo_pocs/*~

