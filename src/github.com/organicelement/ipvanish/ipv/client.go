package ipv
import (
"net/http"
"encoding/json"
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
	res, err := http.Get("http://freegeoip.net/json/")

	if err != nil {
		return geoip, err
	}

	if err = json.NewDecoder(res.Body).Decode(&geoip); err != nil {
		return geoip, err
	}

	return geoip, nil
}
