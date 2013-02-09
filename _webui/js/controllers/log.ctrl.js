function LogCtrl($scope, $http) {
  $scope.refreshLog = function() {
    $http.get("/log").success(function(data, status, headers, config) {
        $scope.logfile = data;
    });
  }
  $scope.interv = setInterval($scope.refreshLog, 2000);
  $scope.refreshLog();
}