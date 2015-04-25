/**
 * Created by jani on 15.4.2015.
 */

var tulospalveluControllers = angular.module('tulospalveluControllers', []);

tulospalveluControllers.controller('KilpailijaController', ['$scope', '$http', function($scope, $http) {
    $http.get('api/kilpailijat/list').success(function(data) {
        $scope.kilpailijat = data;
    });
}]);

tulospalveluControllers.controller('KilpailuListaController', ['$scope', '$http', function($scope, $http) {
    // todo api get
}]);

tulospalveluControllers.controller('KilpailuController', ['$scope', '$http', function($scope, $http) {
    // todo
    $scope.kilpailut = [];
    console.log("lel");
    $http.get('api/kilpailut/list').success(function(data) {
        console.log(data);
        $scope.kilpailut = data;
    }).error(function(data, status, headers, config) {
        // called asynchronously if an error occurs
        // or server returns response with an error status.
        console.log("err");
        console.log(data);
        console.log(status);
        console.log(headers);
    });
}]);