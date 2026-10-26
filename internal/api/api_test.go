package api_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/JaniL/tulospalvelu"
	"github.com/JaniL/tulospalvelu/internal/api"
	"github.com/JaniL/tulospalvelu/internal/geofilter"
	"github.com/JaniL/tulospalvelu/internal/store"
)

type fakeStore struct {
	store.Store

	kilpailijat []store.Kilpailija
	kilpailut   []store.Kilpailu
	lahtolista  []store.LahtolistaSijoitus

	lisattyKilpailija   *store.Kilpailija
	lisattyKisa         *store.Kilpailu
	poistettuKisa       int
	poistettuKilpailija int
}

func (f *fakeStore) Kilpailijat(ctx context.Context) ([]store.Kilpailija, error) {
	return f.kilpailijat, nil
}

func (f *fakeStore) Kilpailut(ctx context.Context) ([]store.Kilpailu, error) {
	return f.kilpailut, nil
}

func (f *fakeStore) Kisa(ctx context.Context, id int) (*store.Kilpailu, error) {
	for _, k := range f.kilpailut {
		if k.ID == id {
			kopio := k
			return &kopio, nil
		}
	}
	return nil, nil
}

func (f *fakeStore) Lahtolista(ctx context.Context, kisaID int) ([]store.LahtolistaSijoitus, error) {
	return f.lahtolista, nil
}

func (f *fakeStore) LisaaKilpailija(ctx context.Context, nimi, kansallisuus, sukupuoli, syntynyt string) (int, error) {
	f.lisattyKilpailija = &store.Kilpailija{Nimi: nimi, Kansallisuus: kansallisuus, Sukupuoli: sukupuoli, Syntynyt: syntynyt}
	return 7, nil
}

func (f *fakeStore) PoistaKilpailija(ctx context.Context, id int) error {
	f.poistettuKilpailija = id
	return nil
}

func (f *fakeStore) LisaaKisa(ctx context.Context, nimi, alkaa string, valiaikapisteita int) (int, error) {
	f.lisattyKisa = &store.Kilpailu{Nimi: nimi, Alkaa: alkaa, Valiaikapisteita: valiaikapisteita}
	return 5, nil
}

func (f *fakeStore) PoistaKisa(ctx context.Context, id int) error {
	f.poistettuKisa = id
	return nil
}

func (f *fakeStore) Ping(ctx context.Context) error {
	return nil
}

func newTestServer(t *testing.T, s store.Store, opts api.Options) *httptest.Server {
	t.Helper()
	if opts.GeoRestrict && opts.Filter == nil {
		f, err := geofilter.New()
		if err != nil {
			t.Fatalf("geofilter.New: %v", err)
		}
		opts.Filter = f
	}
	srv := httptest.NewServer(api.New(s, tulospalvelu.StaticFS, opts).Handler())
	t.Cleanup(srv.Close)
	return srv
}

// fiAddr palauttaa deterministisesti ensimmäisen suomalaisen osoitteen
// upotetusta listasta testifikstuureja varten.
func fiAddr(t *testing.T) string {
	t.Helper()
	f, err := geofilter.New()
	if err != nil {
		t.Fatalf("geofilter.New: %v", err)
	}
	return f.Prefixes()[0].Addr().String()
}

func TestKilpailijatListJSONFormat(t *testing.T) {
	s := &fakeStore{kilpailijat: []store.Kilpailija{
		{ID: 1, Nimi: "Matti M", Kansallisuus: "FI", Sukupuoli: "M", Syntynyt: "1992-11-11"},
	}}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Get(srv.URL + "/api/kilpailijat/list")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	want := `[{"id":1,"nimi":"Matti M","kansallisuus":"FI","sukupuoli":"M","syntynyt":"1992-11-11"}]`
	if got := strings.TrimSpace(string(body)); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("content-type = %s", ct)
	}
}

func TestKilpailijatListEmptyReturnsArray(t *testing.T) {
	s := &fakeStore{kilpailijat: []store.Kilpailija{}}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Get(srv.URL + "/api/kilpailijat/list")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if got := strings.TrimSpace(string(body)); got != "[]" {
		t.Errorf("got %s, want []", got)
	}
}

