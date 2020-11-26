#!/usr/bin/bash

WORKROOT=$(pwd)
cd ${WORKROOT}

# prepare PATH, GOROOT and GOPATH
export PATH=$(pwd)/go/bin:$PATH
export GOROOT=$(pwd)/go
export GOPATH=$(pwd)

ls -l ${GOPATH}/src/github.com/iris3th/terraform-provisioner-biome/biome
go test ${GOPATH}/src/github.com/iris3th/terraform-provisioner-biome/biome -v
if [ $? -ne 0 ];
then
    echo "Failure in biome provisioner unit tests"
    exit 1
fi
echo "Successfully ran the unit tests for biome provisioner"
