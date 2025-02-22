
all: src/N4L.go test
	cd src; make

test: src/N4L
	(cd src; make)
	(cd tests; make)
clean:
	rm *~
	(cd src; make clean)
	(cd examples; make clean)

% : %.go
	go build $<
