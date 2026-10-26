package store

const dateFormat = "2006-01-02"

type Kilpailija struct {
	ID           int    `json:"id"`
	Nimi         string `json:"nimi"`
	Kansallisuus string `json:"kansallisuus"`
	Sukupuoli    string `json:"sukupuoli"`
	Syntynyt     string `json:"syntynyt"`
}

type KisaAika struct {
	ID            int `json:"id"`
	KisaID        int `json:"kisaId"`
	KilpailijaID  int `json:"kilpailijaId"`
	Valiaikapiste int `json:"valiaikapiste"`
	Aika          int `json:"aika"`
}

type Kilpailu struct {
	ID               int        `json:"id"`
	Nimi             string     `json:"nimi"`
	Alkaa            string     `json:"alkaa"`
	Valiaikapisteita int        `json:"valiaikapisteita"`
	Valiaikapisteet  []KisaAika `json:"valiaikapisteet"`
}

type LahtolistaSijoitus struct {
	ID           int        `json:"id"`
	KisaID       int        `json:"kisaId"`
	Kilpailija   Kilpailija `json:"kilpailija"`
	KilpailijaID int        `json:"kilpailijaId"`
	Sijoitus     int        `json:"sijoitus"`
}
