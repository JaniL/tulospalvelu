<?php

  class DatabaseConfig{

    // Valitse käyttämäsi tietokantapalvelin - PostgreSQL (psql) tai MySQL (mysql)
    private static $default_database = 'psql';

    private static function env($key, $default = null) {
      $value = getenv($key);

      if($value === false || $value === '') {
        return $default;
      }

      return $value;
    }

    // Muuta users-ympäristöä asettamalle oikeat arvot KAYTTAJATUNNUS-kohtaan (käyttäjätunnuksesi)
    // ja SALASANA-kohtaan (tietokantasi pääkäyttäjän salasana)
    private static $connection_config = array(
      'psql' => array(
        'resource' => 'pgsql:host=localhost;port=5432;dbname=tulospalvelu',
        'username' => 'postgres',
        'password' => 'postgres'
      ),
      'mysql' => array(
        'resource' => 'mysql:host=localhost;port=3306;dbname=mysql',
        'username' => 'root',
        'password' => 'SALASANA'
      )
    );

    public static function connection_config(){
      $db = self::env('DB_DRIVER', self::$default_database);

      if(!isset(self::$connection_config[$db])) {
        $db = self::$default_database;
      }

      $config = self::$connection_config[$db];

      if($db == 'psql') {
        $host = self::env('DB_HOST', 'localhost');
        $port = self::env('DB_PORT', '5432');
        $name = self::env('DB_NAME', 'tulospalvelu');
        $config['resource'] = 'pgsql:host=' . $host . ';port=' . $port . ';dbname=' . $name;
      }

      if($db == 'mysql') {
        $host = self::env('DB_HOST', 'localhost');
        $port = self::env('DB_PORT', '3306');
        $name = self::env('DB_NAME', 'mysql');
        $config['resource'] = 'mysql:host=' . $host . ';port=' . $port . ';dbname=' . $name;
      }

      $config['username'] = self::env('DB_USER', $config['username']);
      $config['password'] = self::env('DB_PASSWORD', $config['password']);

      $config = array(
        'db' => $db,
        'config' => $config
      );

      return $config;
    }

  }
