#!/usr/bin/env bash

cd ..
go build -o terraform-provisioner-habitat_dev
mv ./terraform-provisioner-habitat_dev ~/.terraform.d/plugins/terraform-provisioner-habitat_dev
chmod +x ~/.terraform.d/plugins/terraform-provisioner-habitat_dev
