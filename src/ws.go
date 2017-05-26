package sonda

import (
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"fmt"
)

type WebServer struct {
	Port 	int
	WebSocket chan string
	WSClients []*websocket.Conn
}

func (w *WebServer) Init() {
	w.WebSocket = make(chan string)
	w.WSClients = make([]*websocket.Conn, 0)

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", w.ws)
	http.HandleFunc("/", home)
	go w.initWebSocket()
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%v", w.Port), nil))
}

func (w *WebServer) initWebSocket() {
	for message := range w.WebSocket {
		for _, c := range w.WSClients {
			err := c.WriteMessage(1, []byte(message))
			if err != nil {
				w.closeWs(c)
			}
		}
	}
}

var upgrader = websocket.Upgrader{} // use default options

func (w *WebServer) closeWs(c *websocket.Conn) {
	for i, v := range w.WSClients {
		if v == c {
			w.WSClients = append(w.WSClients[:i], w.WSClients[i+1:]...)
			break
		}
	}
	c.Close()
}

func (w *WebServer) addWs(c *websocket.Conn) {
	w.WSClients = append(w.WSClients, c)
	c.Close()
}

func (w *WebServer) ws(rw http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	w.addWs(c);
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

var homeTemplate = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8">
    <title>Sonda</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css" integrity="sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ" crossorigin="anonymous">
<!--
    <link rel="stylesheet" href="//netdna.bootstrapcdn.com/bootstrap/3.0.3/css/bootstrap.min.css">-->
    <style>body { font-size: 1rem; padding-top: 1rem; } h1 { margin-bottom: 1rem } h4 { font-size: 1.4rem; } p { color: gray; margin-bottom: 0.3rem; margin-top: 0.8rem; } #info h4 { font-size: 1.3rem; color: #444;  font-weight: normal; }</style>
  <head>
  <body>
    <div class="container">
    <h1><span class="hidden-sm-down">Meteosonda </span>Doubrava</h1>

    <hr />

    <div class="row">
    <div class="col-md-6">

    <div class="row">
        <div class="col-6">
        <p>Průměrná rychlost</p>
        <h4 id="speed_average">&nbsp;</h4>
        </div>

       <div class="col-6">
       <p>Poryv</p>
       <h4 id="speed_max">&nbsp;</h4>
       </div>
    </div>

    <div class="row">
       <div class="col-sm-12">
       <p>Průměrný směr</p>
       <h4 id="direction_average">&nbsp;</h4>
       </div>
    </div>


    <div class="hidden-sm-down" id="info">
    <br /><br />

    <div class="row">
       <div class="col-6">
       <p>Teplota CPU</p>
       <h4 id="temperature_cpu">&nbsp;</h4>
       </div>

       <div class="col-6">
       <p>Teplota GPU</p>
       <h4 id="temperature_gpu">&nbsp;</h4>
       </div>
    </div>

    <div class="row">
       <div class="col-6">
       <p>Zatížení (1m)</p>
       <h4 id="load">&nbsp;</h4>
       </div>


       <div class="col-6">
       <p>Doba běhu</p>
       <h4 id="uptime">&nbsp;</h4>
       </div>
    </div>
    </div>

    </div>

    <div class="col-md-6">
	<div class="row">
	    <div class="col-6">
		<p>Aktuální rychlost</p>
		<h4 id="speed_current">&nbsp;</h4>
	    </div>
	    <div class="col-6">
		<p>Aktuální směr</p>
		<h4 id="direction_current">&nbsp;</h4>
	    </div>
	</div>

         <br />
        <p style="position: absolute; z-index: 999">
	<span id="connecting" class="badge badge-warning">připojuji se...</span>
	<span id="online" class="badge badge-success" style="display:none">online</span>
	<span id="closed" class="badge badge-danger" style="display:none">offline</span>
	</p>

	<div id="chart" style="width: 300px; height: 300px;"></div>
    </div>
    </div>

    <hr />
    <p class="help-block">Zdroj dat <a href="http://85.13.85.12/data.json">http://85.13.85.12/data.json</a></p>

    <script src="http://code.highcharts.com/4.0.0/adapters/standalone-framework.js"></script>
    <script src="http://code.highcharts.com/4.0.0/highcharts.js"></script>
    <script src="http://code.highcharts.com/4.0.0/highcharts-more.js"></script>
    <script>
    window.onload = init;
    setInterval(load_data,30000);
    var data = {};
    var chart;

    function init() {
	for(var i = 0; i < 36; i++) {
	    data[i*10] = 0;
	}
	load_data();
	websocket();
	drawChart();
    }

    function addData(direction) {
	data[direction]++;
	serie = new Array();
	var total = 0;
	for(var key in data) {
	    total += data[key];
	}
	for(var key in data) {
	    serie.push(Math.round((data[key] / total) * 100));
	}
	console.log(serie);
	chart.series[0].setData(serie, true);
    }

    function websocket() {
	var ws = new WebSocket('ws://' + window.location.hostname + ':' + window.location.port + '/ws');
	ws.onopen = function() {
	    document.getElementById('connecting').style.display = 'none';
	    document.getElementById('online').style.display = 'inline';
	};
	ws.onmessage = function (evt) {
	    obj = JSON.parse(evt.data);
	    addData(obj.direction_current.toFixed(0));
            document.getElementById('speed_current').innerHTML = obj.speed_current.toFixed(1) + " m/s";
	    document.getElementById('direction_current').innerHTML = obj.direction_current.toFixed(0) + "°";
        };
        ws.onclose = function() {
	    document.getElementById('connecting').style.display = 'none';
	    document.getElementById('online').style.display = 'none';
	    document.getElementById('closed').style.display = 'inline';
	};
    }

    function load_data() {
	var xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function() {
	    if (xhr.readyState == 4) {
		obj = JSON.parse(xhr.responseText);
		document.getElementById('speed_average').innerHTML = obj.speed_average.toFixed(1) + " m/s";
		document.getElementById('speed_max').innerHTML = obj.speed_max.toFixed(1) + " m/s";
		document.getElementById('direction_average').innerHTML = obj.direction_average.toFixed(0) + "°";
		document.getElementById('temperature_cpu').innerHTML = obj.temperature_cpu.toFixed(0) + "°C";
		document.getElementById('temperature_gpu').innerHTML = obj.temperature_gpu.toFixed(0) + "°C";
		document.getElementById('load').innerHTML = obj.load;
		document.getElementById('uptime').innerHTML = convert(obj.uptime);
	    }
	}

	xhr.open('GET', 'data.json', true);
	xhr.send(null);
    };

    function convert(seconds){
	var days = Math.floor(seconds / 86400);
	var hours = Math.floor((seconds - (days * 86400)) / 3600);
	var minutes = Math.floor((seconds - ((hours * 3600) + (days * 86400))) / 60);
	var seconds = seconds - ((days * 86400) + (hours * 3600) + (minutes * 60));
	var result = new String();

	if((days > 0) === true){result += days + 'd';}
	if((hours > 0) === true){result += ' ' + hours + 'h';}
	if((minutes > 0) === true){result += ' ' + minutes + 'm';}
	if((seconds > 0)){result += ' ' + seconds + 's';}
	return result;
    }

    function drawChart() {
    chart = new Highcharts.Chart({
        chart: {
            renderTo: 'chart',
            polar: true
        },
	credits: {
	    enabled: false
	},
        title: false,
        legend: false,

        pane: {
            startAngle: 0,
            endAngle: 360
        },

        xAxis: {
            // gridLineColor: 'transparent', // lines
            lineColor: 'transparent', // oouter circle
            tickInterval: 45,
            min: 0,
            max: 360,
            labels: false,
            /*categories: ['N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW'],*/
            labels: {
        	formatter: function () {
                    var d = new Array('N', 'NE', 'E', 'SE', 'S', 'SW', 'W', 'NW');
        	    return d[this.value/45];
        	}
            }
        },

        yAxis: {
            gridLineColor: 'transparent', // prostredni
            tickColor: 'transparent',
            alternateGridColor: 'transparent',
            lineColor: 'transparent',
            plotLines: false,
            minorTickColor: 'transparent',
            minorGridLineColor: 'transparent',
            labels: false,
            min: 0,
            maxPadding: 0
        },

        plotOptions: {
            series: {
                pointStart: -5,
                pointInterval: 10
            },
            column: {
                pointPadding: 0,
                groupPadding: 0
            }
        },

        tooltip: {
            formatter: function() {
                return 'Směr <b>' + (this.x+5) + '°</b> v <b>' + this.y + '%</b> času';
            }
        },

        series: [{
            type: 'column',
            name: 'Column',
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
            pointPlacement: 'between'
        }]
    });
    }
    </script>
  </body>
</html>
`))