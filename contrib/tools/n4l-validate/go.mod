module github.com/markburgess/SSTorytime/contrib/tools/n4l-validate

go 1.25.3

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/markburgess/SSTorytime/contrib/internal/n4l v0.0.0
)

require golang.org/x/sys v0.13.0 // indirect

replace github.com/markburgess/SSTorytime/contrib/internal/n4l => ../../internal/n4l
