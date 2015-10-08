# IPVanish server listing with geo information

https://www.ipvanish.com/api/servers.geojson

# Generating go source from json

## IPVanish json
```gojson -name IPVanish -input servers.sm.json -o ipv/types.go -pkg ipv```

## Freegeoip json
```curl http://freegeoip.net/json/|gojson -name GeoIP -o ipvanish/geoip.go -pkg ipv```

[![Analytics](https://ga-beacon.appspot.com/UA-68563453-1/ipvanish/readme.md?flat)](https://github.com/igrigorik/ga-beacon)

# License

Copyright 2015 Organic Element LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
