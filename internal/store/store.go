package store

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Store interface {
	Kilpailijat(ctx context.Context) ([]Kilpailija, error)
	Kilpailut(ctx context.Context) ([]Kilpailu, error)
	Kisa(ctx context.Context, id int) (*Kilpailu, error)
	Lahtolista(ctx context.Context, kisaID int) ([]LahtolistaSijoitus, error)

	LisaaKilpailija(ctx context.Context, nimi, kansallisuus, sukupuoli, syntynyt string) (int, error)
	PaivitaKilpailija(ctx context.Context, k Kilpailija) error
	PoistaKilpailija(ctx context.Context, id int) error

	LisaaKisa(ctx context.Context, nimi, alkaa string, valiaikapisteita int) (int, error)
	PaivitaKisa(ctx context.Context, id int, nimi, alkaa string) error
	PoistaKisa(ctx context.Context, id int) error

	LisaaLahtoSijoitus(ctx context.Context, kisaID, kilpailijaID, sijoitus int) error
	PoistaLahtoSijoitus(ctx context.Context, kisaID, kilpailijaID int) error
	VaihdaSijoitukset(ctx context.Context, kisaID, sijoitus1, sijoitus2 int) error

	LisaaAika(ctx context.Context, kisaID, kilpailijaID, valiaikapiste, aika int) error
	PoistaAika(ctx context.Context, kilpailijaID, kisaID, valiaikapiste int) error

	Ping(ctx context.Context) error
}

type PGStore struct {
	pool *pgxpool.Pool
}

func NewPG(pool *pgxpool.Pool) *PGStore {
	return &PGStore{pool: pool}
}

func parseDate(s string) (time.Time, error) {
	return time.Parse(dateFormat, s)
}

func (s *PGStore) Kilpailijat(ctx context.Context) ([]Kilpailija, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, nimi, kansallisuus, sukupuoli, syntynyt FROM kilpailija ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kilpailijat := []Kilpailija{}
	for rows.Next() {
		var k Kilpailija
		var syntynyt time.Time
		if err := rows.Scan(&k.ID, &k.Nimi, &k.Kansallisuus, &k.Sukupuoli, &syntynyt); err != nil {
			return nil, err
		}
		k.Syntynyt = syntynyt.Format(dateFormat)
		kilpailijat = append(kilpailijat, k)
	}
	return kilpailijat, rows.Err()
}

