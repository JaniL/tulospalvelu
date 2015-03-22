var tulospalvelu = angular.module('tulospalvelu', []);

tulospalvelu.controller('TulospalveluCtrl', function ($scope) {
    $scope.kilpailijat = ["Matti Meikäläinen"];
    $scope.tulokset = [
        {'nimi': 'Matti Meikäläinen',
            'aika': '02:02.22'},
    ];

    $scope.orderProp = 'aika';
});