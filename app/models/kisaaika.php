<?php

// require 'app/models/kilpailu.php';
class KisaAika extends BaseModel {
    // Attribuutit
    public $id, $kisaId, $kilpailijaId, $valiaikapiste, $aika;

    // Konstruktori
    public function __construct($attributes){
        parent::__construct($attributes);
        $this->validators = array('validate_valiaika');
    }

    public function validate_valiaika() {
        $errors = array();
        $kilpailu = Kilpailu::find($this->id);

        if ($kilpailu == null) {
            $errors[] = 'Kilpailua ei ole olemassa!';
            return $errors;
        }

        if ($kilpailu->valiaikapisteita < $this->valiaikapiste) {
            $errors[] = 'Väliaikapistettä ei ole olemassa! (väliaikapisteitä on ' . $kilpailu->valiaikapisteita . ')';
        }

        return $errors;
    }

    public static function all(){
        $query = DB::connection()->prepare('SELECT * FROM KisaAika');
        $query->execute();
        $rows = $query->fetchAll();
        $ajat = array();

        foreach($rows as $row){
            $ajat[] = self::luoOlio($row);
        }

        return $ajat;
    }

    public function save(){
        $query = DB::connection()->prepare('INSERT INTO KisaAika (kisaId, kilpailijaId, valiaikapiste, aika) VALUES (:kisaId, :kilpailijaId, :valiaikapiste, :aika) RETURNING id');
        $query->execute(array('kisaId' => $this->kisaId, 'kilpailijaId' => $this->kilpailijaId, 'valiaikapiste' => $this->valiaikapiste, 'aika' => $this->aika));
        $row = $query->fetch();
        $this->id = $row['id'];
    }

    public static function findByKisa($kisa) {
        return self::findByKisaId($kisa['id']);
    }

    public static function findByKisaId($kisaId){
        $query = DB::connection()->prepare('SELECT * FROM KisaAika WHERE kisaId = :kisaId');
        $query->execute(array(
            'kisaId' => $kisaId
        ));
        $rows = $query->fetchAll();
        $ajat = array();

        foreach($rows as $row){
            $ajat[] = self::luoOlio($row);
        }

        return $ajat;
    }

    public function delete(){
        $query = DB::connection()->prepare('DELETE FROM KisaAika WHERE id = :id');
        $query->execute();
        // todo: tee tässä jotain olion kadottamiseen
    }

    public static function deleteCertain($kilpailijaId,$kilpailuId,$valiaikapiste) {
        $query = DB::connection()->prepare('DELETE FROM KisaAika WHERE kilpailijaId = :kilpailijaId AND kisaId = :kilpailuId AND valiaikapiste = :valiaikapiste');
        $query->execute(array(
            'kilpailijaId' => $kilpailijaId,
            'kilpailuId' => $kilpailuId,
            'valiaikapiste' => $valiaikapiste
        ));
    }

    /**
     * @param $row
     * @return Kilpailija
     */
    private static function luoOlio($row)
    {
        return new KisaAika(array(
            'id' => $row['id'],
            'kisaId' => $row['kisaid'],
            'kilpailijaId' => $row['kilpailijaid'],
            'valiaikapiste' => $row['valiaikapiste'],
            'aika' => $row['aika']
        ));
    }
}