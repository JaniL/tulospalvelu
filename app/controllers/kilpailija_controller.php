<?php
 
  require 'app/models/kilpailija.php';

  class KilpailijaController extends BaseController{

      /**
       * Hakee kaikki kilpailijat ja palauttaa ne jsonina
       */
    public static function all(){
      // make-metodi renderöi app/views-kansiossa sijaitsevia tiedostoja
   	  // View::make('home.html');
      $kilpailijat = Kilpailija::all();
      echo json_encode($kilpailijat);
    }

      /**
       * Lisää kilpailijan tietokantaan. Data otetaan vastaan jsonina.
       */
  public static function lisaa() {
      $data = json_decode(file_get_contents('php://input'), true);

      $nimi = $data['nimi'];
      $kansallisuus = $data['kansallisuus'];
      $syntynyt = $data['syntynyt'];

      $olio = new Kilpailija(array(
          'nimi' => $nimi,
          'kansallisuus' => $kansallisuus,
          'sukupuoli' => 'M',
          'syntynyt' => $syntynyt
      ));

      $olio->save();

      echo "OK";
  }

      /**
       * Poistaa kilpailijan tietokannasta. Vaatii pyynnöltä kilpailija id:n jsonina.
       */
      public static function poista() {
          $data = json_decode(file_get_contents('php://input'), true);

          $id = $data['id'];

          Kilpailija::deleteId($id);

          echo "OK";
      }

      /**
       * Päivittää annetun kilpailijan tiedot.
       */
      public static function paivita() {
          $data = json_decode(file_get_contents('php://input'), true);


          $olio = new Kilpailija($data);

          $olio->update();

          echo "OK";
      }


    /* public static function sandbox(){
      // Testaa koodiasi täällä
      echo 'Hello World!';
    } */
  }
