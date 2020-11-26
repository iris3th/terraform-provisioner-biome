#!/usr/bin/bash

WORKROOT=$(pwd)
cd ${WORKROOT}

# unzip go environment
#go_env="go1.12.6.linux-amd64.tar.gz"
#wget -c https://dl.google.com/go/go1.12.6.linux-amd64.tar.gz
#tar -zxf ./$go_env
#if [ $? -ne 0 ];
#then
#    echo "Failure in extracting go"
#    exit 1
#fi
echo "Successfully installed Go"
#rm -rf ./$go_env

# prepare PATH, GOROOT and GOPATH
#export PATH=$(pwd)/go/bin:$PATH
#export GOROOT=$(pwd)/go
#export GOPATH=$(pwd)

# dependencies
echo "terraform stuff"
go get -u github.com/hashicorp/terraform/plugin
go get -u github.com/hashicorp/terraform/terraform
go get -u github.com/hashicorp/terraform/communicator
echo "msgpack"
go get -u github.com/vmihailenco/msgpack
cd ~/go/src/github.com/vmihailenco/msgpack
git checkout v4
mkdir -p ~/go/src/github.com/vmihailenco/msgpack/v4
mv -f ~/go/src/github.com/vmihailenco/msgpack/* ~/go/src/github.com/vmihailenco/msgpack/v4
cd -
go get -u github.com/hashicorp/hcl
cd ~/go/src/github.com/hashicorp/hcl
git checkout hcl2
cd -
go get github.com/vmihailenco/msgpack/v4
go get -u github.com/hashicorp/terraform/configs
go get -u github.com/mitchellh/go-linereader
go get -u github.com/iris3th/terraform-provisioner-habitat/biome

echo "Installed project and dependencies"

# build
cd ${WORKROOT}

go build -o terraform-provisioner-habitat_dev -v
if [ $? -ne 0 ];
then
    echo "Failure in building habitat provisioner"
    exit 1
fi
echo "Successfully built habitat provisoner"
