<?php

  $routes->get('/', function() {
    HomeController::index();
  });

  $routes->get('/hiekkalaatikko', function() {
    HomeController::sandbox();
  });

  $routes->get('/api/kilpailijat/list', function() {
    KilpailijaController::all();
  });

  $routes->get('/api/kilpailut/list', function() {
      KilpailuController::all();
  });

  $routes->get('/api/kilpailut/id/:kisaId', function($kisaId) {
     KilpailuController::findById($kisaId);
  });

  $routes->post('/api/kilpailut/lisaaSijoitus', function() {
      KilpailuController::lisaaSijoitus();
  });

  $routes->post('/api/kilpailut/lisaaLahtosijoitus', function() {
      KilpailuController::lisaaLahtosijoitus();
  });

  $routes->post('/api/kilpailut/poistaLahtosijoitus', function() {
      KilpailuController::poistaLahtosijoitus();
  });

  $routes->post('/api/kilpailut/vaihdaSijoitukset', function() {
      KilpailuController::vaihdaSijoitukset();
  });

  $routes->post('/api/kilpailut/lisaaAika', function() {
      KilpailuController::lisaaAika();
  });

  $routes->post('/api/kilpailut/poistaAika', function() {
      KilpailuController::poistaAika();
  });

  $routes->post('/api/kilpailijat/lisaa', function() {
    KilpailijaController::lisaa();
  });

  $routes->post('/api/kilpailijat/poista', function() {
      KilpailijaController::poista();
  });

  $routes->post('/api/kilpailut/lisaa', function() {
      KilpailuController::lisaaKisa();
  });

  $routes->post('/api/kilpailut/poista', function() {
      KilpailuController::poistaKisa();
  });