func TestKilpailutListIncludesValiaikapisteet(t *testing.T) {
	s := &fakeStore{kilpailut: []store.Kilpailu{{
		ID: 1, Nimi: "Kumpulan hiihtokisat", Alkaa: "2014-05-05", Valiaikapisteita: 2,
		Valiaikapisteet: []store.KisaAika{{ID: 1, KisaID: 1, KilpailijaID: 1, Valiaikapiste: 1, Aika: 5000}},
	}}}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Get(srv.URL + "/api/kilpailut/list")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	want := `[{"id":1,"nimi":"Kumpulan hiihtokisat","alkaa":"2014-05-05","valiaikapisteita":2,"valiaikapisteet":[{"id":1,"kisaId":1,"kilpailijaId":1,"valiaikapiste":1,"aika":5000}]}]`
	if got := strings.TrimSpace(string(body)); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestKilpailuByID(t *testing.T) {
	kilpailija := store.Kilpailija{ID: 1, Nimi: "Matti M", Kansallisuus: "FI", Sukupuoli: "M", Syntynyt: "1992-11-11"}
	s := &fakeStore{
		kilpailut:  []store.Kilpailu{{ID: 1, Nimi: "Kumpulan hiihtokisat", Alkaa: "2014-05-05", Valiaikapisteita: 2, Valiaikapisteet: []store.KisaAika{}}},
		lahtolista: []store.LahtolistaSijoitus{{ID: 1, KisaID: 1, Kilpailija: kilpailija, KilpailijaID: 1, Sijoitus: 1}},
	}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Get(srv.URL + "/api/kilpailut/id/1")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	want := `{"kilpailu":{"id":1,"nimi":"Kumpulan hiihtokisat","alkaa":"2014-05-05","valiaikapisteita":2,"valiaikapisteet":[]},"lahtolista":[{"id":1,"kisaId":1,"kilpailija":{"id":1,"nimi":"Matti M","kansallisuus":"FI","sukupuoli":"M","syntynyt":"1992-11-11"},"kilpailijaId":1,"sijoitus":1}]}`
	if got := strings.TrimSpace(string(body)); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestKilpailuByUnknownIDReturnsEmptyObject(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Get(srv.URL + "/api/kilpailut/id/999")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if got := strings.TrimSpace(string(body)); got != "{}" {
		t.Errorf("got %s, want {}", got)
	}
}

func TestLisaaKilpailijaReturnsIDAsText(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Post(srv.URL+"/api/kilpailijat/lisaa", "application/json",
		strings.NewReader(`{"nimi":"Kalle K","kansallisuus":"FI","syntynyt":"1990-01-01"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if got := strings.TrimSpace(string(body)); got != "7" {
		t.Errorf("got %q, want 7", got)
	}
	if s.lisattyKilpailija == nil {
		t.Fatal("kilpailijaa ei tallennettu")
	}
	if s.lisattyKilpailija.Sukupuoli != "M" {
		t.Errorf("sukupuoli = %q, want M", s.lisattyKilpailija.Sukupuoli)
	}
	if s.lisattyKilpailija.Nimi != "Kalle K" || s.lisattyKilpailija.Syntynyt != "1990-01-01" {
		t.Errorf("tallennettu kilpailija = %+v", s.lisattyKilpailija)
	}
}

func TestLisaaKisaReturnsIDAsText(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Post(srv.URL+"/api/kilpailut/lisaa", "application/json",
		strings.NewReader(`{"nimi":"Uusi kisa","alkaa":"2026-01-01","valiaikapisteita":3}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if got := strings.TrimSpace(string(body)); got != "5" {
		t.Errorf("got %q, want 5", got)
	}
	if s.lisattyKisa == nil || s.lisattyKisa.Valiaikapisteita != 3 {
		t.Errorf("tallennettu kisa = %+v", s.lisattyKisa)
	}
}

func TestPoistaKilpailijaReturnsOK(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Post(srv.URL+"/api/kilpailijat/poista", "application/json",
		strings.NewReader(`{"id":3}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if got := strings.TrimSpace(string(body)); got != "OK" {
		t.Errorf("got %q, want OK", got)
	}
}

func TestPoistaKisaCallsStoreWithCorrectID(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Post(srv.URL+"/api/kilpailut/poista", "application/json",
		strings.NewReader(`{"id":9}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	if s.poistettuKisa != 9 {
		t.Errorf("poistettuKisa = %d, want 9", s.poistettuKisa)
	}
}

func TestDatabaseConnectionReturnsOK(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Get(srv.URL + "/tietokantayhteys")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if got := strings.TrimSpace(string(body)); got != "OK" {
		t.Errorf("got %q, want OK", got)
	}
}

func TestIndexServesAngularApp(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Get(srv.URL + "/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if !strings.Contains(string(body), "tulospalveluApp") {
		t.Error("index.html ei sisällä angular-sovellusta")
	}
	if !strings.Contains(string(body), `ng-view`) {
		t.Error("index.html ei sisällä ng-view:tä")
	}
}

func TestStaticFilesAreServed(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	for _, path := range []string{
		"/vendor/jquery.min.js",
		"/vendor/bootstrap/css/bootstrap.min.css",
		"/vendor/bootstrap/fonts/glyphicons-halflings-regular.woff",
		"/assets/js/app.js",
		"/assets/js/controllers.js",
		"/assets/css/site.css",
		"/assets/partials/kilpailut.htm",
	} {
		resp, err := http.Get(srv.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			t.Errorf("%s: status = %d", path, resp.StatusCode)
		}
		if len(body) == 0 {
			t.Errorf("%s: tyhjä sisältö", path)
		}
	}
}

func TestGeoRestrictionDisabledAllowsAllPosts(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{})

	resp, err := http.Post(srv.URL+"/api/kilpailijat/poista", "application/json",
		strings.NewReader(`{"id":1}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestGeoRestrictionBlocksNonFinnishIP(t *testing.T) {
	s := &fakeStore{}
	// TrustXFF: asiakkaan IP luetaan X-Forwarded-For-otsakkeesta.
	srv := newTestServer(t, s, api.Options{GeoRestrict: true, TrustXFF: true})

	req, _ := http.NewRequest("POST", srv.URL+"/api/kilpailijat/poista",
		strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403", resp.StatusCode)
	}
	if !strings.Contains(string(body), "Suomen IP") {
		t.Errorf("virheilmoitus = %q", body)
	}
	if s.poistettuKilpailija != 0 {
		t.Error("ei-suomalaisen pyyntö ei saanut päätyä storeen")
	}
}

func TestGeoRestrictionAllowsFinnishIP(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{GeoRestrict: true, TrustXFF: true})

	req, _ := http.NewRequest("POST", srv.URL+"/api/kilpailijat/poista",
		strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", fiAddr(t))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
	if got := strings.TrimSpace(string(body)); got != "OK" {
		t.Errorf("body = %q, want OK", got)
	}
}

func TestGeoRestrictionUsesRightmostXFFValue(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{GeoRestrict: true, TrustXFF: true})

	// Vasemman arvon asiakas voi väärentää; oikeampi (edge-palvelimen
	// lisäämä) on ratkaiseva.
	req, _ := http.NewRequest("POST", srv.URL+"/api/kilpailijat/poista",
		strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "8.8.8.8, "+fiAddr(t))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (oikein XFF-arvo on suomalainen)", resp.StatusCode)
	}
}

func TestGeoRestrictionDoesNotApplyToGetRequests(t *testing.T) {
	s := &fakeStore{kilpailijat: []store.Kilpailija{}}
	srv := newTestServer(t, s, api.Options{GeoRestrict: true, TrustXFF: true})

	req, _ := http.NewRequest("GET", srv.URL+"/api/kilpailijat/list", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200 (GET pysyy julkisena)", resp.StatusCode)
	}
}

func TestGeoRestrictionWithoutXFFTrustUsesRemoteAddr(t *testing.T) {
	s := &fakeStore{}
	// TrustXFF: false → XFF:ää ei luoteta, IP otetaan RemoteAddrista
	// (httptest-palvelimessa 127.0.0.1, ei suomalainen).
	srv := newTestServer(t, s, api.Options{GeoRestrict: true})

	req, _ := http.NewRequest("POST", srv.URL+"/api/kilpailijat/poista",
		strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", fiAddr(t))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403 (RemoteAddr 127.0.0.1 ei ole suomalainen)", resp.StatusCode)
	}
}

func TestGeoRestrictionRejectsMalformedIP(t *testing.T) {
	s := &fakeStore{}
	srv := newTestServer(t, s, api.Options{GeoRestrict: true, TrustXFF: true})

	req, _ := http.NewRequest("POST", srv.URL+"/api/kilpailijat/poista",
		strings.NewReader(`{"id":1}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Forwarded-For", "ei-ole-ip-osoite")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("status = %d, want 403", resp.StatusCode)
	}
}
