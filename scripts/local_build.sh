#!/usr/bin/env bash

cd ..
go build -o terraform-provisioner-biome_dev
mv ./terraform-provisioner-biome_dev ~/.terraform.d/plugins/terraform-provisioner-biome_dev
chmod +x ~/.terraform.d/plugins/terraform-provisioner-biome_dev
