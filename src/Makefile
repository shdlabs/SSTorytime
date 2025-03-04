#

all: N4L N4L-db

N4L: N4L.go
	go build -o $@ $@.go

N4L-db: N4L-db.go ../pkg/SSTorytime/SSTorytime.go
	go build -o $@ $@.go

clean:
	rm -f N4L N4L-db
	rm -f *~ demo_pocs/*~

