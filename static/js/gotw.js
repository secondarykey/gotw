var map;
$(document).ready( function() {
        /**
         * リサイズ時の処理
         */
        function resize() {
            // 画面の大きさを精一杯にして
            var rootWidth  = $(window).width();
            var rootHeight = $(window).height();

            var workHeight = rootHeight-45;
            $('#canvas').width(rootWidth);
            $('#canvas').height(workHeight);
        }

        resize();
        // リサイズ時のイベント
        $(window).bind('resize', function() {
            resize();
        });

        map = L.map('canvas').setView([38.0,140.0],6)
        L.tileLayer('http://{s}.tile.osm.org/{z}/{x}/{y}.png', {
            attribution: '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors'
        }).addTo(map);

        function put(tweet) {
          var coord = tweet.coordinates.coordinates
          var lon = coord[0]
          var lat = coord[1]

          var marker = L.marker([lat,lon]).addTo(map);
          var content = '<div>' + tweet.Text + '</div>';
          marker.bindPopup(content);
        }

        $("#rain").click(function() {
            $.getJSON("./rain.json", function(json){
                for ( var i = 0; i < json.length; ++i ) {
                    put(json[i])
                }
            });
        });
});


