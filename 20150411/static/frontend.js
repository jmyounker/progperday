function create_timeline() {
  $.get('/data',
        function(data) {
            if (data['status'] != 'ok') {
                alert("could not read data from server");
                return;
            }
            var container = document.getElementById('timeline');
            var data = new vis.DataSet(data['data'])
            var options = {};
            new vis.Timeline(container, data, options);
        }
    );
}
