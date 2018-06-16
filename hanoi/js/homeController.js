(function () {

    angular.module('hanoi')
        .controller('homeController', ['$scope', 'dataService', homeController]);

    function homeController($scope, dataService) {
      var self = this;
      $scope.dataService = dataService;

      $scope.ruleStatusIcon = function(rule){
        if (rule.ok && !rule.noop)return 'done';
        if (rule.noop)return 'loop';
        return 'warning';
      }
      $scope.ruleStatusClass = function(rule){
        if (rule.ok && !rule.noop)return 'green lighten-5';
        if (rule.noop)return 'grey lighten-4';
        return 'red lighten-5';
      }

      $scope.groupCircleIcon = function(group){
        if(dataService.failing_groups[group]){
          return 'live-red slow-blink';
        }
        return 'live-green';
      }
    }
})();
