package geolocapi

import (
	"github.com/jpiontek/go-ip-api"
)

type Geoloc struct {
	client goip.Client
}

func New() *Geoloc{
	g := &Geoloc{
		client: goip.NewClient()
	}
}

func (g *Geoloc) searchGeoloc(ip string) (string, string, string){

	client := g.client
	
	city := ""
	country := ""
	isp := ""

	result, err := client.GetLocationForIp(ip)
	if err != nil{
		print()
	} else {
		city=result.City
		country=result.Country
		isp=result.Isp
	}

	return  city, country, isp
}
