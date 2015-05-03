<?php
 
  require 'app/models/kilpailija.php';

  class KilpailijaController extends BaseController{

    public static function all(){
      // make-metodi renderöi app/views-kansiossa sijaitsevia tiedostoja
   	  // View::make('home.html');
      $kilpailijat = Kilpailija::all();
      echo json_encode($kilpailijat);
    }

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

      public static function poista() {
          $data = json_decode(file_get_contents('php://input'), true);

          $id = $data['id'];

          Kilpailija::deleteId($id);

          echo "OK";
      }


    /* public static function sandbox(){
      // Testaa koodiasi täällä
      echo 'Hello World!';
    } */
  }
