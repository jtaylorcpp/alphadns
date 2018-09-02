package alphadns

import (
	"net"
	"reflect"
	"testing"
)

func TestSliceReverse(t *testing.T) {
	s := []string{"hello", "world"}
	rs := reverse(s)
	if rs[0] != "world" || rs[1] != "hello" {
		t.Error("reverse slice broken: ", s, rs)
	}
}

func TestDNSNodes(t *testing.T) {
	node := DNSNode{
		Name:       "test",
		SubDomains: make(map[Domain]*DNSNode),
		Addresses:  []*net.IP{},
	}

	t.Log("node test: ", node)
}

func TestDNSGraph(t *testing.T) {
	graph := &DNSGraph{
		Roots: make(map[Domain]*DNSNode),
	}
	t.Log("graph: ", graph)

	graph.AddRecord("test.example.com", []string{"1.1.1.1", "2.2.2.2"})
	t.Log("graph: ", graph)
	testIPs, _ := graph.GetIPAddresses("test.example.com")
	t.Log("graph ips test.example.com: ", testIPs)
	if !reflect.DeepEqual([]string{"1.1.1.1", "2.2.2.2"},
		[]string{testIPs[0].String(), testIPs[1].String()}) {
		t.Error("graphs ips not called back correctly")
	}

	graph.AddRecord("*.example.com", []string{"3.3.3.3", "4.4.4.4"})
	t.Log("graph: ", graph)
	wildcard1Ips, _ := graph.GetIPAddresses("test1.example.com")
	if !reflect.DeepEqual([]string{"3.3.3.3", "4.4.4.4"},
		[]string{wildcard1Ips[0].String(), wildcard1Ips[1].String()}) {
		t.Error("graph wildcard ips not called back correctly")
	}

	_, ok := graph.GetIPAddresses("google.com")
	t.Log("domain 'google.com' should not exist: ", ok)
	if ok {
		t.Fatal("google.com should not return true")
	}
}
