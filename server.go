package alphadns

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/miekg/dns"
)

// DNS Server is the stuct used to start the network service.
type DNSServer struct {
	ServerAddress string
	ServerPort    string
	TTL           uint32
	DNSRecords    *DNSGraph
}

func (ds *DNSServer) AddRecord(domainString string, ips []string) {
	ds.DNSRecords.AddRecord(domainString, ips)
}

func (ds *DNSServer) localLookup(r *dns.Msg) (*dns.Msg, bool) {
	m := new(dns.Msg)

	log.Println("dns questions: ",r.Question)
	if len(r.Question) > 1 {
		log.Println("more than 1 dns question: ", r.Question)
	}

	domain := r.Question[0].Name
	log.Println("looking up: ", domain)

	ips, ok := ds.DNSRecords.GetIPAddresses(domain)
	if !ok {
		return m, false
	}

	m.SetReply(r)
	m.Authoritative = true
	answers := make([]dns.RR, len(ips))

	for ipIdx, ip := range ips {
		record := new(dns.A)
		record.Hdr = dns.RR_Header{
			Name:   domain,
			Rrtype: dns.TypeA,
			Class:  dns.ClassINET,
			Ttl:    ds.TTL,
		}
		record.A = *ip
		answers[ipIdx] = record
	}

	m.Answer = answers
	return m, true
}

func (ds *DNSServer) dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	log.Println("dns request: ", r.String())
	msg, ok := ds.localLookup(r)
	if ok {
		log.Println("returning msg: ", msg.String())
		w.WriteMsg(msg)
	}
	log.Println("did not find entry in server")
}

func (ds *DNSServer) Start() {
	dns.HandleFunc(".", ds.dnsHandler)
	go func() {
		server := &dns.Server{
			Addr: ds.ServerAddress + ":" + ds.ServerPort,
			Net:  "udp",
		}
		log.Println("starting udp dns server")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalln("Failed to set udp listener: ", err.Error())
		}
	}()

	go func() {
		server := &dns.Server{
			Addr: ds.ServerAddress + ":" + ds.ServerPort,
			Net:  "tcp",
		}
		log.Println("starting tcp dns server")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatalln("Failed to set tcp listener: ", err.Error())
		}
	}()

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%d) received, stopping\n", s)
		}
	}

}
