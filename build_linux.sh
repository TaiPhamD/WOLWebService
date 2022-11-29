#!/bin/bash

# create build/dist if it doesn't exist
if [ -d build ]; then
    rm build/dist/wolwebservice
else
    mkdir -p build/dist
fi
pushd server
go build -o ../build/dist
echo "compiled binary is located in build/dist"
# pop back to the root folder
popd

