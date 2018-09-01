package alphadns

import (
	"net"
	"reflect"
	"testing"
)

func TestDNSNodes(t *testing.T) {
	addr1 := net.ParseIP("0.0.0.0").To4()
	leaf1 := DNSLeaf{
		Name:      "test",
		Addresses: []*net.IP{&addr1},
	}

	addr2 := net.ParseIP("1.1.1.1").To4()
	leaf2 := DNSLeaf{
		Name:      "*",
		Addresses: []*net.IP{&addr2},
	}

	if !reflect.DeepEqual(*leaf1.GetIPAddresses([]Domain{"test"})[0],
		net.ParseIP("0.0.0.0").To4()) {
		t.Error("test should be 0.0.0.0")
	}

	node1 := DNSBranch{
		Name:       "example",
		SubDomains: map[Domain]DNSNode{leaf1.Name: &leaf1},
	}

	node1Address := node1.GetIPAddresses([]Domain{"example", "test"})
	if !reflect.DeepEqual(*node1Address[0], net.ParseIP("0.0.0.0").To4()) {
		t.Error("test.example should be 0.0.0.0")
	}

	node2Address := node1.GetIPAddresses([]Domain{"example", "test2"})
	if len(node2Address) != 0 {
		t.Error("test2.example should not exist")
	}

	node1.SubDomains[leaf2.Name] = &leaf2

	node3Address := node1.GetIPAddresses([]Domain{"example", "*"})
	t.Log(node3Address)
	if !reflect.DeepEqual(*node3Address[0], net.ParseIP("1.1.1.1").To4()) {
		t.Error("*.example should be 1.1.1.1")
	}

	node4Address := node1.GetIPAddresses([]Domain{"example", "notdefined"})
	t.Log(node4Address)
	if !reflect.DeepEqual(*node4Address[0], net.ParseIP("1.1.1.1").To4()) {
		t.Error("notdefined.example should be widcard *.example or 1.1.1.1")
	}
}
