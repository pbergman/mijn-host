package mijnhost

import (
	"context"

	"github.com/libdns/libdns"
)

func (p *Provider) ListZones(ctx context.Context) ([]libdns.Zone, error) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	domains, err := p.getClient().GetDomains(ctx)

	if err != nil {
		return nil, err
	}

	var zones = make([]libdns.Zone, len(domains))

	for i, c := 0, len(domains); i < c; i++ {
		zones[i] = libdns.Zone{
			Name: fqdn(domains[i].Domain),
		}
	}

	return zones, nil
}