func (s *PGStore) ajatByKisa(ctx context.Context, kisaID int) ([]KisaAika, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, kisaid, kilpailijaid, valiaikapiste, aika FROM kisaaika WHERE kisaid = $1 ORDER BY id`, kisaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ajat := []KisaAika{}
	for rows.Next() {
		var a KisaAika
		if err := rows.Scan(&a.ID, &a.KisaID, &a.KilpailijaID, &a.Valiaikapiste, &a.Aika); err != nil {
			return nil, err
		}
		ajat = append(ajat, a)
	}
	return ajat, rows.Err()
}

func (s *PGStore) Kilpailut(ctx context.Context) ([]Kilpailu, error) {
	rows, err := s.pool.Query(ctx, `SELECT id, nimi, alkaa, valiaikapisteita FROM kisa ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	kilpailut := []Kilpailu{}
	for rows.Next() {
		k, err := scanKilpailu(rows)
		if err != nil {
			return nil, err
		}
		kilpailut = append(kilpailut, *k)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range kilpailut {
		ajat, err := s.ajatByKisa(ctx, kilpailut[i].ID)
		if err != nil {
			return nil, err
		}
		kilpailut[i].Valiaikapisteet = ajat
	}
	return kilpailut, nil
}

func (s *PGStore) Kisa(ctx context.Context, id int) (*Kilpailu, error) {
	k, err := scanKilpailu(s.pool.QueryRow(ctx, `SELECT id, nimi, alkaa, valiaikapisteita FROM kisa WHERE id = $1`, id))
	if err != nil {
		return nil, err
	}
	ajat, err := s.ajatByKisa(ctx, id)
	if err != nil {
		return nil, err
	}
	k.Valiaikapisteet = ajat
	return k, nil
}

func scanKilpailu(row pgx.Row) (*Kilpailu, error) {
	var k Kilpailu
	var alkaa time.Time
	if err := row.Scan(&k.ID, &k.Nimi, &alkaa, &k.Valiaikapisteita); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	k.Alkaa = alkaa.Format(dateFormat)
	k.Valiaikapisteet = []KisaAika{}
	return &k, nil
}

func (s *PGStore) Lahtolista(ctx context.Context, kisaID int) ([]LahtolistaSijoitus, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT l.id, l.kisaid, l.kilpailijaid, l.sijoitus,
		       k.id, k.nimi, k.kansallisuus, k.sukupuoli, k.syntynyt
		FROM kisalahtolista l
		INNER JOIN kilpailija k ON k.id = l.kilpailijaid
		WHERE l.kisaid = $1
		ORDER BY l.sijoitus`, kisaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lahtolista := []LahtolistaSijoitus{}
	for rows.Next() {
		var l LahtolistaSijoitus
		var syntynyt time.Time
		if err := rows.Scan(&l.ID, &l.KisaID, &l.KilpailijaID, &l.Sijoitus,
			&l.Kilpailija.ID, &l.Kilpailija.Nimi, &l.Kilpailija.Kansallisuus, &l.Kilpailija.Sukupuoli, &syntynyt); err != nil {
			return nil, err
		}
		l.Kilpailija.Syntynyt = syntynyt.Format(dateFormat)
		lahtolista = append(lahtolista, l)
	}
	return lahtolista, rows.Err()
}

func (s *PGStore) LisaaKilpailija(ctx context.Context, nimi, kansallisuus, sukupuoli, syntynyt string) (int, error) {
	syntynytDate, err := parseDate(syntynyt)
	if err != nil {
		return 0, fmt.Errorf("syntynyt: %w", err)
	}
	var id int
	err = s.pool.QueryRow(ctx,
		`INSERT INTO kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES ($1, $2, $3, $4) RETURNING id`,
		nimi, kansallisuus, sukupuoli, syntynytDate).Scan(&id)
	return id, err
}

func (s *PGStore) PaivitaKilpailija(ctx context.Context, k Kilpailija) error {
	syntynyt, err := parseDate(k.Syntynyt)
	if err != nil {
		return fmt.Errorf("syntynyt: %w", err)
	}
	_, err = s.pool.Exec(ctx,
		`UPDATE kilpailija SET nimi = $1, kansallisuus = $2, sukupuoli = $3, syntynyt = $4 WHERE id = $5`,
		k.Nimi, k.Kansallisuus, k.Sukupuoli, syntynyt, k.ID)
	return err
}

func (s *PGStore) PoistaKilpailija(ctx context.Context, id int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM kisaaika WHERE kilpailijaid = $1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM kisalahtolista WHERE kilpailijaid = $1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM kilpailija WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) LisaaKisa(ctx context.Context, nimi, alkaa string, valiaikapisteita int) (int, error) {
	alkaaDate, err := parseDate(alkaa)
	if err != nil {
		return 0, fmt.Errorf("alkaa: %w", err)
	}
	var id int
	err = s.pool.QueryRow(ctx,
		`INSERT INTO kisa (nimi, alkaa, valiaikapisteita) VALUES ($1, $2, $3) RETURNING id`,
		nimi, alkaaDate, valiaikapisteita).Scan(&id)
	return id, err
}

func (s *PGStore) PaivitaKisa(ctx context.Context, id int, nimi, alkaa string) error {
	alkaaDate, err := parseDate(alkaa)
	if err != nil {
		return fmt.Errorf("alkaa: %w", err)
	}
	_, err = s.pool.Exec(ctx, `UPDATE kisa SET nimi = $1, alkaa = $2 WHERE id = $3`, nimi, alkaaDate, id)
	return err
}

func (s *PGStore) PoistaKisa(ctx context.Context, id int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM kisaaika WHERE kisaid = $1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM kisalahtolista WHERE kisaid = $1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx, `DELETE FROM kisa WHERE id = $1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) LisaaLahtoSijoitus(ctx context.Context, kisaID, kilpailijaID, sijoitus int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var exists bool
	if err := tx.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM kisalahtolista WHERE kisaid = $1 AND sijoitus = $2)`,
		kisaID, sijoitus).Scan(&exists); err != nil {
		return err
	}
	if exists {
		if _, err := tx.Exec(ctx,
			`UPDATE kisalahtolista SET sijoitus = sijoitus + 1 WHERE sijoitus >= $1 AND kisaid = $2`,
			sijoitus, kisaID); err != nil {
			return err
		}
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO kisalahtolista (kisaid, kilpailijaid, sijoitus) VALUES ($1, $2, $3)`,
		kisaID, kilpailijaID, sijoitus); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) PoistaLahtoSijoitus(ctx context.Context, kisaID, kilpailijaID int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var sijoitus int
	err = tx.QueryRow(ctx,
		`DELETE FROM kisalahtolista WHERE kisaid = $1 AND kilpailijaid = $2 RETURNING sijoitus`,
		kisaID, kilpailijaID).Scan(&sijoitus)
	if errors.Is(err, pgx.ErrNoRows) {
		return tx.Commit(ctx)
	}
	if err != nil {
		return err
	}

	if _, err := tx.Exec(ctx,
		`UPDATE kisalahtolista SET sijoitus = sijoitus - 1 WHERE sijoitus > $1 AND kisaid = $2`,
		sijoitus, kisaID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) VaihdaSijoitukset(ctx context.Context, kisaID, sijoitus1, sijoitus2 int) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`UPDATE kisalahtolista SET sijoitus = 0 WHERE sijoitus = $1 AND kisaid = $2`,
		sijoitus1, kisaID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`UPDATE kisalahtolista SET sijoitus = $1 WHERE sijoitus = $2 AND kisaid = $3`,
		sijoitus1, sijoitus2, kisaID); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`UPDATE kisalahtolista SET sijoitus = $1 WHERE sijoitus = 0 AND kisaid = $2`,
		sijoitus2, kisaID); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (s *PGStore) LisaaAika(ctx context.Context, kisaID, kilpailijaID, valiaikapiste, aika int) error {
	_, err := s.pool.Exec(ctx,
		`INSERT INTO kisaaika (kisaid, kilpailijaid, valiaikapiste, aika) VALUES ($1, $2, $3, $4)`,
		kisaID, kilpailijaID, valiaikapiste, aika)
	return err
}

func (s *PGStore) PoistaAika(ctx context.Context, kilpailijaID, kisaID, valiaikapiste int) error {
	_, err := s.pool.Exec(ctx,
		`DELETE FROM kisaaika WHERE kilpailijaid = $1 AND kisaid = $2 AND valiaikapiste = $3`,
		kilpailijaID, kisaID, valiaikapiste)
	return err
}

func (s *PGStore) Ping(ctx context.Context) error {
	return s.pool.Ping(ctx)
}
