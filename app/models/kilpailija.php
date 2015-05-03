<?php

class Kilpailija extends BaseModel{

    // Attribuutit
	public $id, $nimi, $kansallisuus, $sukupuoli, $syntynyt;

  	// Konstruktori
	public function __construct($attributes){
		parent::__construct($attributes);
		$this->validators = array('validate_nimi', 'validate_kansallisuus', 'validate_sukupuoli');
	}

    /**
     * Validoi olion nimi-attribuutin
     * @return array Mahdolliset virheilmoitukset
     */
	public function validate_nimi() {
		$errors = array();
		if ($this->nimi == '' || $this->nimi == null) {
			$errors[] = "Nimi ei saa olla tyhjä!";
		}
		if(strlen($this->nimi) < 3){
			$errors[] = 'Nimen pituuden tulee olla vähintään kolme merkkiä!';
		}
		return $errors;
	}

    /**
     * Validoi olion kansallisuus-attribuutin
     * @return array Mahdolliset virheilmoitukset
     */
	public function validate_kansallisuus() {
		$errors = array();
		if ($this->kansallisuus == '' || $this->kansallisuus == null) {
			$errors[] = "Kansallisuus ei saa olla tyhjä!";
		}
		if(strlen($this->kansallisuus) > 1){
			$errors[] = 'Kansallisuuden tulee olla vähintään kaksi merkkiä!';
		}
		return $errors;
	}

    /**
     * Validoi olion sukupuoli-attribuutin
     * @return array Mahdolliset virheilmoitukset
     */
	public function validate_sukupuoli() {
		$errors = array();
		if (!($this->sukupuoli == 'M' || $this->sukupuoli == 'N')) {
			$errors[] = "Sukupuoli määritelty väärin";
		}
		return $errors;
	}

    /**
     * Palauttaa kaikki kilpailijat tietokannasta
     * @return array Kaikki tietokannassa olevat kilpailijat
     */
	public static function all(){
		$query = DB::connection()->prepare('SELECT * FROM Kilpailija');
		$query->execute();
		$rows = $query->fetchAll();
		$kilpailijat = array();

		foreach($rows as $row){
			$kilpailijat[] = new Kilpailija(array(
				'id' => $row['id'],
				'nimi' => $row['nimi'],
				'kansallisuus' => $row['kansallisuus'],
				'sukupuoli' => $row['sukupuoli'],
				'syntynyt' => $row['syntynyt']
				));
		}

		return $kilpailijat;
	}

    /**
     * Palauttaa kilpailijan, jolla on pyydetty id
     * @param $id Kilpailijan id
     * @return Kilpailija|null
     */
	public static function find($id){
		$query = DB::connection()->prepare('SELECT * FROM Kilpalija WHERE id = :id LIMIT 1');
		$query->execute(array('id' => $id));
		$row = $query->fetch();

		if($row){
			$kilpailija = new Kilpailija(array(
				'id' => $row['id'],
				'nimi' => $row['nimi'],
				'kansallisuus' => $row['kansallisuus'],
				'sukupuoli' => $row['sukupuoli'],
				'syntynyt' => $row['syntynyt']
				));

			return $kilpailija;
		}

		return null;
	}

    /**
     * Tallentaa olion tietokantaan
     */
	public function save(){
		$query = DB::connection()->prepare('INSERT INTO Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES (:nimi, :kansallisuus, :sukupuoli, :syntynyt) RETURNING id');
		$query->execute(array('nimi' => $this->nimi, 'kansallisuus' => $this->kansallisuus, 'sukupuoli' => $this->sukupuoli, 'syntynyt' => $this->syntynyt));
		$row = $query->fetch();
		$this->id = $row['id'];
	}

    /**
     * Päivittää olemassaolevan olion tietokantaan
     */
	public function update(){
		$query = DB::connection()->prepare('UPDATE Kilpailija SET nimi = :nimi, kansallisuus = :kansallisuus, sukupuoli = :sukupuoli, syntynyt = :syntynyt WHERE id = :id');
		$query->execute(array('nimi' => $this->nimi, 'kansallisuus' => $this->kansallisuus, 'sukupuoli' => $this->sukupuoli, 'syntynyt' => $this->syntynyt, 'id' => $this->id));
	}

    /**
     * Poistaa kilpailijan tietokannasta
     */
	public function delete(){
		$query = DB::connection()->prepare('DELETE FROM Kilpailija WHERE id = :id');
		$query->execute();
		// todo: tee tässä jotain olion kadottamiseen
	}

    /**
     * Poistaa tietokannasta kilpailijan jolla on määritelty id
     * @param $id Kilpailijan id
     */
    public static function deleteId($id) {
        $query = DB::connection()->prepare('DELETE FROM Kilpailija WHERE id = :id');
        $query->execute(array(
            'id' => $id
        ));
    }

}
