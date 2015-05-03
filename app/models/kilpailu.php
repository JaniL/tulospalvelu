<?php
/**
 * Created by PhpStorm.
 * User: jani
 * Date: 15.4.2015
 * Time: 11.44
 */

require 'app/models/kisaaika.php';
class Kilpailu extends BaseModel {
    // Attribuutit
    public $id, $nimi, $alkaa, $valiaikapisteita, $valiaikapisteet;

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
            $kilpailut[] = self::luoOlio($row);
        }

        return $kilpailut;
    }

    public static function find($id){
        $query = DB::connection()->prepare('SELECT * FROM Kisa WHERE id = :id LIMIT 1');
        $query->execute(array('id' => $id));
        $row = $query->fetch();

        if($row){
            $kilpailu = self::luoOlio($row);

            return $kilpailu;
        }

        return null;
    }

    /**
     * @param $row
     * @return Kilpailu
     */
    private static function luoOlio($row)
    {
        return new Kilpailu(array(
            'id' => $row['id'],
            'nimi' => $row['nimi'],
            'alkaa' => $row['alkaa'],
            'valiaikapisteita' => $row['valiaikapisteita'],
            'valiaikapisteet' => KisaAika::findByKisaId($row['id'])
            // 'valiaikapisteet' => self::generoiValiaikaArray(KisaAika::findByKisaId($row['id']))
        ));
    }

    /* private static function generoiValiaikaArray($findByKisaId)
    {
        return null;
    } */

    public function save(){
            $query = DB::connection()->prepare('INSERT INTO Kisa (nimi, alkaa, valiaikapisteita) VALUES (:nimi, :alkaa, :valiaikapisteita) RETURNING id');
            $query->execute(array(
                'nimi' => $this->nimi,
                'alkaa' => $this->alkaa,
                'valiaikapisteita' => $this->valiaikapisteita

            ));
            $row = $query->fetch();
            $this->id = $row['id'];
    }

/*    public function update(){
        $query = DB::connection()->prepare('UPDATE Kilpailija (nimi, kansallisuus, sukupuoli, syntynyt) VALUES (:nimi, :kansallisuus, :sukupuoli, :syntynyt) RETURNING id');
        $query->execute(array('nimi' => $this->nimi, 'kansallisuus' => $this->kansallisuus, 'sukupuoli' => $this->sukupuoli, 'syntynyt' => $this->syntynyt));
        $row = $query->fetch();
        $this->id = $row['id'];
    }*/

    public function delete(){
        $tempId = $this->id;
        $query = DB::connection()->prepare('DELETE FROM Kisa WHERE id = :id');
        $query->execute(array(
            'id' => $tempId
        ));

        $query = DB::connection()->prepare('DELETE FROM KisaAika WHERE kisaId = :id');
        $query->execute(array(
            'id' => $tempId
        ));

        $query = DB::connection()->prepare('DELETE FROM KisaLahtoLista WHERE kisaId = :id');
        $query->execute(array(
            'id' => $tempId
        ));
        // todo: tee tässä jotain olion kadottamiseen
    }
}