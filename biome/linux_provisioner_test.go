package biome

import (
	"github.com/hashicorp/terraform/communicator"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const linuxDefaultSystemdUnitFileContents = `[Unit]
Description=Biome Supervisor

[Service]
ExecStart=/bin/bio sup run --peer host1 --peer 1.2.3.4 --auto-update
Restart=on-failure
[Install]
WantedBy=default.target`

const linuxCustomSystemdUnitFileContents = `[Unit]
Description=Biome Supervisor

[Service]
ExecStart=/bin/bio sup run --listen-ctl 192.168.0.1:8443 --listen-gossip 192.168.10.1:9443 --listen-http 192.168.20.1:8080 --peer host1 --peer host2 --peer 1.2.3.4 --peer 5.6.7.8 --peer foo.example.com
Restart=on-failure
Environment="BIO_SUP_GATEWAY_AUTH_TOKEN=ea7-beef"
Environment="BIO_AUTH_TOKEN=dead-beef"
[Install]
WantedBy=default.target`

func TestLinuxProvisioner_linuxInstallbiome(t *testing.T) {
	cases := map[string]struct {
		Config   map[string]interface{}
		Commands map[string]bool
	}{
		"Installation with sudo": {
			Config: map[string]interface{}{
				"version":     "0.79.1",
				"auto_update": true,
				"use_sudo":    true,
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'curl --silent -L0 https://raw.githubusercontent.com/biome-sh/biome/master/components/bio/install.sh > install.sh'": true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bash ./install.sh -v 0.79.1'":                                                                                          true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio install core/busybox'":                                                                                             true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio pkg exec core/busybox adduser -D -g \"\" bio'":                                                                     true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'rm -f install.sh'":                                                                                                     true,
			},
		},
		"Installation without sudo": {
			Config: map[string]interface{}{
				"version":     "0.79.1",
				"auto_update": true,
				"use_sudo":    false,
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true /bin/bash -c 'curl --silent -L0 https://raw.githubusercontent.com/biome-sh/biome/master/components/bio/install.sh > install.sh'": true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true /bin/bash -c 'bash ./install.sh -v 0.79.1'":                                                                                          true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true /bin/bash -c 'bio install core/busybox'":                                                                                             true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true /bin/bash -c 'bio pkg exec core/busybox adduser -D -g \"\" bio'":                                                                     true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true /bin/bash -c 'rm -f install.sh'":                                                                                                     true,
			},
		},
		"Installation with biome license acceptance": {
			Config: map[string]interface{}{
				"version":        "0.81.0",
				"accept_license": true,
				"auto_update":    true,
				"use_sudo":       true,
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'curl --silent -L0 https://raw.githubusercontent.com/biome-sh/biome/master/components/bio/install.sh > install.sh'": true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bash ./install.sh -v 0.81.0'":                                                                                          true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio install core/busybox'":                                                                                             true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio pkg exec core/busybox adduser -D -g \"\" bio'":                                                                     true,
				"env BIO_LICENSE=accept sudo -E /bin/bash -c 'bio -V'":                                                                                                                                        true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'rm -f install.sh'":                                                                                                     true,
			},
		},
	}

	o := new(terraform.MockUIOutput)
	c := new(communicator.MockCommunicator)

	for k, tc := range cases {
		c.Commands = tc.Commands

		p, err := decodeConfig(
			schema.TestResourceDataRaw(t, Provisioner().(*schema.Provisioner).Schema, tc.Config),
		)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		err = p.linuxInstallbiome(o, c)
		if err != nil {
			t.Fatalf("Test %q failed: %v", k, err)
		}
	}
}

func TestLinuxProvisioner_linuxStartbiome(t *testing.T) {
	cases := map[string]struct {
		Config   map[string]interface{}
		Commands map[string]bool
		Uploads  map[string]string
	}{
		"Start systemd biome with sudo": {
			Config: map[string]interface{}{
				"version":      "0.79.1",
				"auto_update":  true,
				"use_sudo":     true,
				"service_name": "bio-sup",
				"peer":         "--peer host1",
				"peers":        []interface{}{"1.2.3.4"},
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio install core/bio-sup/0.79.1'":                             true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'systemctl enable bio-sup && systemctl start bio-sup'":         true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'mv /tmp/bio-sup.service /etc/systemd/system/bio-sup.service'": true,
			},

			Uploads: map[string]string{
				"/tmp/bio-sup.service": linuxDefaultSystemdUnitFileContents,
			},
		},
		"Start systemd biome without sudo": {
			Config: map[string]interface{}{
				"version":      "0.79.1",
				"auto_update":  true,
				"use_sudo":     false,
				"service_name": "bio-sup",
				"peer":         "--peer host1",
				"peers":        []interface{}{"1.2.3.4"},
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true /bin/bash -c 'bio install core/bio-sup/0.79.1'":                     true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true /bin/bash -c 'systemctl enable bio-sup && systemctl start bio-sup'": true,
			},

			Uploads: map[string]string{
				"/etc/systemd/system/bio-sup.service": linuxDefaultSystemdUnitFileContents,
			},
		},
		"Start unmanaged biome with sudo": {
			Config: map[string]interface{}{
				"version":      "0.81.0",
				"license":      "accept-no-persist",
				"auto_update":  true,
				"use_sudo":     true,
				"service_type": "unmanaged",
				"peer":         "--peer host1",
				"peers":        []interface{}{"1.2.3.4"},
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio install core/bio-sup/0.81.0'":                                                                                true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'mkdir -p /bio/sup/default && chmod o+w /bio/sup/default'":                                                        true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c '(setsid bio sup run --peer host1 --peer 1.2.3.4 --auto-update > /bio/sup/default/sup.log 2>&1 <&1 &) ; sleep 1'": true,
			},

			Uploads: map[string]string{
				"/etc/systemd/system/bio-sup.service": linuxDefaultSystemdUnitFileContents,
			},
		},
		"Start biome with custom config": {
			Config: map[string]interface{}{
				"version":            "0.79.1",
				"auto_update":        false,
				"use_sudo":           true,
				"service_name":       "bio-sup",
				"peer":               "--peer host1 --peer host2",
				"peers":              []interface{}{"1.2.3.4", "5.6.7.8", "foo.example.com"},
				"listen_ctl":         "192.168.0.1:8443",
				"listen_gossip":      "192.168.10.1:9443",
				"listen_http":        "192.168.20.1:8080",
				"builder_auth_token": "dead-beef",
				"gateway_auth_token": "ea7-beef",
				"ctl_secret":         "bad-beef",
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true BIO_AUTH_TOKEN=dead-beef sudo -E /bin/bash -c 'bio install core/bio-sup/0.79.1'":                             true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true BIO_AUTH_TOKEN=dead-beef sudo -E /bin/bash -c 'systemctl enable bio-sup && systemctl start bio-sup'":         true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true BIO_AUTH_TOKEN=dead-beef sudo -E /bin/bash -c 'mv /tmp/bio-sup.service /etc/systemd/system/bio-sup.service'": true,
			},

			Uploads: map[string]string{
				"/tmp/bio-sup.service": linuxCustomSystemdUnitFileContents,
			},
		},
	}

	o := new(terraform.MockUIOutput)
	c := new(communicator.MockCommunicator)

	for k, tc := range cases {
		c.Commands = tc.Commands
		c.Uploads = tc.Uploads

		p, err := decodeConfig(
			schema.TestResourceDataRaw(t, Provisioner().(*schema.Provisioner).Schema, tc.Config),
		)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		err = p.linuxStartbiome(o, c)
		if err != nil {
			t.Fatalf("Test %q failed: %v", k, err)
		}
	}
}

func TestLinuxProvisioner_linuxUploadRingKey(t *testing.T) {
	cases := map[string]struct {
		Config   map[string]interface{}
		Commands map[string]bool
	}{
		"Upload ring key": {
			Config: map[string]interface{}{
				"version":          "0.79.1",
				"auto_update":      true,
				"use_sudo":         true,
				"service_name":     "bio-sup",
				"peers":            []interface{}{"1.2.3.4"},
				"ring_key":         "test-ring",
				"ring_key_content": "dead-beef",
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'echo -e \"dead-beef\" | bio ring key import'": true,
			},
		},
	}

	o := new(terraform.MockUIOutput)
	c := new(communicator.MockCommunicator)

	for k, tc := range cases {
		c.Commands = tc.Commands

		p, err := decodeConfig(
			schema.TestResourceDataRaw(t, Provisioner().(*schema.Provisioner).Schema, tc.Config),
		)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		err = p.linuxUploadRingKey(o, c)
		if err != nil {
			t.Fatalf("Test %q failed: %v", k, err)
		}
	}
}

func TestLinuxProvisioner_linuxStartbiomeService(t *testing.T) {
	cases := map[string]struct {
		Config   map[string]interface{}
		Commands map[string]bool
		Uploads  map[string]string
	}{
		"Start biome service with sudo": {
			Config: map[string]interface{}{
				"version":          "0.79.1",
				"auto_update":      false,
				"use_sudo":         true,
				"service_name":     "bio-sup",
				"peers":            []interface{}{"1.2.3.4"},
				"ring_key":         "test-ring",
				"ring_key_content": "dead-beef",
				"service": []interface{}{
					map[string]interface{}{
						"name":      "core/foo",
						"topology":  "standalone",
						"strategy":  "none",
						"channel":   "stable",
						"user_toml": "[config]\nlisten = 0.0.0.0:8080",
						"bind": []interface{}{
							map[string]interface{}{
								"alias":   "backend",
								"service": "bar",
								"group":   "default",
							},
						},
					},
					map[string]interface{}{
						"name":      "core/bar",
						"topology":  "standalone",
						"strategy":  "rolling",
						"channel":   "staging",
						"user_toml": "[config]\nlisten = 0.0.0.0:443",
					},
				},
			},

			Commands: map[string]bool{
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio pkg install core/foo  --channel stable'":                                                                        true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'mkdir -p /bio/user/foo/config'":                                                                                     true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'mv /tmp/user-a5b83ec1b302d109f41852ae17379f75c36dff9bc598aae76b6f7c9cd425fd76.toml /bio/user/foo/config/user.toml'": true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio svc load core/foo  --topology standalone --strategy none --channel stable --bind backend:bar.default'":          true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio pkg install core/bar  --channel staging'":                                                                       true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'mkdir -p /bio/user/bar/config'":                                                                                     true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'mv /tmp/user-6466ae3283ae1bd4737b00367bc676c6465b25682169ea5f7da222f3f078a5bf.toml /bio/user/bar/config/user.toml'": true,
				"env BIO_NONINTERACTIVE=true BIO_NOCOLORING=true sudo -E /bin/bash -c 'bio svc load core/bar  --topology standalone --strategy rolling --channel staging'":                                 true,
			},

			Uploads: map[string]string{
				"/tmp/user-a5b83ec1b302d109f41852ae17379f75c36dff9bc598aae76b6f7c9cd425fd76.toml": "[config]\nlisten = 0.0.0.0:8080",
				"/tmp/user-6466ae3283ae1bd4737b00367bc676c6465b25682169ea5f7da222f3f078a5bf.toml": "[config]\nlisten = 0.0.0.0:443",
			},
		},
	}

	o := new(terraform.MockUIOutput)
	c := new(communicator.MockCommunicator)

	for k, tc := range cases {
		c.Commands = tc.Commands
		c.Uploads = tc.Uploads

		p, err := decodeConfig(
			schema.TestResourceDataRaw(t, Provisioner().(*schema.Provisioner).Schema, tc.Config),
		)
		if err != nil {
			t.Fatalf("Error: %v", err)
		}

		var errs []error
		for _, s := range p.Services {
			err = p.linuxStartbiomeService(o, c, s)
			if err != nil {
				errs = append(errs, err)
			}
		}

		if len(errs) > 0 {
			for _, e := range errs {
				t.Logf("Test %q failed: %v", k, e)
				t.Fail()
			}
		}
	}
}
