(function () {

    angular.module('hanoi')
        .controller('ruleController', ['$scope', 'dataService', ruleController]);

    function ruleController($scope, dataService) {
      var self = this;
      $scope.dataService = dataService;

      $scope.typeToName = function(rule){
        switch (rule.type){
          case 'latest':
            return rule.name;
        }
        return rule.type;
      }

      $scope.ruleIcon = function(rule){
        switch (rule.type){
          case 'latest':
            return 'code';
        }
        return ''
      }

      $scope.ruleDetails = function(rule){

        if (rule.explanation_exp){
          return eval('(function (rule,context) {return '+rule.explanation_exp+';})').call(null,rule,$scope);
        }

        switch (rule.type){
          case 'latest':
            return rule.condition;
        }
        return "Cannot get information for unknown type: " + rule.type;
      }

      $scope.tooltip = function(rule){
        var output = 'Sensors: ';
        for (var property in rule.sensors) {
            if (rule.sensors.hasOwnProperty(property)) {
                output += property + " -> " + rule.sensors[property] + "\n";
            }
        }
        return output;
      }
    }
})();
