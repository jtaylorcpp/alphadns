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
	DNSRecords    *DNSGraph
}

func (ds *DNSServer) AddRecord(domainString string, ips []string) {
	ds.DNSRecords.AddRecord(domainString, ips)
}

func dnsHandler(w dns.ResponseWriter, r *dns.Msg) {
	log.Println("dns request: ", r.String())
}

func (ds *DNSServer) Start() {
	dns.HandleFunc(".", dnsHandler)
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
