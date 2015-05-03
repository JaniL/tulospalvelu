/**
 * Created by jani on 15.4.2015.
 */
var tulospalveluApp = angular.module('tulospalveluApp', [
    'ngRoute',
    'angular-jwt',
    'tulospalveluControllers',
    'angular-loading-bar'
]);

/**
 * Määritellään angularin reititys
 */
tulospalveluApp.config(['$routeProvider','$httpProvider', 'jwtInterceptorProvider', function($routeProvider, $httpProvider, jwtInterceptorProvider) {
    /*jwtInterceptorProvider.tokenGetter = function() {
        return localStorage.getItem('id_token');
    };
    $httpProvider.interceptors.push('jwtInterceptor');*/
    $routeProvider.
        when('/kilpailijat', {
            templateUrl: 'assets/partials/kilpailijat.htm',
            controller: 'KilpailijaController'
        }).
        when('/kilpailut', {
            templateUrl: 'assets/partials/kilpailut.htm',
            controller: 'KilpailutController'
        }).
        when('/kisa/:kisaId', {
            templateUrl: 'assets/partials/kilpailu.htm',
            controller: 'KilpailuController'
        }).
        otherwise({
            redirectTo: '/kilpailut'
        });
}]);