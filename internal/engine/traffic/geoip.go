package traffic

import (
	_ "embed"
	"net"
	"os"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

//go:embed assets/GeoLite2-Country.mmdb
var embeddedGeoIPDB []byte

type GeoIPResolver struct {
	mu     sync.RWMutex
	reader *geoip2.Reader
}

func NewGeoIPResolver(customPath string) (*GeoIPResolver, error) {
	if customPath != "" {
		if _, statErr := os.Stat(customPath); statErr == nil {
			reader, err := geoip2.Open(customPath)
			if err != nil {
				return nil, err
			}
			return &GeoIPResolver{reader: reader}, nil
		}
	}

	if len(embeddedGeoIPDB) > 0 {
		reader, err := geoip2.FromBytes(embeddedGeoIPDB)
		if err != nil {
			return nil, err
		}
		return &GeoIPResolver{reader: reader}, nil
	}

	return &GeoIPResolver{}, nil
}

func (g *GeoIPResolver) Close() error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.reader != nil {
		err := g.reader.Close()
		g.reader = nil
		return err
	}
	return nil
}

func (g *GeoIPResolver) LookupCountry(ipStr string) string {
	if ipStr == "" {
		return ""
	}
	parsed := net.ParseIP(ipStr)
	if parsed == nil {
		return ""
	}
	if parsed.IsLoopback() || parsed.IsPrivate() || parsed.IsLinkLocalUnicast() || parsed.IsUnspecified() {
		return ""
	}

	g.mu.RLock()
	reader := g.reader
	g.mu.RUnlock()

	if reader == nil {
		return ""
	}

	record, err := reader.Country(parsed)
	if err != nil || record == nil {
		return ""
	}
	return record.Country.IsoCode
}
