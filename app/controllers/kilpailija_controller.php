<?php

  class KilpailijaController extends BaseController{

    public static function list(){
      // make-metodi renderöi app/views-kansiossa sijaitsevia tiedostoja
   	  // View::make('home.html');
      $kilpailijat = Kilpailija::all();
      echo json_encode($kilpailijat);
    }

    /* public static function sandbox(){
      // Testaa koodiasi täällä
      echo 'Hello World!';
    } */
  }
