(function () {
  'use strict';

  angular.module('dashboardApp').controller('ChartController', ['$scope', '$timeout', 'events', 'drawer', function($scope, $timeout, events, drawer) {
    $scope.fields = ["latency", "count_received"];
    $scope.loaded = false;
    $scope.dataAvailable = false;

    $scope.startFetching = function() {
      var conn = new WebSocket("ws://" + window.location.host + "/events");

      conn.onclose = function() {
        console.log("Connection closed.");
      };

      conn.onopen = function() {
        console.log("Connection opened.");
        $scope.startRendering();
      };

      conn.onmessage = function(e) {
        $scope.$apply(function() {
          var message = JSON.parse(e.data);
          $scope.addEvent(message, true);
        });
      };
    };

    $scope.startRendering = function() {
      setInterval(function(){
        for (var chartId in $scope.charts) {
          if (!$scope.charts[chartId].rendered) {
            if ($scope.charts[chartId].svg) {
              $scope.charts[chartId].svg.selectAll('*').remove();
            }

            $scope.drawChart(chartId);
          }
        }
        $scope.stats = [];
        for (var framework in drawer.stats) {
          var stat = drawer.stats[framework];
          stat.framework = framework;
          $scope.stats.push(stat);
        }
        $scope.colors = drawer.colors;
        $scope.stats_count = $scope.stats.sort(function(el) { return -el.avgCount; });
        $scope.stats_latency = $scope.stats.sort(function(el) { return el.avgLatency; });
        $scope.$apply();
      }, 1000);
    };

    $scope.addChart = function(eventName, field) {
      $scope.charts[eventName + field] = {
        events: {},
        allEvents: [],
        field: field,
        eventName: eventName,
        sumLat: 0,
        sumCount: 0,
        count: 0,
        minLat: Infinity,
        maxCount: 0
      };
      $scope.dataAvailable = true;
    };

    $scope.addEvent = function(event, update) {
      event.time = new Date(event.second * 1000);
      event.latency = event.value / event.count_sent;
      for (var i=0; i<$scope.fields.length; i++){
        var chartId = event.eventName + $scope.fields[i];
        if (!$scope.charts[chartId]) {
          $scope.addChart(event.eventName, $scope.fields[i]);
        }
        var events = $scope.charts[chartId].events;
        if (!events[event.framework]) {
          events[event.framework] = {};
        }
        if (update === true && events[event.framework][event.second]) {
          var old_event = events[event.framework][event.second]
          old_event.value += event.value;
          old_event.count_sent += event.count_sent;
          old_event.count_received += event.count_received;
          old_event.latency = old_event.value / old_event.count_sent;
          events[event.framework][event.second] = old_event
        } else {
          events[event.framework][event.second] = event;
        }
        $scope.charts[chartId].allEvents.push(event);
        $scope.charts[chartId].events = events;
        $scope.charts[chartId].rendered = false;
        update = false;
      }
    };

    events.fetch(function(data){
      $scope.charts = {};

      if (data) {
        data.forEach(function(event){
          $scope.addEvent(event);
        });
      }

      $scope.loaded = true;

      for(var chartId in $scope.charts) {
        $scope.drawChart(chartId);
      }

      $scope.startFetching();
    });

    $scope.render = function(last) {
      if (last) {
        $timeout(function() {
          for(var chartId in $scope.charts) {
            $scope.charts[chartId].rendered = false;
            $scope.charts[chartId].init = false;
            $scope.drawChart(chartId);
          }
        }, 0);
      }
    };

    $scope.drawChart = function(chartId) {
      if ($scope.charts[chartId].field === 'latency') {
        drawer.drawChart($scope.charts[chartId], 'line');
      } else {
        drawer.drawChart($scope.charts[chartId], 'area', ['count_received', 'count_sent']);
      }
    };
  }]);
}());
