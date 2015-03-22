-- foreign keys
ALTER TABLE KisaAika DROP CONSTRAINT KisaAika_Kilpailija;

ALTER TABLE KisaAika DROP CONSTRAINT KisaAika_Kisa;

ALTER TABLE KisaLahtolista DROP CONSTRAINT KisaLahtolista_Kilpailija;

ALTER TABLE KisaLahtolista DROP CONSTRAINT KisaLahtolista_Kisa;

-- tables
DROP TABLE Kilpailija;
DROP TABLE Kisa;
DROP TABLE KisaAika;
DROP TABLE KisaLahtolista;