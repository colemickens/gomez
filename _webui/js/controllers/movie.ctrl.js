

function MovieCtrl($scope, $location, $http) {
  xhr2 = new XMLHttpRequest();

  $scope.refreshRooms = function() {   
    $http.get("/api/movies")
      .success(function(data, status, headers, config) {
        $scope.movies = data;
        console.log(data);
      });
  }

  $scope.sortMethod = function(movie) {
    console.log(movie);
    if (movie.title != "") {
      return movie.title;
    } else {
      return movie.search_string;
    }
  }

  $scope.play = function() { $scope.jwplayer.play(); }
  $scope.seek = function() {}
  $scope.fullscreen = function() {
    // make that div fullscreen
  }

  $scope.revealPlayer = function(movie_id) {
    $scope.format = "flv";
    $scope.offset = 0;
    $scope.width = 0;
    $scope.height = 0;
    $scope.jwplayer = jwplayer("mediaplayer");


    $scope.player = $("#mediaplayer").first();
    
    $http.get("/api/info?id=" + movie_id).success(function(data, status, headers, config) {
      $scope.movie = data;

      $scope.width = videoStream($scope.movie).width;
      $scope.height = videoStream($scope.movie).height;

      $("#player").reveal();

      $scope.swap_video(true);

      $("#seekbar").on("change", function() {
        if($scope.jwplayer.getState() != "PAUSED"){
          $scope.jwplayer.pause();
        }
      });
      $("#seekbar").on("change", _.debounce($scope.seek, 1000));

      $scope.jwplayer.onTime(function(e) {
        $scope.position = e.position;
        $("#seekbar")[0].value = $scope.offset + Math.round($scope.position);
      });
      
      $("#seekbar").attr('max', $scope.movie.ffprobe.format.duration);
      $("#seekbar").val(0);

      //TODO: Attach this to onPlay/onPause?
      clearInterval($scope.interval);
    });

    $scope.swap_video = function(e) {
      var _oldWidth = $scope.width;
      $scope.width = $("#playercont").width();
      $scope.height = $scope.height * ($scope.width/_oldWidth);
      $scope.streamUrl = "/api/stream?id=" + $scope.movie.id;
      $scope.streamUrl += "&format=" + $scope.format;
      $scope.streamUrl += "&width=" + Math.round($scope.width);
      $scope.streamUrl += "&height=" + Math.round($scope.height);
      $scope.streamUrl += "&offset=" + $scope.offset;
      $scope.streamUrl += "&bitrate=" + $scope.bitrate;
      // swap_video on anything changing

      if(e) {
        $scope.jwplayer.setup({
          flashplayer: "/lib/mediaplayer-5.10/player.swf",
          controlbar: "none",
        });
      }
      $scope.jwplayer.resize($scope.width, $scope.height);
      $scope.jwplayer.load([{
        file: $scope.streamUrl,
        duration: $scope.movie.ffprobe.format.duration - $scope.offset,
        provider: "video",
      }]);

      $scope.jwplayer.play();
    };

    $scope.seek = function() {
      $scope.offset = parseInt($("#seekbar").val());
      console.log("new offset", $scope.offset);
      $scope.swap_video(false);
      /*
      $("#video")[0].load();
      $("#video")[0].play();
      $("#seekbar").value = value;
      */
    };
  };

  $scope.revealInfo = function(movie_id) {
    $http.get("/api/info?id=" + movie_id).success(function(data, status, headers, config) {
      $scope.movie = data;
      $("#info").reveal();
    });
  };

  $scope.interval = setInterval($scope.refreshRooms, 10000)

  $scope.genres = [ "movies-for-whiners", "good movies" ];
  $scope.refreshRooms();
}