(function () {

    angular.module('hanoi')
        .controller('pageController', ['$scope', 'dataService', pageController]);

    function pageController($scope, dataService) {
      var self = this;
      $scope.dataService = dataService;

      $scope.typeToName = function(rule){
        switch (rule.page.type){
          case 'GMAIL':
            return rule.name;
        }
        return rule.page.type;
      }

      $scope.pageIcon = function(rule){
        switch (rule.page.type){
          case 'GMAIL':
            return 'email';
        }
        return ''
      }

      $scope.ruleStatusIcon = function(ok){
        if (ok)return 'done';
        return 'warning';
      }
      $scope.ruleStatusClass = function(ok){
        if (ok)return 'green lighten-5';
        return 'red lighten-5';
      }

      $scope.pageDetails = function(rule){
        switch (rule.page.type){
          case 'GMAIL':
            return "Page " + rule.page.address;
        }
        return "Cannot get information for unknown type: " + rule.page.type;
      }
    }
})();
