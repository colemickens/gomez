var gomez = angular.module('gomez', []);

gomez.config(function($routeProvider, $locationProvider) {
    //$locationProvider.html5Mode(true);

    $routeProvider.
        when('/',        { templateUrl: '/partials/home.html',    controller: HomeCtrl    }).
        when('/log',     { templateUrl: '/partials/log.html',     controller: LogCtrl     }).
        when('/movies',  { templateUrl: '/partials/movies.html',  controller: MovieCtrl   }).
        when('/tvshows', { templateUrl: '/partials/tvshows.html', controller: TvShowCtrl  }).
        when('/music',   { templateUrl: '/partials/music.html',   controller: MusicCtrl   }).

        otherwise({redirectTo: '/'});
});

gomez.filter('posterurl', function() {
	return function(input) {
		if(input == "") {
			return "/images/missing.jpg";
		}
		return input;
	};
});

gomez.filter('duration', function() {
	return function(input) {
		d = Number(input);
		var h = Math.floor(d / 3600);
		var m = Math.floor(d % 3600 / 60);
		var s = Math.floor(d % 3600 % 60);
		return ((h > 0 ? h + " hours " : "") + (m > 0 ? (h > 0 && m < 10 ? "0" : "") + m + " minutes " : "0:") + (s < 10 ? "0" : "") + s + " seconds");
	};
});

gomez.filter('filesize', function() {
	return function(input, precision) {
		bytes = Number(input);
		var kilobyte = 1024;
		var megabyte = kilobyte * 1024;
		var gigabyte = megabyte * 1024;
		var terabyte = gigabyte * 1024;

		if ((bytes >= 0) && (bytes < kilobyte)) {
			return bytes + ' B';
		} else if ((bytes >= kilobyte) && (bytes < megabyte)) {
			return (bytes / kilobyte).toFixed(precision) + ' KB';
		} else if ((bytes >= megabyte) && (bytes < gigabyte)) {
			return (bytes / megabyte).toFixed(precision) + ' MB';
		} else if ((bytes >= gigabyte) && (bytes < terabyte)) {
			return (bytes / gigabyte).toFixed(precision) + ' GB';
		} else if (bytes >= terabyte) {
			return (bytes / terabyte).toFixed(precision) + ' TB';
		} else {
			return bytes + ' B';
		}
		return input;
	};
});

angular.element(document).ready(function() { 
	//$(document).foundationTooltips();
	$(document).foundationNavigation();
	$(document).foundationTopBar();
	$(document).foundationCustomForms();
	console.log("activated TopBar");
});