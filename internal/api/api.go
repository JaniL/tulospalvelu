package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"net/netip"
	"strconv"
	"strings"

	"github.com/JaniL/tulospalvelu/internal/geofilter"
	"github.com/JaniL/tulospalvelu/internal/store"
)

// Options ohjaa API:n käyttäytymistä ympäristössä, jossa se ajetaan.
type Options struct {
	// GeoRestrict estää kirjoitusoperaatiot (POST) muista kuin Suomen
	// IP-osoitteista.
	GeoRestrict bool
	// TrustXFF hakee asiakkaan IP-osoitteen X-Forwarded-For-otsakkeen
	// oikeimmasta arvosta (scaleway-edgen takana). Muuten käytetään
	// RemoteAddria (esim. lokaali kehitys).
	TrustXFF bool
	// Filter on Suomen IP-lista. Pakollinen, jos GeoRestrict on päällä.
	Filter *geofilter.Filter
}

type API struct {
	store   store.Store
	static  fs.FS
	options Options
}

func New(s store.Store, static fs.FS, opts Options) *API {
	if opts.GeoRestrict && opts.Filter == nil {
		panic("api.New: GeoRestrict vaatii Filterin")
	}
	return &API{store: s, static: static, options: opts}
}

func (a *API) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/kilpailijat/list", a.kilpailijatList)
	mux.HandleFunc("POST /api/kilpailijat/lisaa", a.kilpailijaLisaa)
	mux.HandleFunc("POST /api/kilpailijat/poista", a.kilpailijaPoista)
	mux.HandleFunc("POST /api/kilpailijat/paivita", a.kilpailijaPaivita)

	mux.HandleFunc("GET /api/kilpailut/list", a.kilpailutList)
	mux.HandleFunc("GET /api/kilpailut/id/{kisaId}", a.kilpailuByID)
	mux.HandleFunc("POST /api/kilpailut/lisaaSijoitus", a.lahtosijoitusLisaa)
	mux.HandleFunc("POST /api/kilpailut/lisaaLahtosijoitus", a.lahtosijoitusLisaa)
	mux.HandleFunc("POST /api/kilpailut/poistaLahtosijoitus", a.lahtosijoitusPoista)
	mux.HandleFunc("POST /api/kilpailut/vaihdaSijoitukset", a.sijoituksetVaihda)
	mux.HandleFunc("POST /api/kilpailut/lisaaAika", a.aikaLisaa)
	mux.HandleFunc("POST /api/kilpailut/poistaAika", a.aikaPoista)
	mux.HandleFunc("POST /api/kilpailut/lisaa", a.kisaLisaa)
	mux.HandleFunc("POST /api/kilpailut/poista", a.kisaPoista)
	mux.HandleFunc("POST /api/kilpailut/paivita", a.kisaPaivita)

	mux.HandleFunc("GET /tietokantayhteys", a.dbTest)
	mux.HandleFunc("GET /hiekkalaatikko", a.sandbox)
	mux.HandleFunc("GET /{$}", a.index)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServerFS(mustSub(a.static, "assets"))))
	mux.Handle("GET /vendor/", http.StripPrefix("/vendor/", http.FileServerFS(mustSub(a.static, "vendor"))))

	return a.requireFinnishIP(mux)
}

// clientIP palauttaa pyynnön tekijän IP-osoitteen. XFF-luottamisen
// yhteydessä oikeinmpaa arvoa käytetään, koska jokainen välityspalvelin
// lisää todellisen lähdeosoitteen listan loppuun (vasemman arvon
// asiakas voi väärentää).
func clientIP(r *http.Request, trustXFF bool) string {
	if trustXFF {
		if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			parts := strings.Split(xff, ",")
			for i := len(parts) - 1; i >= 0; i-- {
				ip := strings.TrimSpace(parts[i])
				if ip != "" {
					return ip
				}
			}
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return strings.TrimSpace(r.RemoteAddr)
	}
	return host
}

