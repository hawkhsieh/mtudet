Here's a Go implementation of an mtu detector. Sometimes, mtu can be smaller or larger than 1500. This program allows you to use it on your embedded devices. Its advantage is that it can be easily compiled to ARM Linux without needing to find a toolchain or deal with complex Makefiles.



mac:

docker run --rm -it -v "$PWD":/go/src/mtudet -w /go/src/mtudet golang:1.18 bash -c "GOOS=darwin GOARCH=amd64 go build"

arm:

docker run --rm -it -v "$PWD":/go/src/mtudet -w /go/src/mtudet golang:1.18 bash -c "GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build"
