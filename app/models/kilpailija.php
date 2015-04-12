<?php

class Kilpailija extends BaseModel{

    // Attribuutit
	public $id, $nimi, $kansallisuus, $sukupuoli, $syntynyt

  	// Konstruktori
	public function __construct($attributes){
		parent::__construct($attributes);
		$this->validators = array('validate_nimi', 'validate_kansallisuus', 'validate_sukupuoli');
	}

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

	public function validate_kansallisuus() {
		$errors = array();
		if ($this->kansallisuus == '' || $this->kansallisuus == null) {
			$errors[] = "Kansallisuus ei saa olla tyhjä!";
		}
		if(strlen($this->kansallisuus) < 3){
			$errors[] = 'Kansallisuuden tulee olla vähintään kolme merkkiä!';
		}
		return $errors;
	}

	public function validate_sukupuoli() {
		$errors = array();
		if (!($this->sukupuoli == 'M' || $this->sukupuoli == 'N')) {
			$errors[] = "Sukupuoli määritelty väärin";
		}
		return $errors;
	}

	public static function all(){
		$query = DB::connection()->prepare('SELECT * FROM Kilpailija');
		$query->execute();
		$rows = $query->fetchAll();
		$games = array();

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

	public function save(){
		$query = DB::connection()->prepare('INSERT INTO Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES (:nimi, :kansallisuus, :sukupuoli, :syntynyt) RETURNING id');
		$query->execute(array('nimi' => $this->nimi, 'kansallisuus' => $this->kansallisuus, 'sukupuoli' => $this->sukupuoli, 'syntynyt' => $this->syntynyt));
		$row = $query->fetch();
		$this->id = $row['id'];
	}

	public function update(){
		$query = DB::connection()->prepare('UPDATE Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES (:nimi, :kansallisuus, :sukupuoli, :syntynyt) RETURNING id');
		$query->execute(array('nimi' => $this->nimi, 'kansallisuus' => $this->kansallisuus, 'sukupuoli' => $this->sukupuoli, 'syntynyt' => $this->syntynyt));
		$row = $query->fetch();
		$this->id = $row['id'];
	}

	public function delete(){
		$query = DB::connection()->prepare('DELETE FROM Kilpailija WHERE id == :id');
		$query->execute();
		// todo: tee tässä jotain olion kadottamiseen
	}

}
