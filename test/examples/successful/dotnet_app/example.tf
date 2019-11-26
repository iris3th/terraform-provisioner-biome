provider "aws" {
  region     = "${var.region}"
  access_key = "${var.access_key}"
}

resource "aws_instance" "ms-hab-sqlserver" {
  ami = "${lookup(var.win_amis, var.region)}"
  instance_type = "t2.large"
  key_name = "${var.key_name}"
  tags = {
    Name = "ms-hab-sqlserver"
  }
  user_data = <<EOF
  <powershell>
  net user ${var.win_username} '${var.win_password}' /add /y
  net localgroup administrators ${var.win_username} /add
  winrm quickconfig -q
  winrm set winrm/config/winrs '@{MaxMemoryPerShellMB="300"}'
  winrm set winrm/config '@{MaxTimeoutms="1800000"}'
  winrm set winrm/config/service '@{AllowUnencrypted="true"}'
  winrm set winrm/config/service/auth '@{Basic="true"}'
  netsh advfirewall firewall add rule name="WinRM 5985" protocol=TCP dir=in localport=5985 action=allow
  netsh advfirewall firewall add rule name="WinRM 5986" protocol=TCP dir=in localport=5986 action=allow
  net stop winrm
  sc.exe config winrm start=auto
  net start winrm
  </powershell>
  EOF

  provisioner "habitat_dev" {
    peer = ""
    accept_license = true

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
}

resource "aws_instance" "ms-hab-appserver" {
  ami = "${lookup(var.win_amis, var.region)}"
  instance_type = "t2.large"
  key_name = "${var.key_name}"
  tags = {
    Name = "ms-hab-appserver"
  }
  user_data = <<EOF
  <powershell>
  net user ${var.win_username} '${var.win_password}' /add /y
  net localgroup administrators ${var.win_username} /add
  winrm quickconfig -q
  winrm set winrm/config/winrs '@{MaxMemoryPerShellMB="300"}'
  winrm set winrm/config '@{MaxTimeoutms="1800000"}'
  winrm set winrm/config/service '@{AllowUnencrypted="true"}'
  winrm set winrm/config/service/auth '@{Basic="true"}'
  netsh advfirewall firewall add rule name="WinRM 5985" protocol=TCP dir=in localport=5985 action=allow
  netsh advfirewall firewall add rule name="WinRM 5986" protocol=TCP dir=in localport=5986 action=allow
  net stop winrm
  sc.exe config winrm start=auto
  net start winrm
  setx HAB_FEAT_INSTALL_HOOK ON /m
  </powershell>
  EOF

  provisioner "habitat_dev" {
    peer = "${aws_instance.ms-hab-sqlserver.private_ip}"
    accept_license = true
    
    service {
      name = "mwrock/contosouniversity"
      topology = "standalone"
      user_toml = "${file("conf/user.toml")}"
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
}

output "App-Server-IP" {
  value = ["${aws_instance.ms-hab-appserver.public_ip}"]
}
output "SQL-Server-IP" {
  value = ["${aws_instance.ms-hab-sqlserver.public_ip}"]
}


