<?php
/**
 * Created by PhpStorm.
 * User: jani
 * Date: 15.4.2015
 * Time: 11.58
 */

// require 'app/models/kilpailu.php';
require 'app/models/lahtolistasijoitus.php';
class KilpailuController extends  BaseController {
    /**
     * Palauttaa kaikki kilpailut jsonina.
     */
    public static function all() {
        $kilpailut = Kilpailu::all();

        header("Content-Type: application/json");
        echo json_encode($kilpailut);
    }

    /**
     * Palauttaa pyydetyn kilpailun tiedot jsonina.
     * @param $kisaId Pyydetyn kilpailun id
     */
    public static function findById($kisaId) {

        $kisa = Kilpailu::find($kisaId);

        if ($kisa == null) {
            echo "{}";
            return;
        }

        $lahtolista = LahtolistaSijoitus::findBykisaId($kisaId);

        header("Content-Type: application/json");
        echo json_encode(array(
            'kilpailu' => $kisa,
            'lahtolista' => $lahtolista,
            // 'valiajat' => array()
        ));

    }

    /**
     * Palauttaa halutun kilpailun lähtölistan jsonina
     * @param $kisaId Kilpailun id
     */
    public static function findLahtolistaByKisaId($kisaId) {
        $sijoitukset = LahtolistaSijoitus::findBykisaId($kisaId);

        header("Content-Type: application/json");
        echo json_encode($sijoitukset);
    }


    /**
     * Lisää ilmoitetulle kilpailijalle lähtösijoituksen haluttuun kilpailuun
     */
    public static function lisaaLahtosijoitus() {
        $data = json_decode(file_get_contents('php://input'), true);

        $kilpailija = $data['id'];
        $sijoitus = $data['sijoitus'];
        $kisaId = $data['kisaId'];

        $olio = new LahtolistaSijoitus(array(
            'kisaId' => $kisaId,
            'kilpailijaId' => $kilpailija,
            'sijoitus' => $sijoitus
        ));

        $olio->save();
        echo "OK";
    }

    /**
     * Poistaa ilmoitetun kilpailijan lähtösijoituksen kilpailusta
     */
    public static function poistaLahtosijoitus() {
        $data = json_decode(file_get_contents('php://input'), true);

        $kilpailija = $data['id'];
        $kisaId = $data['kisaId'];

        LahtolistaSijoitus::deleteWhereKisaIdAndKilpailijaId($kisaId,$kilpailija);

        echo "OK";
    }

    /**
     * Vaihtaa annettujen lähtösijoitusten paikkoja keskenään
     */
    public static function vaihdaSijoitukset() {
        $data = json_decode(file_get_contents('php://input'), true);

        $kilpailija1 = $data['id1'];
        $kilpailija2 = $data['id2'];
        $kisaId = $data['kisaId'];

        LahtolistaSijoitus::vaihdaSijoitukset($kisaId,$kilpailija1,$kilpailija2);

        echo "OK";
    }

    /**
     * Lisää kilpailijalle ajan annetussa kisan välipisteessä
     */
    public static function lisaaAika() {
        $data = json_decode(file_get_contents('php://input'), true);

        /*
         * kilpailuId: $scope.data.kilpailu.id,
            kilpailijaId: $scope.aika.nimi,
            valipiste: $scope.aika.valipiste,
            aika: ms
         */

        $kilpailuId = $data['kilpailuId'];
        $kilpailijaId = $data['kilpailijaId'];
        $valipiste = $data['valipiste'];
        $aika = $data['aika'];

        $olio = new KisaAika(array(
            'kisaId' => $kilpailuId,
            'kilpailijaId' => $kilpailijaId,
            'valiaikapiste' => $valipiste,
            'aika' => $aika
        ));

        $olio->save();

        echo "OK";
    }

    /**
     * Poistaa kilpailijalle annetun ajan
     */
    public static function poistaAika() {
        $data = json_decode(file_get_contents('php://input'), true);

        $kilpailijaId = $data['kilpailijaId'];
        $kilpailuId = $data['kilpailuId'];
        $valipiste = $data['valipiste'];

        KisaAika::deleteCertain($kilpailijaId,$kilpailuId,$valipiste);

        echo "OK";
    }

    /**
     * Lisää järjestelmään kilpailu
     */
    public static function lisaaKisa() {
        $data = json_decode(file_get_contents('php://input'), true);

        $nimi = $data['nimi'];
        $alkaa = $data['alkaa'];
        $valiaikapisteita = $data['valiaikapisteita'];

        $olio = new Kilpailu(array(
            'nimi' => $nimi,
            'alkaa' => $alkaa,
            'valiaikapisteita' => $valiaikapisteita
        ));

        $olio->save();

        echo $olio->id;
    }

    /**
     * Poistaa järjestelmästä pyydetty kilpailu
     */
    public static function poistaKisa() {
        $data = json_decode(file_get_contents('php://input'), true);

        $id = $data['id'];

        $k = Kilpailu::find($id);

        $k->delete();

        echo "OK";
    }

    /**
     * Päivittää kilpailun tiedot.
     */
    public static function muokkaaKisa() {
        $data = json_decode(file_get_contents('php://input'), true);

        $olio = new Kilpailu($data);

        $olio->update();

        echo "OK";
    }
}