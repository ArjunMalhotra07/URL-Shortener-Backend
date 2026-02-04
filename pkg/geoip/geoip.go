package geoip

import (
	"errors"
	"net"

	"github.com/oschwald/geoip2-golang"
)

var ErrInvalidIP = errors.New("invalid IP address")

type GeoResult struct {
	Country string
	City    string
	Region  string
}

type GeoIPService struct {
	db *geoip2.Reader
}

func NewGeoIPService(dbPath string) (*GeoIPService, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &GeoIPService{db: db}, nil
}

func (g *GeoIPService) Close() error {
	return g.db.Close()
}

func (g *GeoIPService) Lookup(ipStr string) (*GeoResult, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, ErrInvalidIP
	}

	record, err := g.db.City(ip)
	if err != nil {
		return nil, err
	}

	result := &GeoResult{
		Country: record.Country.Names["en"],
		City:    record.City.Names["en"],
	}

	if len(record.Subdivisions) > 0 {
		result.Region = record.Subdivisions[0].Names["en"]
	}

	return result, nil
}

// NullGeoIPService is a no-op implementation for when GeoIP is not configured
type NullGeoIPService struct{}

func NewNullGeoIPService() *NullGeoIPService {
	return &NullGeoIPService{}
}

func (g *NullGeoIPService) Lookup(ipStr string) (*GeoResult, error) {
	return &GeoResult{}, nil
}

func (g *NullGeoIPService) Close() error {
	return nil
}

// GeoIPLookup interface for dependency injection
type GeoIPLookup interface {
	Lookup(ipStr string) (*GeoResult, error)
	Close() error
}
