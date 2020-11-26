
# Terraform Provisioner for Biome
A [Terraform](https://terraform.io) provisioner to install and configure [Biome](https://Biome.sh)

This is a development repository for adding features to the [provisioner](https://github.com/hashicorp/terraform/tree/master/builtin/provisioners/Biome). Primary goal at the moment is to add support for installation and configuration of Biome on Windows.

Contents
--------

* [Installation](#installation)
* [Example](#example)
* [Documentation](#documentation)


Installation
------------

### Requirements

This Biome provisioner is based on [Terraform](https://www.terraform.io), this means that you need


- [Terraform](https://www.terraform.io/downloads.html) >=0.11.11
- You must have winrm access to Windows machine with admin rights for installation

### Installation from binaries (recommended)

The recommended way to install Biome provisioner is to use the binary distributions from the [Releases](https://github.com/iris3th/terraform-provisioner-Biome/releases) page. The packages are available for macOS, Windows and Linux.

Download and uncompress the latest release for your OS. This example uses the macOS binary.

```sh
> wget https://github.com/iris3th/terraform-provisioner-Biome/releases/download/0.1/terraform-provisioner-Biome_dev_v0.1_darwin_amd64.tar.gz
> tar -xvf terraform-provisioner-Biome_dev*.tar.gz
```

Now copy the binary to the Terraform's plugins folder, if is your first plugin maybe isn't present.

```sh
> mkdir -p ~/.terraform.d/plugins/
> mv terraform-provisioner-Biome_dev*/terraform-provisioner-Biome_dev* ~/.terraform.d/plugins/
```

### Build from source

If you wish to compile the provisioner from source code, you'll first need [Go](http://www.golang.org) installed on your machine (version >=1.11.4 is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`. In addition you will need to install required Terraform packages on your machine from Hashicorp Terraform [repository](https://github.com/hashicorp/terraform).

Clone Provisioner repository to: `$GOPATH/src/github.com/iris3th/terraform-provisioner-Biome`

```
cd $GOPATH/src/github.com/iris3th/
git clone https://github.com/iris3th/terraform-provisioner-Biome.git
cd terraform-provisioner-Biome
go build -o terraform-provisioner-Biome_dev
mv ./terraform-provisioner-Biome_dev ~/.terraform.d/plugins/terraform-provisioner-Biome_dev
chmod +x ~/.terraform.d/plugins/terraform-provisioner-Biome_dev
```


Example
--------
The following example shows how to install SQL Server Biome service on a Windows machine

```
provisioner "Biome_dev" {
    peer = ""

    service {
      name = "core/sqlserver"
      topology = "standalone"
    }

    connection {
      type = "winrm"
      timeout = "10m"
      user = "${var.win_username}"
      password = "${var.win_password}"
    }
}
```

To attach this SQL server with an application as Biome peer and establish service binding, the following can be used in the application server's provisioning

```
provisioner "Biome_dev" {
    peer = "${aws_instance.ms-hab-sqlserver.private_ip}"

    service {
      name = "mwrock/contosouniversity"
      topology = "standalone"
      binds = [
        "database:sqlserver.default"
      ]
    }

    connection {
      type = "winrm"
      timeout = "15m"
      user = "${var.win_username}"
      password = "${var.win_password}"
    }
  }
 ```

Documentation
--------
## Arguments
There are 2 configuration levels, supervisor and service.  Values placed directly within the `provisioner` block are supervisor configs, and values placed inside a `service` block are service configs.  Services can also take a `bind` block to configure runtime bindings.

### Supervisor
* `permanent_peer`: Whether this supervisor should be marked as a permanent peer. Optional (Defaults to false)
* `listen_gossip`: IP and port to listen for gossip traffic.  Optional (Defaults to "0.0.0.0:9638")
* `listen_http`: IP and port for the HTTP API service.  Optional (Defaults to "0.0.0.0:9631")
* `peer`: IP or FQDN of a supervisor instance to peer with.  Optional (Defaults to none)
* `ring_key`: Key for encrypting the supervisor ring traffic.  Optional (Defaults to none)
* `skip_install`: Skips the installation Biome, if it's being installed another way.  Optional (Defaults to no)
### Service
* `name`: A package identifier of the Biome package to start (eg `core/nginx`, `core/nginx/1.11.10` or `core/nginx/1.11.10/20170215233218`).  Required.
* `strategy`: Update strategy to use. Possible values "at-once", "rolling" or "none".  Optional (Defaults to "none")
* `topology`: Topology to start service in.  Possible values "standalone" or "leader".  Optional (Defaults to "standalone")
* `channel`: Channel in a remote depot to watch for package updates.  Optional
* `group`: Service group to join.  Optional (Defaults to "default")
* `url`: URL of the remote Depot to watch.  Optional (Defaults to the public depot)
* `binds`:  Array of binding statements (eg "backend:nginx.default").  Optional
* `user_toml`: TOML formatted user configuration for the service.  Easiest to source from a file (eg `user_toml = "${file("conf/redis.toml")}"`).  Optional

### Bind
* `service`: The target service to bind.
* `group`: The target group to bind.
* `alias`: The alias for the binding.

**This format for declaring bindings is optional.  It can be used in place of or along side the `binds = ["alias:service.group"]` method of declaring binds.  This format might be easier to manage when populating one or more of the bind parameters dynamically.

Example:
```
service {
  name = "core/haproxy"
  group = "${var.environment}"

  bind {
    alias = "backend"
    service = "nginx"
    group = "${var.environment}"
  }
}
```
This block will generate the option `--bind backend:nginx.default` when starting the haproxy service.


## Kudos & Prior Art

This provisioner is heavily inspired by prior art from [nsdavidson/terraform-provision-Biome](https://github.com/nsdavidson/terraform-provisioner-Biome) later added to [Terraform](https://github.com/hashicorp/terraform/pull/16280).
