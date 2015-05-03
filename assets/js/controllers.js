/**
 * Created by jani on 15.4.2015.
 */

function parseDate(date) {
    return date.toISOString().substring(0, 10);
}

var tulospalveluControllers = angular.module('tulospalveluControllers', []);

tulospalveluControllers.controller('KilpailijaController', ['$scope', '$http', function($scope, $http) {
    $http.get('api/kilpailijat/list').success(function(data) {
        $scope.kilpailijat = data;
    });

    $scope.lisaaKilpailija = function() {
        var data = {
            nimi: $scope.kilpailija.nimi,
            kansallisuus: $scope.kilpailija.kansallisuus,
            syntynyt: parseDate($scope.kilpailija.syntymaaika)
        };
        $http.post('api/kilpailijat/lisaa', data).success(function(res) {
            data['id'] = parseInt(res);
            $scope.kilpailijat.push(data);
        });
    }

    $scope.poistaKilpailija = function(id) {
        return;
        $http.post('api/kilpailijat/poista', {
            id: id
        }).success(function(res) {
           $scope.kilpailijat = $scope.kilpailijat.filter(function(k) {
              return k.id != id;
           });
        });
    }

    $scope.editMode = null;

    $scope.muokkaaKilpailija = function() {
        return;
    }
}]);

tulospalveluControllers.controller('KilpailuController', ['$scope', '$http', '$routeParams', function($scope, $http, $routeParams) {

    $scope.getNumber = function(n) {
        return new Array(n);
    };

    $scope.valipiste = 0;

    $scope.showTab = function(n) {
        console.log("triggered: " + n);
        if (n==0) {
            $("#lahtolista").show();
            $("#valipiste").hide();
            return;
        }
        /* $scope.valipisteet = $scope.data.kilpailu.valiaikapisteet.filter(function(v) {
            return v.valiaikapiste == n;
        }).map(function(k) {
            k['kilpailija'] = $scope.etsiKilpailija(k['kilpailijaId']);
            return k;
        }); */
        $scope.valipiste = n;
        console.log($scope.data.kilpailu.valiaikapisteet);
        console.log($scope.valipisteet);
        $("#valipiste").show();
        $("#lahtolista").hide();
    };

    $scope.etsiKilpailija = function(id) {
        var kilpailija = $scope.data.lahtolista.filter(function(k) {
            return k.kilpailija.id == id;
        });
        if (kilpailija.length==0) {
            return null;
        } else {
            return kilpailija[0].kilpailija;
        }
    }

    $scope.lisaaLahtolistaan = function(kilpailija,sijoitus) {
        /* var kilpailija = $scope.kilpailijat.filter(function(k) {
            return k.id = kilpailija;
        })[0]; */

        if ($scope.etsiKilpailija(kilpailija) != null) {
            alert("Kilpailija on jo lähtölistassa!");
            return;
        }

        $http.post('api/kilpailut/lisaaLahtosijoitus', {
            id: kilpailija,
            sijoitus: sijoitus,
            kisaId: $scope.data.kilpailu.id
        }).success(function(data) {
            $scope.haeKilpailu();
        });
    };

    $scope.poistaLahtolistasta = function(kilpailijaId) {
        $http.post('api/kilpailut/poistaLahtosijoitus', {
            id: kilpailijaId,
            kisaId: $scope.data.kilpailu.id
        }).success(function(data) {
            $scope.haeKilpailu();
        })
    };

    $scope.lisaaAika = function() {

        if ($scope.data.kilpailu.valiaikapisteita < $scope.aika.valipiste) {
            alert("Kyseistä välipistettä ei ole olemassa!");
            return;
        }

        if ($scope.data.kilpailu.valiaikapisteet.filter(function(v) {
                return v.kilpailijaId == $scope.aika.nimi && v.valiaikapiste == $scope.aika.valipiste;
            }).length > 0) {
            alert("Kilpailijalle on jo asetettu väliaika!");
            return;
        }

        var ms = ($scope.aika.minuutit*60*1000) + ($scope.aika.sekunnit*1000) + ($scope.aika.millisekunnit);
        $http.post('api/kilpailut/lisaaAika', {
            kilpailuId: $scope.data.kilpailu.id,
            kilpailijaId: $scope.aika.nimi,
            valipiste: $scope.aika.valipiste,
            aika: ms
        }).success(function(data) {
            $scope.haeKilpailu();
        });
    };

    $scope.poistaAika = function(kilpailijaId, valipiste) {
        $http.post('api/kilpailut/poistaAika', {
            kilpailijaId: kilpailijaId,
            kilpailuId: $scope.data.kilpailu.id,
            valipiste: valipiste
        }).success(function(data) {
            $scope.haeKilpailu();
        })
    };

    $scope.aika = function(ms) {
        var min = (ms/1000/60) << 0,
            sec = (ms/1000) % 60;

        min = min.toString();
        if (min.length==1) {
            min = "0" + min;
        }
        sec = sec.toFixed(2).toString();
        if (sec.length==4) {
            sec = "0" + sec;
        }

        return min + ":" + sec;
    }

    $scope.vaihdaSijoitukset = function(k1,k2) {
        $http.post('api/kilpailut/vaihdaSijoitukset', {
            id1: k1,
            id2: k2,
            kisaId: $scope.data.kilpailu.id
        }).success(function(data) {
            $scope.haeKilpailu();
        });
    };

    // todo api get
    $scope.haeKilpailu = function() {
        $http.get('api/kilpailut/id/' + $routeParams.kisaId).success(function(data) {
            console.log(data.kilpailu.valiaikapisteita);
            console.log($scope.getNumber(data.kilpailu.valiaikapisteita));
            console.log(data);
            $scope.data = data;
            $scope.data.kilpailu.valiaikapisteet = $scope.data.kilpailu.valiaikapisteet.map(function(k) {
                k['kilpailija'] = $scope.etsiKilpailija(k['kilpailijaId']);
                return k;
            });
        });
    }

    $scope.haeKilpailu();

    $http.get('api/kilpailijat/list').success(function(data) {
        $scope.kilpailijat = data;
    });
}]);

tulospalveluControllers.controller('KilpailutController', ['$scope', '$http', function($scope, $http) {
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

    $scope.lisaaKisa = function() {
        var data = {
            nimi: $scope.kisa.nimi,
            alkaa: parseDate($scope.kisa.pvm),
            valiaikapisteita: $scope.kisa.valipisteet
        };
        $http.post('api/kilpailut/lisaa', data).success(function(res) {
            data['id'] = parseInt(res);
            $scope.kilpailut.push(data);
        })
    }

    $scope.poistaKisa = function(id) {
        $http.post('api/kilpailut/poista', {
            id: id
        }).success(function(res) {
           $scope.kilpailut = $scope.kilpailut.filter(function(k) {
               return k.id != id;
           });
        });
    }
}]);