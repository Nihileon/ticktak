#!/bin/bash

RUN_NAME="ticktak_server"
mkdir -p output/bin output/conf output/script
cp -r conf/*  output/conf/
cp -r script/* output/script

go build -a -o output/bin/${RUN_NAME}