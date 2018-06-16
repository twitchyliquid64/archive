(function() {

    var app = angular.module('hanoi').factory('dataService', ['$rootScope', '$http', '$interval', '$window', dataService]);
7
    GET_SENSORS_URL = "/sensors"
    GET_RULES_URL = "/rules"
    GET_PAGERS_URL = "/paging_rules"
    GET_DIFF_URL = "/diff"

    function dataService($rootScope, $http, $interval, $window){
      var self = this;
      self.sensors = {};
      self.rules = {};
      self.paging_rules = {};
      self.failing_groups = {};
      self.update_key = 0;
      self.loaded = false;
      self.load_in_flight = false;

      self.fetchPagingRules = function(){
        return $http.get(GET_PAGERS_URL, {}).then(function (response) {
          self.paging_rules = response.data;
          $rootScope.$broadcast('pagers-loaded');
        }, self.errorCb);
      };

      self.fetchSensors = function(){
        return $http.get(GET_SENSORS_URL, {}).then(function (response) {
          self.sensors = response.data;
          $rootScope.$broadcast('sensors-loaded');
        }, self.errorCb);
      };

      self.fetchRules = function(){
        return $http.get(GET_RULES_URL, {}).then(function (response) {
          self.rules = response.data;
          self.failing_groups = {};
          Object.keys(self.rules).forEach(function(group,index) {
            for (var i = 0; i < self.rules[group].length; i++){
              rule = self.rules[group][i];
              if (!rule.ok){
                self.failing_groups[group] = true;
              }
            }
          });
          $rootScope.$broadcast('rules-loaded');
        }, self.errorCb);
      };

      self.errorCb = function(response) {
        self.load_in_flight = false;
        console.log("Err: ", response);
      }

      self.updateComponent = function(component){
        switch (component){
          case 'all':
            Promise.all([self.fetchSensors(), self.fetchRules(), self.fetchPagingRules()]).then(function(ok){
              self.loaded = true;
            });
          case 'sensors':
            self.fetchSensors();
          case 'rules':
            self.fetchRules();
          case 'pagers':
          case 'paging_rules':
            self.fetchPagingRules();
        }
      }

      self.fetchAndApplyDiff = function(){
       if (self.load_in_flight)return;
       self.load_in_flight = true;

        $http.get(GET_DIFF_URL + "?update_key=" + self.update_key, {}).then(function (response) {
          self.load_in_flight = false;
          if (!response.data.up_to_date){
            for (var i = 0; i < response.data.updates.length; i++) {
              self.updateComponent(response.data.updates[i]);
            }
          }
          self.update_key = response.data.key;
        }, self.errorCb);
      }

      $interval(self.fetchAndApplyDiff, 2000);
      self.updateComponent('all');
      return self;
    };

})();
