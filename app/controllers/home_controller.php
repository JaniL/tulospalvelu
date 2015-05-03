<?php

    require 'app/models/lahtolistasijoitus.php';
  class HomeController extends BaseController{

    public static function index(){
      // make-metodi renderöi app/views-kansiossa sijaitsevia tiedostoja
   	  View::make('home.html');
    }

    public static function sandbox(){
      // Testaa koodiasi täällä
        echo json_encode(LahtolistaSijoitus::findBykisaId(1));
      //echo 'Hello World!';
    }
  }
