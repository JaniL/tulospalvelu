var tulospalvelu = angular.module('tulospalvelu', []);

tulospalvelu.controller('TulospalveluCtrl', ['$scope', '$filter', function ($scope,$filter) {
    $scope.kilpailijat = ["Matti Meikäläinen"];
    $scope.tulokset = [
        {'nimi': 'Matti Meikäläinen',
            'aika': 122022},
    ];

    $scope.lahtolista = ["Matti Meikäläinen", "Jukka Kovahiihtäjä", "Dope Matti"];

    $scope.aika = function(ms) {
        var min = (ms/1000/60) << 0,
            sec = (ms/1000) % 60;

        return "0" + min + ":0" + sec.toFixed(2);
    }

    $scope.lisaaAika = function(aika) {
        $scope.tulokset.push({
            'nimi': aika.nimi,
            'aika': aika.millisekunnit + (aika.sekunnit*1000) + (aika.minuutit*60*1000)
        });

        $scope.tulokset = $filter('orderBy')($scope.tulokset, 'aika');
    }

    $scope.lisaaLahtolistaan = function(kilpailija) {
        $scope.lahtolista.splice(kilpailija.lahtosijoitus-1,0,kilpailija.nimi);
    }

    $scope.orderProp = 'aika';
}]);

tulospalvelu.controller('KilpailijaCtrl', function($scope) {
    $scope.kilpailijat = [
        {
            'nimi': 'Matti V',
            'kansallisuus': 'FI',
            'syntymaaika': new Date(1992,11,11).toLocaleDateString()
        }
    ];

    $scope.lisaaKilpailija = function(kilpailija) {
        $scope.kilpailijat.push({
            'nimi': kilpailija.nimi,
            'kansallisuus': kilpailija.kansallisuus,
            'syntymaaika': kilpailija.syntymaaika.toLocaleDateString()
        });
    }
})