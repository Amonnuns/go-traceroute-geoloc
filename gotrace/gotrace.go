package gotrace

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
	"github.com/jpiontek/go-ip-api"
	"golang.org/x/net/icmp"
)

var (
	localhost string = "0.0.0.0"
	ICMP4Proto int = 1
	ICMP6Proto int = 58
)


type GoTracer struct {
	firstTTL int
	maxTTL int
	dstIP *net.IPAddr 
}

type ICMPAnswer struct {
	Peer net.Addr
	bufferSize int
	message *icmp.Message
} 

func createICMPAnswer(bufferSize int, 
	sourceAddr net.Addr,
	readBf []byte) *ICMPAnswer{

	icmpAnswer := new(ICMPAnswer)
	icmpAnswer.Peer = sourceAddr
	icmpAnswer.bufferSize = bufferSize

	rm, err := icmp.ParseMessage(ICMP4Proto, readBf[:bufferSize])
		if err != nil {
			log.Fatal(err)
		}

	icmpAnswer.message = rm

	return icmpAnswer
}



func sendICMPRequest(message []byte, conn *icmp.PacketConn,
	client goip.Client, dstIp *net.IPAddr ){

	wg :=  &sync.WaitGroup{}
	m :=  &sync.Mutex{}

	for ttl := 1; ttl<30; ttl++{ 

		conn.IPv4PacketConn().SetTTL(ttl)

		if _, err := conn.WriteTo(message, dstIp); err != nil{
			log.Fatal(err)
		}

		var durationSlice []time.Duration
		readBf := make([]byte, 1500)

		wg.Add(3)
		for i := 0; i<3; i++ {
			go func(readBf []byte, wg *sync.WaitGroup,
				 m *sync.Mutex) {

				defer wg.Done()
				start := time.Now()
				n, peer, err := conn.ReadFrom(readBf)
				if err != nil {
					println(err)
				}
				duration := time.Since(start)

				
				m.Lock()
				durationSlice = append(durationSlice,duration)
				m.Unlock()
			}(readBf, wg, m)
		}
		wg.Wait()

		

			ip := peer.String()
			city, country, isp := searchGeoloc(ip, client)

			fmt.Printf("%v %v  %v %v %v %v", ttl, ip, city, country, isp, duration)
			println()
		

		if ip == dstIp.String(){
			break;
		}
	}
}