#!/bin/bash

go build -o $OUTPUT/main ./example/main.go

cp -r ./conf $OUTPUT/