#!/bin/bash

rm dtp_*.go

cd ./generator

go run . browser_protocol.json js_protocol.json


go fmt ./files

cp ./files/*.* ./../

