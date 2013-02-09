function TvShowCtrl($scope, $http) {
  $scope.refreshRooms = function() {
    $http.get("/api/tvshows")
      .success(function(data, status, headers, config) {
        $scope.files = data;
      });
  }
  $scope.refreshRooms();
}