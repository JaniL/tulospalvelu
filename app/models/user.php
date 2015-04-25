<?php
/**
 * Created by PhpStorm.
 * User: jani
 * Date: 19.4.2015
 * Time: 21.01
 */

class User extends BaseModel{
    // Attribuutit
    public $id, $username, $password, $active;

    // Konstruktori
    public function __construct($attributes){
        parent::__construct($attributes);
        // $this->validators = array('validate_nimi', 'validate_kansallisuus', 'validate_sukupuoli');
    }


}