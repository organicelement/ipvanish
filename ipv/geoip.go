package ipv

type GeoIP struct {
	City        string  `json:"city"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	IP          string  `json:"ip_address"`
	ISP         string  `json:"isp"`
	HostName    string  `json:"hostname"`
	GeoHash     string  `json:"geohash"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	TimeZone    string  `json:"time_zone"`
	ZipCode     string  `json:"zip_code"`
}
