# IPVanish server listing with geo information

https://www.ipvanish.com/api/servers.geojson

# Generating go source from json

## IPVanish json
gojson -name IPVanish -input servers.sm.json -o ipv/types.go -pkg ipv

## Freegeoip json
curl http://freegeoip.net/json/|gojson -name GeoIP -o ipvanish/geoip.go -pkg ipv