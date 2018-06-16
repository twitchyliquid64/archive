(function () {

    angular.module('hanoi')
        .controller('sensorsController', ['$scope', 'dataService', sensorsController]);

    function sensorsController($scope, dataService) {
      var self = this;
      $scope.dataService = dataService;

      $scope.typeToName = function(sensor){
        switch (sensor.type){
          case 'latency':
            return 'Latency';
          case 'http_status':
            return "HTTP Status Code: " + sensor.domain + sensor.url;
          case 'pushtart_pub':
            if (sensor.component2) {
              return "Pushtart stat: " + sensor.component1 + "." + sensor.component2;
            }
            return "Pushtart Stats RPC: " + sensor.component1;
          case 'pushtart_tart_rpc':
            if (sensor.component3) {
              return "Tart RPC: " + sensor.component1 + "." + sensor.component2 + "." + sensor.component3 + " (" + sensor.push_url + ")";
            }
            if (sensor.component2) {
              return "Tart RPC: " + sensor.component1 + "." + sensor.component2 + " (" + sensor.push_url + ")";
            }
            return "Tart RPC: " + sensor.component1 + " (" + sensor.push_url + ")";
          case 'pushtart_tartstat_rpc':
            if (sensor.component3) {
              return "TartStat RPC: " + sensor.component1 + "." + sensor.component2 + "." + sensor.component3 + " (" + sensor.push_url + ")";
            }
            if (sensor.component2) {
              return "TartStat RPC: " + sensor.component1 + "." + sensor.component2 + " (" + sensor.push_url + ")";
            }
            return "TartStat RPC: " + sensor.component1 + " (" + sensor.push_url + ")";
        }
        return sensor.type;
      }

      $scope.sensorIcon = function(sensor){
        switch (sensor.type){
          case 'latency':
            return 'schedule';
          case 'http_status':
            return 'http';
          case 'pushtart_pub':
            return 'list';
          case 'pushtart_tart_rpc':
            return 'call_made';
          case 'pushtart_tartstat_rpc':
            return 'view_headline';
        }
        return 'code';
      }

      $scope.formatBytes = function(bytes,decimals) {
         if(bytes == 0) return '0B';
         var k = 1024; // or 1024 for binary
         var dm = decimals + 1 || 3;
         var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
         var i = Math.floor(Math.log(bytes) / Math.log(k));
         return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
      }


      $scope.sensorDetails = function(sensor){
        if (!sensor.ok) {
          return sensor.err_msg;
        }

        if (sensor.explanation_exp){
          if (sensor.last_result.val == undefined){
            return 'Sensor not run.';
          }
          return eval('(function (sensor,context) {return '+sensor.explanation_exp+';})').call(null,sensor,$scope);
        }

        switch (sensor.type){
          case 'latency':
            res = '. last result = unknown.';
            if (sensor.last_result.val) {
              res = '. last result = ' + sensor.last_result.val + 'ms.';
            }
            return 'destination = ' + sensor.dest + res;
          case 'http_status':
            if (sensor.last_result.val) {
              return "last status code = " + sensor.last_result.val + ".";
            }
            return "Sensor not run.";
          case 'pushtart_tart_rpc':
            if (sensor.last_result.val == undefined){
              return 'Sensor not run.';
            }
            return String(sensor.last_result.val);
          case 'pushtart_tartstat_rpc':
            if (sensor.last_result.val == undefined){
              return 'Sensor not run.';
            }
            return String(sensor.last_result.val);
        }
        if (sensor.last_result.val) {
          return sensor.last_result.val;
        }
        return "Sensor not run.";      }
    }
})();
