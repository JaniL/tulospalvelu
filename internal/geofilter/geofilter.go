// Package geofilter sisältää Suomen IP-osoiteavaruuden CIDR-listan ja
// siihen perustuvan tarkistuksen. Lista on rakennettu RIPE NCC:n
// delegated-statistiikkatiedostosta (CC BY 4.0) ja sisältää kaikki
// Suomeen rekisteröidyt IPv4- ja IPv6-lohkot.
//
// Listan päivitys:
//
//	curl -sS https://ftp.ripe.net/pub/stats/ripencc/delegated-ripencc-extended-latest \
//	  | awk -F'|' '$2=="FI" && $3=="ipv4" { p=32; n=$5; while (n>1) { n/=2; p-- }; print $4"/"p } $2=="FI" && $3=="ipv6" { print $4"/"$5 }' \
//	  > internal/geofilter/fi_networks.txt
package geofilter

import (
	"bufio"
	"embed"
	"fmt"
	"net/netip"
	"strings"
)

//go:embed fi_networks.txt
var fiNetworksFS embed.FS

// Filter tarkistaa, onko osoite osa Suomen IP-avaruutta.
type Filter struct {
	prefixes []netip.Prefix
}

// New jäsentää upotetun CIDR-listan. Lista on validoitu jäsentämisen
// yhteydessä, joten virhe palautetaan vain, jos upotettu tiedosto on
// vioittunut.
func New() (*Filter, error) {
	data, err := fiNetworksFS.ReadFile("fi_networks.txt")
	if err != nil {
		return nil, fmt.Errorf("fi_networks.txt: %w", err)
	}

	f := &Filter{}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		p, err := netip.ParsePrefix(line)
		if err != nil {
			return nil, fmt.Errorf("virheellinen CIDR-rivi %q: %w", line, err)
		}
		f.prefixes = append(f.prefixes, p)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(f.prefixes) == 0 {
		return nil, fmt.Errorf("fi_networks.txt on tyhjä")
	}
	return f, nil
}

// Prefixes palauttaa jäsenistetyn CIDR-listan (diagnostiikkaa ja testejä varten).
func (f *Filter) Prefixes() []netip.Prefix {
	return f.prefixes
}

// IsFinnish ilmoittaa, kuuluuko osoite johonkin Suomeen rekisteröityyn
// verkkolohkoon. Tuntematon osoite palauttaa false.
func (f *Filter) IsFinnish(addr netip.Addr) bool {
	if !addr.IsValid() {
		return false
	}
	for _, p := range f.prefixes {
		if p.Contains(addr) {
			return true
		}
	}
	return false
}
