# This how we want to name the binary output
BINARY='rest-server'

# These are the values we want to pass for VERSION and BUILD
# git tag 1.0.1
# git commit -am "One more change after the tags"
VERSION=`git describe --tags`
BUILD=`date +%FT%T%z`

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS_app=-ldflags "-X 'main.Version=${VERSION}' -X main.Build=${BUILD}"

# Builds the project
build:
	go build ${LDFLAGS_app} -o ${BINARY} main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS_app} -o bin/${BINARY}
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS_app} -o bin/${BINARY}.exe