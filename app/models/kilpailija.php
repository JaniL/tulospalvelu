<?php

class Kilpailija extends BaseModel{

    // Attribuutit
	public $id, $nimi, $kansallisuus, $sukupuoli, $syntynyt

  	// Konstruktori
	public function __construct($attributes){
		parent::__construct($attributes);
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

}
