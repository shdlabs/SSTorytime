module github.com/markburgess/SSTorytime/contrib/tools/godoc2n4l

go 1.25.3

require github.com/markburgess/SSTorytime/contrib/internal/scraper v0.0.0

require (
	github.com/PuerkitoBio/goquery v1.10.3 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	golang.org/x/net v0.39.0 // indirect
)

replace github.com/markburgess/SSTorytime/contrib/internal/scraper => ../../internal/scraper
