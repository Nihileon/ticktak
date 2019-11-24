#!/bin/bash

RUN_NAME="ticktak_server"
mkdir -p output/conf output/script output/log
go build -a -o output/${RUN_NAME}
cp -r conf/*  output/conf/
cp -r script/* output/script