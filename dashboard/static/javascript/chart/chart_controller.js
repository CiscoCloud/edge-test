(function () {
  'use strict';

  angular.module('dashboardApp').controller('ChartController', ['$scope', '$timeout', 'events', function($scope, $timeout, events) {
    $scope.fields = ["value", "count"];
    $scope.loaded = false;

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
          $scope.addEvent(message);
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
    };

    $scope.addEvent = function(event) {
      event.time = new Date(event.second * 1000);
      for (var i=0; i<$scope.fields.length; i++){
        var chartId = event.eventName + $scope.fields[i];
        if (!$scope.charts[chartId]) {
          $scope.addChart(event.eventName, $scope.fields[i]);
        }
        var events = $scope.charts[chartId].events;
        if (!events[event.operation]) {
          events[event.operation] = [];
        }
        events[event.operation].push(event);
        events[event.operation].sort(function(a, b){
          return a.second - b.second;
        });
        $scope.charts[chartId].allEvents.push(event);
        $scope.charts[chartId].events = events;
        $scope.charts[chartId].rendered = false;
        if (event.operation == "avg10second") {
          if (event.value < $scope.charts[chartId].minLat) {
            $scope.charts[chartId].minLat = event.value;
          }
          if (event.count > $scope.charts[chartId].maxCount) {
            $scope.charts[chartId].maxCount = event.count;
          }
          $scope.charts[chartId].sumLat += event.value;
          $scope.charts[chartId].sumCount += event.count;
          $scope.charts[chartId].count = $scope.charts[chartId].count + 1;
          $scope.charts[chartId].avgLat = ($scope.charts[chartId].sumLat / $scope.charts[chartId].count).toFixed(2);
          $scope.charts[chartId].avgCount = ($scope.charts[chartId].sumCount / $scope.charts[chartId].count).toFixed(2);
        }
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
      var chart = $scope.charts[chartId];
      if (!chart.init) {
        chart.margin = {top: 20, right: 20, bottom: 30, left: 80};
        chart.width = 550 - chart.margin.left - chart.margin.right;
        chart.height = 250 - chart.margin.top - chart.margin.bottom;

        chart.x = d3.scale.linear()
            .range([0, chart.width]);

        chart.y = d3.scale.linear()
            .range([chart.height, 0]);

        chart.xAxis = d3.svg.axis()
            .scale(chart.x)
            .orient("bottom")
            .tickFormat(function(d){return d3.time.format('%X')(new Date(d));})
            .ticks(9);

        chart.yAxis = d3.svg.axis()
            .scale(chart.y)
            .orient("left");

        chart.line = d3.svg.line()
            .interpolate("basis")
            .x(function(d) { return chart.x(d.time); })
            .y(function(d) { return chart.y(d[chart.field]); });

        chart.svg = d3.select("#chart_" + chart.field + "_" + chart.eventName)
            .attr("width", chart.width + chart.margin.left + chart.margin.right)
            .attr("height", chart.height + chart.margin.top + chart.margin.bottom)
          .append("g")
            .attr("transform", "translate(" + chart.margin.left + "," + chart.margin.top + ")");
        chart.init = true;
      }

      chart.x.domain(d3.extent(chart.allEvents, function(d) {
        return d.time;
      }));
      chart.y.domain(d3.extent(chart.allEvents, function(d) {
        return d[chart.field];
      }));

      chart.svg.append("g")
          .attr("class", "x axis")
          .attr("transform", "translate(0," + chart.height + ")")
          .call(chart.xAxis);

      chart.svg.append("g")
          .attr("class", "y axis")
          .call(chart.yAxis)
          .append("text")
          .attr("transform", "rotate(-90)")
          .attr("y", 6)
          .attr("dy", ".71em")
          .style("text-anchor", "end")
          .text(chart.field);

      var color = d3.scale.category10();
      var operations = ["avg10second", "avg30second", "avg1minute", "avg5minute", "avg10minute", "avg15minute"];
      for (var i = 0; operations < i.length; i++) {
        operations[i] = (i+1).toString();
      }

      for (var operation in chart.events) {
        color.domain(operations);
        var path = chart.svg.append("path")
            .datum(chart.events[operation])
            .attr("class", "line")
            .attr("d", chart.line)
            .style("stroke", function(d) { return color(operation) });
      }

      $scope.charts[chartId] = chart;
      $scope.charts[chartId].rendered = true;
    };
  }]);
}());
