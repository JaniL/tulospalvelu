/**
 * Created by jani on 19.4.2015.
 */

var user = angular.module('user', []);

user.factory('AuthService', function($http, Session) {
   var authService = {};

    authService.login = function (credentials) {
        return $http
            .post('api/login', credentials)
            .then(function(res) {
                // todo
            });
    };

    authService.register = function (credentials) {
        return $http
            .post('api/register', credentials)
            .then(function(res) {
                // todo
            });
    };
});

user.service('JWT', function () {
    this.create = function (token) {
        this.token = token;
    };
    this.destroy = function () {
        this.token = null;
    };
})