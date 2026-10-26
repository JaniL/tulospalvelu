package geofilter_test

import (
	"net/netip"
	"testing"

	"github.com/JaniL/tulospalvelu/internal/geofilter"
)

func newFilter(t *testing.T) *geofilter.Filter {
	t.Helper()
	f, err := geofilter.New()
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(f.Prefixes()) == 0 {
		t.Fatal("lista on tyhjä")
	}
	return f
}

func TestListContainsOnlyValidCIDRs(t *testing.T) {
	newFilter(t)
}

func TestFinnishPrefixAddressIsAccepted(t *testing.T) {
	f := newFilter(t)
	for _, p := range f.Prefixes() {
		if !f.IsFinnish(p.Addr()) {
			t.Fatalf("listan oma osoite %v ei läpäissyt tarkistusta", p.Addr())
		}
	}
}

func TestNonFinnishAddressIsRejected(t *testing.T) {
	f := newFilter(t)
	for _, s := range []string{
		"8.8.8.8",     // Google, USA
		"1.1.1.1",     // Cloudflare, AU
		"2606:4700::", // Cloudflare IPv6, USA
	} {
		addr := netip.MustParseAddr(s)
		if f.IsFinnish(addr) {
			t.Errorf("ei-suomalainen osoite %s hyväksyttiin", s)
		}
	}
}

func TestIPv4MappedNonFinnishAddressIsRejected(t *testing.T) {
	f := newFilter(t)
	addr := netip.MustParseAddr("::ffff:8.8.8.8")
	if f.IsFinnish(addr) {
		t.Error("IPv6-muotoinen ei-suomalainen osoite hyväksyttiin")
	}
}

func TestInvalidAddressIsRejected(t *testing.T) {
	f := newFilter(t)
	if f.IsFinnish(netip.Addr{}) {
		t.Error("validoimaton osoite hyväksyttiin")
	}
}
