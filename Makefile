
all: cmd/N4L test

cmd/bin/N4L:
	(cd cmd; make)
	(cd cmd/demo_pocs; make)

test: cmd/bin/N4L
	(cd cmd; make)
	(cd cmd/demo_pocs; make)
	(cd tests; make)
clean:
	rm -f *~ \#* N4L
	(cd cmd; make clean)
	(cd examples; make clean)
	(cd cmd/demo_pocs; make clean)

ramdisk:
ramdb:
	(cd contrib; sh ramify.sh)
	(cd contrib; sh makeramdb.sh)

db:
	sh contrib/makedb.sh
