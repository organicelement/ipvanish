package ipv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func GetServers() []IPVServer {
	var servers []IPVServer

	res, err := http.Get("https://www.ipvanish.com/api/servers.geojson")

	if err != nil {
		return nil
	}

	if err = json.NewDecoder(res.Body).Decode(&servers); err != nil {
		return nil
	}

	return servers
}

func GetLocation() (GeoIP, error) {
	var geoip GeoIP

	req := fmt.Sprintf("https://ipdata.info/lookup/%s", GetPublicIP())
	res, err := http.Get(req)

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return geoip, err
	}

	if err != nil {
		log.Printf("Error: %s", err.Error())
		return geoip, err
	}

	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(&geoip); err != nil {
		log.Printf("There was an issue [arsing the json: %d", err.Error())
		return geoip, err
	}

	return geoip, nil
}

func GetPublicIP() string {
	res, err := http.Get("http://checkip.amazonaws.com/")

	if err != nil {
		return ""
	}

	defer res.Body.Close()
	ipaddr, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(ipaddr))
}
