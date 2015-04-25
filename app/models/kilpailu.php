<?php
/**
 * Created by PhpStorm.
 * User: jani
 * Date: 15.4.2015
 * Time: 11.44
 */

class Kilpailu extends BaseModel {
    // Attribuutit
    public $id, $nimi, $alkaa, $valiaikapisteita;

    // Konstruktori
    public function __construct($attributes){
        parent::__construct($attributes);
        // $this->validators = array('validate_nimi', 'validate_kansallisuus', 'validate_sukupuoli');
    }

    public static function all(){
        $query = DB::connection()->prepare('SELECT * FROM Kisa');
        $query->execute();
        $rows = $query->fetchAll();
        $kilpailut = array();

        foreach($rows as $row){
            $kilpailut[] = new Kilpailu(array(
                'id' => $row['id'],
                'nimi' => $row['nimi'],
                'alkaa' => $row['alkaa'],
                'valiaikapisteita' => $row['valiaikapisteita']
            ));
        }

        return $kilpailut;
    }

    public static function find($id){
        $query = DB::connection()->prepare('SELECT * FROM Kisa WHERE id = :id LIMIT 1');
        $query->execute(array('id' => $id));
        $row = $query->fetch();

        if($row){
            $kilpailu = new Kilpailu(array(
                'id' => $row['id'],
                'nimi' => $row['nimi'],
                'alkaa' => $row['alkaa'],
                'valiaikapisteita' => $row['valiaikapisteita']
            ));

            return $kilpailu;
        }

        return null;
    }

/*    public function save(){
        $query = DB::connection()->prepare('INSERT INTO Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES (:nimi, :kansallisuus, :sukupuoli, :syntynyt) RETURNING id');
        $query->execute(array('nimi' => $this->nimi, 'kansallisuus' => $this->kansallisuus, 'sukupuoli' => $this->sukupuoli, 'syntynyt' => $this->syntynyt));
        $row = $query->fetch();
        $this->id = $row['id'];
    }*/

/*    public function update(){
        $query = DB::connection()->prepare('UPDATE Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES (:nimi, :kansallisuus, :sukupuoli, :syntynyt) RETURNING id');
        $query->execute(array('nimi' => $this->nimi, 'kansallisuus' => $this->kansallisuus, 'sukupuoli' => $this->sukupuoli, 'syntynyt' => $this->syntynyt));
        $row = $query->fetch();
        $this->id = $row['id'];
    }*/

 /*   public function delete(){
        $query = DB::connection()->prepare('DELETE FROM Kilpailija WHERE id == :id');
        $query->execute();
        // todo: tee tässä jotain olion kadottamiseen
    }*/
}