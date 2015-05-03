<?php

require 'app/models/kilpailu.php';
require 'app/models/kilpailija.php';
class LahtolistaSijoitus extends BaseModel {
    // Attribuutit
    public $id, $kisaId, $kilpailija, $sijoitus, $kilpailijaId;

    // Konstruktori
    public function __construct($attributes){
        parent::__construct($attributes);
        // $this->validators = array('validate_nimi', 'validate_kansallisuus', 'validate_sukupuoli');
    }

    public static function findBykisaId($kisaId) {
        $query = DB::connection()->prepare('SELECT * FROM KisaLahtoLista INNER JOIN Kilpailija on KisaLahtoLista.kilpailijaId=Kilpailija.id INNER JOIN Kisa on KisaLahtoLista.kisaId=Kisa.id WHERE kisaId = :id');
        $query->execute(array('id' => $kisaId));
        $rows = $query->fetchAll();
        $lahtolista = array();

        foreach($rows as $row) {
            $lahtolista[] = self::luoOlio($row);
        }

        return $lahtolista;
    }

    public static function findBykisaIdAndSijoitus($kisaId,$sijoitus) {
        $query = DB::connection()->prepare('SELECT * FROM KisaLahtoLista INNER JOIN Kilpailija on KisaLahtoLista.kilpailijaId=Kilpailija.id INNER JOIN Kisa on KisaLahtoLista.kisaId=Kisa.id WHERE kisaId = :id AND sijoitus = :sijoitus LIMIT 1');
        $query->execute(array(
            'id' => $kisaId,
            'sijoitus' => $sijoitus
        ));
        $row = $query->fetch();

        if ($row) {
            return self::luoOlio($row);
        }

        return null;
    }

    public function save(){

        Kint::dump(self::findBykisaIdAndSijoitus($this->kisaId,$this->sijoitus));
        if (self::findBykisaIdAndSijoitus($this->kisaId,$this->sijoitus) != null) {
            $query = DB::connection()->prepare('UPDATE KisaLahtoLista SET sijoitus = (sijoitus+1) WHERE sijoitus >= :sijoitus AND kisaId = :kisaId');
            $query->execute(array('sijoitus' => $this->sijoitus, 'kisaId' => $this->kisaId));
        }

        $query = DB::connection()->prepare('INSERT INTO KisaLahtolista (kisaId, kilpailijaId, sijoitus) VALUES (:kisaId, :kilpailijaId, :sijoitus) RETURNING id');
        $query->execute(array(
            'kisaId' => $this->kisaId,
            'kilpailijaId' => $this->kilpailijaId,
            'sijoitus' => $this->sijoitus
        ));
        $row = $query->fetch();
        $this->id = $row['id'];
    }

    public function update(){
        $query = DB::connection()->prepare('UPDATE Kilpailija (kisaId, kilpailijaId, sijoitus) VALUES (:kisaId, :kilpailijaId, :sijoitus) RETURNING id');
        $query->execute(array(
            'kisaId' => $this->kisaId,
            'kilpailijaId' => $this->kilpailija->id,
            'sijoitus' => $this->sijoitus
        ));
        $row = $query->fetch();
        $this->id = $row['id'];
    }

    /**
     * @param $row
     * @return LahtolistaSijoitus
     */
    private static function luoOlio($row)
    {
        return new LahtolistaSijoitus(array(
            'id' => $row['id'],
            /*'kisa' => new Kilpailu(array(
                'id' => $row['9'],
                'nimi' => $row['10'],
                'alkaa' => $row['11'],
                'valiaikapisteita' => $row['12']
            )),*/
            'kisaId' => $row['1'],
            'kilpailija' => new Kilpailija(array(
                'id' => $row['kilpailijaid'],
                'nimi' => $row['5'],
                'kansallisuus' => $row['kansallisuus'],
                'sukupuoli' => $row['sukupuoli'],
                'syntynyt' => $row['syntynyt']
            )),
            'kilpailijaId' => $row['kilpailijaid'],
            'sijoitus' => $row['sijoitus']
        ));
    }

    public function asetaSijoitus($uusiSijoitus) {
        // todo
    }

    public function delete(){
        $tempSijoitus = $this['sijoitus'];
        $query = DB::connection()->prepare('DELETE FROM KisaLahtoLista WHERE id = :id');
        $query->execute();

        $query = DB::connection()->prepare('UPDATE KisaLahtoLista SET sijoitus = (sijoitus-1) WHERE sijoitus > :sijoitus');
        $query->execute(array('sijoitus' => $tempSijoitus));
        // todo: tee tässä jotain olion kadottamiseen
    }

    public static function deleteWhereKisaIdAndKilpailijaId($kisa, $kilpailija) {
        $query = DB::connection()->prepare('DELETE FROM KisaLahtoLista WHERE kisaId = :id AND kilpailijaId = :kilId RETURNING sijoitus');
        $query->execute(array(
            'id' => $kisa,
            'kilId' => $kilpailija
        ));
        $row = $query->fetch();
        $sijoitus = $row['sijoitus'];

        $query = DB::connection()->prepare('UPDATE KisaLahtoLista SET sijoitus = (sijoitus-1) WHERE sijoitus > :sijoitus');
        $query->execute(array('sijoitus' => $sijoitus));
    }



    public static function vaihdaSijoitukset($kisaId,$sijoitus1,$sijoitus2) {
        $query = DB::connection()->prepare('UPDATE KisaLahtoLista SET sijoitus = 0 WHERE sijoitus = :sijoitus1 AND kisaId = :kisaId');
        $query->execute(array(
            'sijoitus1' => $sijoitus1,
            'kisaId' => $kisaId
        ));

        $query = DB::connection()->prepare('UPDATE KisaLahtoLista SET sijoitus = :sijoitus1 WHERE sijoitus = :sijoitus2 AND kisaId = :kisaId');
        $query->execute(array(
            'sijoitus1' => $sijoitus1,
            'sijoitus2' => $sijoitus2,
            'kisaId' => $kisaId
        ));

        $query = DB::connection()->prepare('UPDATE KisaLahtoLista SET sijoitus = :sijoitus2 WHERE sijoitus = 0 AND kisaId = :kisaId');
        $query->execute(array(
            'sijoitus2' => $sijoitus2,
            'kisaId' => $kisaId
        ));
    }

    public function swap($sijoitus) {
        LahtolistaSijoitus::vaihdaSijoitukset($this['kisaId'],$this['sijoitus'],$sijoitus['sijoitus']);
    }
}