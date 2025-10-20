module main

go 1.24.2

replace SSTorytime => ../pkg/SSTorytime

require (
	SSTorytime v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.9
	golang.org/x/text v0.24.0
)

require (
	github.com/arl/statsviz v0.7.2 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/lmittmann/tint v1.1.2 // indirect
)
