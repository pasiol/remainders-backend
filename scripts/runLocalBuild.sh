#!/bin/sh
cd ..
make compile
go build .
../bin/rest-server