<?php
/**
 * Created by PhpStorm.
 * User: jani
 * Date: 15.4.2015
 * Time: 11.58
 */

require 'app/models/kilpailu.php';
class KilpailuController extends  BaseController {
    public static function all() {
        $kilpailut = Kilpailu::all();
        echo json_encode($kilpailut);
    }
}