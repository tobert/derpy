$.ajax({
  url: "/1.0/metric/",
  dataType: "json"
}).done(function (data, status, xhr) {
  console.log(data)

  lseries = [
    { color: 'steelblue', data: [], name: "Read Latency" },
    { color: 'lightblue', data: [], name: "Write Latency" }
  ];

  iseries = [
    { color: 'steelblue', data: [], name: "Read IOPS" },
    { color: 'lightblue', data: [], name: "Write IOPS" }
  ];

  data.forEach(function (m,i) {
    time = new Date(m.Time).getTime()/1000;
    lseries[0].data[i] = { x: time, y: m.ReadMs };
    lseries[1].data[i] = { x: time, y: m.WriteMs };
    iseries[0].data[i] = { x: time, y: m.ReadComplete };
    iseries[1].data[i] = { x: time, y: m.WriteComplete };
  });

  var latency_graph = new Rickshaw.Graph( {
	  element: document.getElementById("latency_chart"),
	  width: 800,
	  height: 400,
	  renderer: 'line',
	  series: lseries
  } );
  addStuff(latency_graph, "latency_legend", "latency_x_axis", "latency_y_axis", true);

  var iops_graph = new Rickshaw.Graph( {
  	element: document.getElementById("iops_chart"),
  	width: 800,
  	height: 400,
  	renderer: 'line',
  	series: iseries
  } );
  addStuff(iops_graph, "iops_legend", "iops_x_axis", "iops_y_axis", false);
});


function addStuff(graph, legendid, xaxisid, yaxisid, bling) {
  var x_ticks, shelving, highlighter;
  var hoverDetail = new Rickshaw.Graph.HoverDetail({ graph: graph });

  var legend = new Rickshaw.Graph.Legend({
	  graph: graph,
	  element: document.getElementById(legendid)
  });

  // ordering requires jquery-ui, but still doesn't seem to work :/
  shelving = new Rickshaw.Graph.Behavior.Series.Toggle({ graph: graph, legend: legend });
  highlighter = new Rickshaw.Graph.Behavior.Series.Highlight({ graph: graph, legend: legend });

  x_ticks = new Rickshaw.Graph.Axis.Time({
    graph: graph,
  });

  var y_ticks = new Rickshaw.Graph.Axis.Y({
    graph: graph,
    orientation: "left",
    element: document.getElementById(yaxisid)
  });

  graph.render();
}
