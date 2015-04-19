INSERT INTO User (username, password) VALUES ('tsoha', '');

INSERT INTO Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES ('Matti M', 'FI', 'M', '1992-11-11');
INSERT INTO Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES ('Jouni A', 'FI', 'M', '1970-01-01');
INSERT INTO Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES ('Eeva V', 'FI', 'N', '1993-05-11');

INSERT INTO Kisa (nimi, alkaa, valiaikapisteita) VALUES ('Kumpulan hiihtokisat', '2014-05-05', 2);

INSERT INTO KisaLahtolista (kisaId, kilpailijaId, sijoitus) VALUES (1,1,1);

INSERT INTO KisaAika (kisaId,kilpailijaId,valiaikapiste,aika) VALUES (1,1,1,interval '5 seconds');
INSERT INTO KisaAika (kisaId,kilpailijaId,valiaikapiste,aika) VALUES (1,1,2,interval '10 seconds');