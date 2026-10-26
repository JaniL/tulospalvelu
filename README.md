# Hiihtokisojen tulospalvelu

**Huom.**

Kieli vaihtunut Rubysta PHP:hen.

http://jjluukko.users.cs.helsinki.fi/tulospalvelu/

Tietokantasovellus-harjoitustyö

Aiheen kuvaus: http://advancedkittenry.github.io/suunnittelu_ja_tyoymparisto/aiheet/Hiihtokisojen_tulospalvelu.html

Dokumentaatio: https://github.com/JaniL/tulospalvelu/raw/master/doc/dokumentaatio.pdf

Kisojen listaus / etusivu: http://jjluukko.users.cs.helsinki.fi/tulospalvelu/#/kilpailut

Kisatulokset: http://jjluukko.users.cs.helsinki.fi/tulospalvelu/#/kisa/1

Kilpailijoiden listaus: http://jjluukko.users.cs.helsinki.fi/tulospalvelu/#/kilpailijat

## Go-rewrite (branch `go-rewrite`)

Sovelluksesta on Go-versio, joka ei ole yhdistettynä masteriin. Se tarjoaa saman
JSON API:n samilla reiteillä ja sama AngularJS-käyttöliittymä on upotettu
binääriin (`go:embed`).

Lokaali ajo (Podman Compose):

```bash
podman compose -f compose.go.yaml up --build
```

- Sovellus: http://localhost:8080/
- Tietokannan skeema luodaan automaattisesti käynnistyksessä
- Testidata lisätään, jos tietokanta on tyhjä (`SEED_TEST_DATA=true`)

Lokaali kehitys ilman kontteja (vaatii paikallisen PostgreSQL:n):

```bash
make run
```

Testit ja tarkistukset:

```bash
make fmt vet test
```

Keskeiset ympäristömuuttujat:

| Muuttuja         | Merkitys                                                               |
|------------------|------------------------------------------------------------------------|
| `DATABASE_URL`   | PostgreSQL-yhteysmerkkijono (pakollinen)                                |
| `SEED_TEST_DATA` | `true` lisää testidatan tyhjään tietokantaan                           |
| `GEO_RESTRICT`   | `true` estää kirjoitusoperaatiot (POST) muilta kuin Suomen IP-osoitteilta |
| `TRUST_XFF`      | `true` lukee asiakkaan IP:n `X-Forwarded-For`-otsakkeesta (editin takana) |

Suomen IP-lista (`internal/geofilter/fi_networks.txt`) on rakennettu RIPE NCC:n
delegated-statistiikasta (CC BY 4.0); päivityskomento on paketin dokumentaatiossa.

### Scaleway

`terraform/` sisältää Terraform-konfiguraation (Serverless Container +
Serverless SQL Database + Container Registry). Deploy-ohje: [DEPLOY.md](DEPLOY.md).

## PHP-versio (Podman / Docker)

Projektin voi ajaa lokaalisti konteissa (Podman Compose tai Docker Compose).

Käynnistys:

```bash
podman compose up --build
```

Sovellus:

- http://localhost:8080/
- API-esimerkki: http://localhost:8080/api/kilpailut/list

Mitä käynnistyksessä tapahtuu:

- `db`-palvelu käynnistää PostgreSQL:n
- taulut luodaan automaattisesti
- testidata lisätään automaattisesti ensimmäisellä käynnistyksellä

Tietokannan alustus ajetaan vain kerran (kun volume on tyhjä). Jos haluat alustaa tietokannan uudelleen:

```bash
podman compose down -v
podman compose up --build
```
