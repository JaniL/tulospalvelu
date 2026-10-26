package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

var schemaStatements = []string{
	`CREATE TABLE IF NOT EXISTS kilpailija (
		id serial NOT NULL,
		nimi varchar(200) NOT NULL,
		kansallisuus varchar(2) NOT NULL,
		sukupuoli varchar(1) NOT NULL,
		syntynyt date NOT NULL,
		CONSTRAINT kilpailija_pk PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS kisa (
		id serial NOT NULL,
		nimi varchar(200) NOT NULL,
		alkaa date NOT NULL,
		valiaikapisteita int NOT NULL,
		CONSTRAINT kisa_pk PRIMARY KEY (id)
	)`,
	`CREATE TABLE IF NOT EXISTS kisaaika (
		id serial NOT NULL,
		kisaid int NOT NULL,
		kilpailijaid int NOT NULL,
		valiaikapiste int NOT NULL,
		aika int NOT NULL,
		CONSTRAINT kisaaika_pk PRIMARY KEY (id),
		CONSTRAINT kisaaika_kilpailija FOREIGN KEY (kilpailijaid) REFERENCES kilpailija (id),
		CONSTRAINT kisaaika_kisa FOREIGN KEY (kisaid) REFERENCES kisa (id)
	)`,
	`CREATE TABLE IF NOT EXISTS kisalahtolista (
		id serial NOT NULL,
		kisaid int NOT NULL,
		kilpailijaid int NOT NULL,
		sijoitus int NOT NULL,
		CONSTRAINT kisalahtolista_pk PRIMARY KEY (id),
		CONSTRAINT kisalahtolista_kilpailija FOREIGN KEY (kilpailijaid) REFERENCES kilpailija (id),
		CONSTRAINT kisalahtolista_kisa FOREIGN KEY (kisaid) REFERENCES kisa (id)
	)`,
	`CREATE TABLE IF NOT EXISTS "User" (
		id serial NOT NULL,
		username varchar(20) NOT NULL,
		password varchar(100) NOT NULL,
		active boolean NOT NULL,
		CONSTRAINT user_pk PRIMARY KEY (id)
	)`,
}

var seedStatements = []string{
	`INSERT INTO "User" (username, password, active) VALUES ('tsoha', '', true)`,
	`INSERT INTO kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES ('Matti M', 'FI', 'M', '1992-11-11')`,
	`INSERT INTO kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES ('Jouni A', 'FI', 'M', '1970-01-01')`,
	`INSERT INTO kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES ('Eeva V', 'FI', 'N', '1993-05-11')`,
	`INSERT INTO kisa (nimi, alkaa, valiaikapisteita) VALUES ('Kumpulan hiihtokisat', '2014-05-05', 2)`,
	`INSERT INTO kisalahtolista (kisaid, kilpailijaid, sijoitus) VALUES (1, 1, 1)`,
	`INSERT INTO kisaaika (kisaid, kilpailijaid, valiaikapiste, aika) VALUES (1, 1, 1, 5000)`,
	`INSERT INTO kisaaika (kisaid, kilpailijaid, valiaikapiste, aika) VALUES (1, 1, 2, 10000)`,
}

func Migrate(ctx context.Context, pool *pgxpool.Pool, seed bool) error {
	for _, stmt := range schemaStatements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("skeeman luonti: %w", err)
		}
	}
	if !seed {
		return nil
	}
	var count int
	if err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM kilpailija`).Scan(&count); err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	for _, stmt := range seedStatements {
		if _, err := pool.Exec(ctx, stmt); err != nil {
			return fmt.Errorf("testidatan lisays: %w", err)
		}
	}
	return nil
}
