package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"github.com/jpiontek/go-ip-api"
)

var (
	localhost string = "0.0.0.0"
)



func main(){

	rawIp := flag.String("addr", "127.0.0.1",
		"Address that we want do discover the route to")
	flag.Parse()

	//client for geolocip-api
	client := goip.NewClient()
	
	dstIp, err := net.ResolveIPAddr("ip4:icmp", *rawIp)
	if err != nil {
		log.Fatal(err)
	}

	connV4, err := icmp.ListenPacket("ip4:icmp", localhost)
	if err != nil{
		log.Fatal(err)
	}
	defer connV4.Close()

	message := icmp.Message{
					Type:  ipv4.ICMPTypeEcho,
					Code: 0,
					Body: &icmp.Echo{
						ID: os.Getpid() & 0xffff,
						Seq: 1,
						Data: []byte("HELLO-STRANGE"),
					},
				}
	
	msgBytes, err := message.Marshal(nil)
	if err != nil{
		log.Fatal(err)
	}
	
	for i := 2; i<30; i++{ 

		connV4.IPv4PacketConn().SetTTL(i)

		if _, err := connV4.WriteTo(msgBytes, dstIp); err != nil{
			log.Fatal(err)
		}

		readBf := make([]byte, 1500)
		n, peer, err := connV4.ReadFrom(readBf)
		if err != nil {
			println(err)
		}

		rm, err := icmp.ParseMessage(58, readBf[:n])
		if err != nil {
			log.Fatal(err)
		}

		switch rm.Type {
		case ipv4.ICMPTypeTimeExceeded:
			fmt.Printf("we have receive a timeExceeded from: %v", peer)
		default:
			peerIp := peer.String()
			result, err := client.GetLocationForIp(peerIp)
			city := ""
			country := ""

			if err != nil{
				print()
			}else {
				city=result.City
				country=result.Country
			}
			fmt.Printf("%v  %v %v  TTL:%v", peer, city, country, i)
			fmt.Println()
		}

		if peer.String() == dstIp.String(){
			break;
		}
	}

	
}
