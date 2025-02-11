
all: src/N4L.go test
	cd src; make

test: src/N4L
	(cd src; make)
	(cd tests; make)
clean:
	(cd src;	rm N4L)
	(cd examples; rm *_test_log *~)

% : %.go
	go build $<
