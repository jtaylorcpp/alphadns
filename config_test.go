package alphadns

import (
	"testing"
)

const config string = `
serverAddress: 127.0.0.1
serverPort: 8053
domains:
- name: test1.example.com
  addresses:
  - 1.1.1.1
  - 1.1.1.2
  - 1.1.1.3
- name: test2.example.com
  addresses:
  - 2.2.2.1
  - 2.2.2.2
  - 2.2.2.3
`

func TestParseConfig(t *testing.T) {
	dnsConfig := parseConfigFile(config)
	t.Log("parsed config: ", dnsConfig)

	dnsServer := DNSServerFromConfig(dnsConfig)
	t.Log("new server: ", dnsServer)
}
