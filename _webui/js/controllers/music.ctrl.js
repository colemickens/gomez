function MusicCtrl($scope, $http){ 
    $scope.refreshRooms = function() {
    $http.get("/api/music")
      .success(function(data, status, headers, config) {
        $scope.files = data;
      });
  }

  $scope.genres = [ "music-for-whiners", "good musics" ];
}