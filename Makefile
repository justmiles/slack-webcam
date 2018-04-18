
build:
	rm -rf build
	mkdir -p build
	go build main.go
	mv main build/slack-webcam.linux-amd64
