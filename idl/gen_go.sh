#!/bin/bash
rm -rf gen-go

thrift -r --gen go base.thrift
thrift -r --gen go mars.thrift

rm -rf ../src/bc
cp -rf gen-go/* ../src/
rm -rf gen-go
