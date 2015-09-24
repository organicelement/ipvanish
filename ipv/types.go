package ipv

type IPVServer struct {
	Geometry struct {
		Coordinates []float64 `json:"coordinates"`
		Type        string    `json:"type"`
	} `json:"geometry"`
	Properties struct {
		Capacity             int    `json:"capacity"`
		City                 string `json:"city"`
		Continent            string `json:"continent"`
		ContinentCode        string `json:"continentCode"`
		Country              string `json:"country"`
		CountryCode          string `json:"countryCode"`
		Hostname             string `json:"hostname"`
		IP                   string `json:"ip"`
		Latitude             string `json:"latitude"`
		Longitude            string `json:"longitude"`
		Marker_cluster_small string `json:"marker-cluster-small"`
		Marker_color         string `json:"marker-color"`
		Online               bool   `json:"online"`
		Region               string `json:"region"`
		RegionAbbr           string `json:"regionAbbr"`
		RegionCode           string `json:"regionCode"`
		Title                string `json:"title"`
		Visible              bool   `json:"visible"`
	} `json:"properties"`
	Type string `json:"type"`
	Distance float64
	Latency float64
}
