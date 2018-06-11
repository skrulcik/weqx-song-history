setup:
	go get 'golang.org/x/net/html'

run:
	go run song-history.go

.PHONY: setup run
