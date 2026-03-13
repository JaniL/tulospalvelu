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

## Podman / Docker

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