// requireFinnishIP sallii kirjoitusoperaatiot (POST) vain Suomen
// IP-osoitteista, kun GeoRestrict on päällä. GET-pyynnöt ja staattiset
// tiedostot pysyvät julkisina.
func (a *API) requireFinnishIP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !a.options.GeoRestrict || r.Method != http.MethodPost {
			next.ServeHTTP(w, r)
			return
		}
		addr, err := netip.ParseAddr(clientIP(r, a.options.TrustXFF))
		if err != nil || !a.options.Filter.IsFinnish(addr) {
			log.Printf("kirjoitusoperaatio estetty: ip=%q url=%s", clientIP(r, a.options.TrustXFF), r.URL.Path)
			http.Error(w, "kirjoitusoperaatiot sallitaan vain Suomen IP-osoitteista", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func writeText(w http.ResponseWriter, s string) {
	fmt.Fprint(w, s)
}

func storeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
}

func decodeBody[T any](r *http.Request) (T, error) {
	var v T
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}

type kilpailijaLisaaReq struct {
	Nimi         string `json:"nimi"`
	Kansallisuus string `json:"kansallisuus"`
	Syntynyt     string `json:"syntynyt"`
}

type idReq struct {
	ID int `json:"id"`
}

type kisaLisaaReq struct {
	Nimi             string `json:"nimi"`
	Alkaa            string `json:"alkaa"`
	Valiaikapisteita int    `json:"valiaikapisteita"`
}

type lahtosijoitusLisaaReq struct {
	ID       int `json:"id"`
	Sijoitus int `json:"sijoitus"`
	KisaID   int `json:"kisaId"`
}

type lahtosijoitusPoistaReq struct {
	ID     int `json:"id"`
	KisaID int `json:"kisaId"`
}

type vaihdaSijoituksetReq struct {
	ID1    int `json:"id1"`
	ID2    int `json:"id2"`
	KisaID int `json:"kisaId"`
}

type aikaLisaaReq struct {
	KilpailuID   int `json:"kilpailuId"`
	KilpailijaID int `json:"kilpailijaId"`
	Valipiste    int `json:"valipiste"`
	Aika         int `json:"aika"`
}

type aikaPoistaReq struct {
	KilpailijaID int `json:"kilpailijaId"`
	KilpailuID   int `json:"kilpailuId"`
	Valipiste    int `json:"valipiste"`
}

func (a *API) kilpailijatList(w http.ResponseWriter, r *http.Request) {
	kilpailijat, err := a.store.Kilpailijat(r.Context())
	if err != nil {
		storeError(w, err)
		return
	}
	writeJSON(w, kilpailijat)
}

func (a *API) kilpailijaLisaa(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[kilpailijaLisaaReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	id, err := a.store.LisaaKilpailija(r.Context(), req.Nimi, req.Kansallisuus, "M", req.Syntynyt)
	if err != nil {
		storeError(w, err)
		return
	}
	writeText(w, strconv.Itoa(id))
}

func (a *API) kilpailijaPoista(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[idReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.PoistaKilpailija(r.Context(), req.ID); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) kilpailijaPaivita(w http.ResponseWriter, r *http.Request) {
	k, err := decodeBody[store.Kilpailija](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.PaivitaKilpailija(r.Context(), k); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) kilpailutList(w http.ResponseWriter, r *http.Request) {
	kilpailut, err := a.store.Kilpailut(r.Context())
	if err != nil {
		storeError(w, err)
		return
	}
	writeJSON(w, kilpailut)
}

func (a *API) kilpailuByID(w http.ResponseWriter, r *http.Request) {
	kisaID, err := strconv.Atoi(r.PathValue("kisaId"))
	if err != nil {
		http.Error(w, "virheellinen id", http.StatusBadRequest)
		return
	}
	kisa, err := a.store.Kisa(r.Context(), kisaID)
	if err != nil {
		storeError(w, err)
		return
	}
	if kisa == nil {
		writeText(w, "{}")
		return
	}
	lahtolista, err := a.store.Lahtolista(r.Context(), kisaID)
	if err != nil {
		storeError(w, err)
		return
	}
	writeJSON(w, struct {
		Kilpailu   *store.Kilpailu            `json:"kilpailu"`
		Lahtolista []store.LahtolistaSijoitus `json:"lahtolista"`
	}{Kilpailu: kisa, Lahtolista: lahtolista})
}

func (a *API) lahtosijoitusLisaa(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[lahtosijoitusLisaaReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.LisaaLahtoSijoitus(r.Context(), req.KisaID, req.ID, req.Sijoitus); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) lahtosijoitusPoista(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[lahtosijoitusPoistaReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.PoistaLahtoSijoitus(r.Context(), req.KisaID, req.ID); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) sijoituksetVaihda(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[vaihdaSijoituksetReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.VaihdaSijoitukset(r.Context(), req.KisaID, req.ID1, req.ID2); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) aikaLisaa(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[aikaLisaaReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.LisaaAika(r.Context(), req.KilpailuID, req.KilpailijaID, req.Valipiste, req.Aika); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) aikaPoista(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[aikaPoistaReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.PoistaAika(r.Context(), req.KilpailijaID, req.KilpailuID, req.Valipiste); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) kisaLisaa(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[kisaLisaaReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	id, err := a.store.LisaaKisa(r.Context(), req.Nimi, req.Alkaa, req.Valiaikapisteita)
	if err != nil {
		storeError(w, err)
		return
	}
	writeText(w, strconv.Itoa(id))
}

func (a *API) kisaPoista(w http.ResponseWriter, r *http.Request) {
	req, err := decodeBody[idReq](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.PoistaKisa(r.Context(), req.ID); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) kisaPaivita(w http.ResponseWriter, r *http.Request) {
	k, err := decodeBody[store.Kilpailu](r)
	if err != nil {
		storeError(w, err)
		return
	}
	if err := a.store.PaivitaKisa(r.Context(), k.ID, k.Nimi, k.Alkaa); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) dbTest(w http.ResponseWriter, r *http.Request) {
	if err := a.store.Ping(r.Context()); err != nil {
		storeError(w, err)
		return
	}
	writeText(w, "OK")
}

func (a *API) sandbox(w http.ResponseWriter, r *http.Request) {
	lahtolista, err := a.store.Lahtolista(r.Context(), 1)
	if err != nil {
		storeError(w, err)
		return
	}
	writeJSON(w, lahtolista)
}

func (a *API) index(w http.ResponseWriter, r *http.Request) {
	data, err := fs.ReadFile(a.static, "index.html")
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "index.html puuttuu", http.StatusInternalServerError)
			return
		}
		storeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(data)
}
