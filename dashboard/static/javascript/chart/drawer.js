(function () {
  'use strict';

  angular.module('dashboardApp').factory('drawer', function(){
    var Drawer = {};

    Drawer.initChart = function(chart, type, fields) {
      chart = Drawer.prepareChart(chart);
      if (type === "area"){
        chart = Drawer.prepareArea(chart, fields);
      } else {
        chart = Drawer.prepareLine(chart);
      }
      chart = Drawer.prepareSVG(chart);
      chart.init = true;
      return chart;
    };

    Drawer.prepareChart = function(chart) {
      chart.margin = {top: 20, right: 20, bottom: 30, left: 60};
      chart.width = 825 - chart.margin.left - chart.margin.right;
      chart.height = 375 - chart.margin.top - chart.margin.bottom;

      chart.x = d3.scale.linear()
          .range([0, chart.width]);

      chart.y = d3.scale.linear()
          .range([chart.height, 0]);

      chart.xAxis = d3.svg.axis()
          .scale(chart.x)
          .orient("bottom")
          .tickFormat(function(d){return d3.time.format('%X')(new Date(d));});

      chart.yAxis = d3.svg.axis()
          .scale(chart.y)
          .orient("left");
      return chart
    };

    Drawer.prepareArea = function(chart, fields) {
      chart.figure = d3.svg.area()
          .x(function(d) { return chart.x(d.time); })
          .y0(function(d) { return chart.y(d[fields[0]]); })
          .y1(function(d) { return chart.y(d[fields[1]]); });
      return chart;
    };

    Drawer.prepareLine = function(chart) {
      chart.figure = d3.svg.line()
          .x(function(d) { return chart.x(d.time); })
          .y(function(d) { return chart.y(d[chart.field]); });
      return chart;
    };

    Drawer.prepareSVG = function(chart) {
      chart.svg = d3.select("#chart_" + chart.field + "_" + chart.eventName)
           .attr("width", chart.width + chart.margin.left + chart.margin.right)
           .attr("height", chart.height + chart.margin.top + chart.margin.bottom)
         .append("g")
           .attr("transform", "translate(" + chart.margin.left + "," + chart.margin.top + ")");
      return chart;
    };

    Drawer.calculateStats = function(data, framework) {
      if (!Drawer.stats) {
        Drawer.stats = {};
      }

      if (!Drawer.stats[framework]) {
        Drawer.stats[framework] = {};
      }
      Drawer.stats[framework].avgCount = d3.mean(data, function(el) {return el.count_received}).toFixed(2);
      Drawer.stats[framework].maxCount = d3.max(data, function(el) {return el.count_received});
      Drawer.stats[framework].avgLatency = d3.mean(data, function(el) {return el.latency}).toFixed(2);
      Drawer.stats[framework].minLatency = d3.min(data, function(el) {return el.latency}).toFixed(2);
    };

    Drawer.drawChart = function(chart, type, fields) {
      if (!chart.init) {
        chart = Drawer.initChart(chart, type, fields);
      }

      chart.xAxis.ticks(5);

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
      var frameworks = ['Golang', 'Dropwizard', 'Finagle', 'Spray', 'Play', 'Unfiltered'];

      for (var framework in chart.events) {
        color.domain(frameworks);
        if (!Drawer.colors) {
          Drawer.colors = {};
        }

        Drawer.colors[framework] = color(framework);
        var data = [];
        for(var key in chart.events[framework]) {
          data.push(chart.events[framework][key]);
        }

        Drawer.calculateStats(data, framework);
        var path = chart.svg.append("path")
            .datum(data)
            .attr("class", type)
            .attr("d", chart.figure)
            .style("stroke", function(d) { return color(framework) });
        if (type === 'area') {
          path.style("fill", function(d) { return color(framework) });
        }

        var pointRadius = 3;

        var tooltip = d3.select('.tooltip');

        var tooltipText = function(el) {
          var time = d3.time.format('%X')(new Date(el.attr('data-time')));
          var text = "";
          if (type === 'line') {
            text += 'Latency: ' + el.attr('data-value') + '<br>';
          }
          if (type === 'area') {
            text += 'Received: ' + el.attr('data-received') + '<br>';
            text += 'Sent: ' + el.attr('data-sent') + '<br>';
          }
          return time + ' ' + el.attr('data-framework') + '<br>' + text;
        };

        for (var i = 0; i < data.length; i+=60) {
          var point = chart.svg.append('circle')
            .attr('class', 'data-point')
            .style('opacity', 0.5)
            .attr('cx', function() { return chart.x(data[i].time); })
            .attr('cy', function() { return chart.y(data[i][chart.field]); })
            .attr('r', function() { return pointRadius; })
            .attr('data-framework', data[i].framework)
            .attr('data-time', data[i].time)
            .attr('data-value', data[i].latency)
            .attr('data-received', data[i].count_received)
            .attr('data-sent', data[i].count_sent);
          point.on('mouseover', function(){
            tooltip
              .style("left", d3.event.pageX + "px")
              .style("top", (d3.event.pageY - 10) + "px")
              .style('opacity', 1)
              .html(tooltipText(d3.select(this)));
          });
          point.on('mouseout', function() {
            tooltip
              .style("left", 0)
              .style("top", 0)
              .style('opacity', 0);
          });
        }

      }

      chart.rendered = true;
    };

    return Drawer;
  });
}());
