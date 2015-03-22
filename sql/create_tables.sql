-- tables
-- Table: Kilpailija
CREATE TABLE Kilpailija (
    id int  NOT NULL,
    nimi varchar(200)  NOT NULL,
    kansallisuus varchar(2)  NOT NULL DEFAULT FI,
    sukupuoli varchar(1)  NOT NULL,
    syntynyt date  NOT NULL,
    CONSTRAINT Kilpailija_pk PRIMARY KEY (id)
);



-- Table: Kisa
CREATE TABLE Kisa (
    id int  NOT NULL,
    nimi varchar(200)  NOT NULL,
    alkaa date  NOT NULL,
    valiaikapisteita int  NOT NULL,
    CONSTRAINT Kisa_pk PRIMARY KEY (id)
);



-- Table: KisaAika
CREATE TABLE KisaAika (
    id int  NOT NULL,
    kisaId int  NOT NULL,
    kilpailijaId int  NOT NULL,
    valiaikapiste int  NOT NULL,
    aika interval  NOT NULL,
    CONSTRAINT KisaAika_pk PRIMARY KEY (id)
);



-- Table: KisaLahtolista
CREATE TABLE KisaLahtolista (
    id int  NOT NULL,
    kisaId int  NOT NULL,
    kilpailijaId int  NOT NULL,
    sijoitus int  NOT NULL,
    CONSTRAINT KisaLahtolista_pk PRIMARY KEY (id)
);







-- foreign keys
-- Reference:  KisaAika_Kilpailija (table: KisaAika)


ALTER TABLE KisaAika ADD CONSTRAINT KisaAika_Kilpailija
    FOREIGN KEY (kilpailijaId)
    REFERENCES Kilpailija (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  KisaAika_Kisa (table: KisaAika)


ALTER TABLE KisaAika ADD CONSTRAINT KisaAika_Kisa
    FOREIGN KEY (kisaId)
    REFERENCES Kisa (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  KisaLahtolista_Kilpailija (table: KisaLahtolista)


ALTER TABLE KisaLahtolista ADD CONSTRAINT KisaLahtolista_Kilpailija
    FOREIGN KEY (kilpailijaId)
    REFERENCES Kilpailija (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;

-- Reference:  KisaLahtolista_Kisa (table: KisaLahtolista)


ALTER TABLE KisaLahtolista ADD CONSTRAINT KisaLahtolista_Kisa
    FOREIGN KEY (kisaId)
    REFERENCES Kisa (id)
    NOT DEFERRABLE
    INITIALLY IMMEDIATE
;
