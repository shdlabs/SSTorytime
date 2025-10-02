module github.com/shdlabs/SSTorytime

go 1.24.2

replace github.com/shdlabs/SSTorytime/services/text2n4l => ./services/text2n4l

replace github.com/shdlabs/SSTorytime/services/sstorytime => ./services/sstorytime

require (
	github.com/shdlabs/SSTorytime/services/sstorytime v0.0.0-00010101000000-000000000000
	github.com/shdlabs/SSTorytime/services/text2n4l v0.0.0-00010101000000-000000000000
)

require github.com/lib/pq v1.10.9 // indirect
