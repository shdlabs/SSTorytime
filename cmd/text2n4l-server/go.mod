module text2n4l-server

go 1.21

replace github.com/shdlabs/SSTorytime/services/text2n4l => ../../services/text2n4l

replace github.com/shdlabs/SSTorytime/services/sstorytime => ../../services/sstorytime

require github.com/shdlabs/SSTorytime/services/text2n4l v0.0.0-00010101000000-000000000000

require (
	github.com/lib/pq v1.10.9 // indirect
	github.com/shdlabs/SSTorytime/services/sstorytime v0.0.0-00010101000000-000000000000 // indirect
)
