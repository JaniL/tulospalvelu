<?php

  $routes->get('/', function() {
    HomeController::index();
  });

  $routes->get('/hiekkalaatikko', function() {
    HomeController::sandbox();
  });

  $routes->get('/api/kilpailijat/list', function() {
  	KilpailijaController::list();
  })