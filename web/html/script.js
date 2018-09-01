var form = document.getElementById('form');
form.addEventListener('submit', function(e) {
    e.preventDefault();
    var source = form.elements['source'].value;
    var result = document.getElementById('result');
    var data = {
        source: source
    };
    var xhr = new XMLHttpRequest();
    xhr.open('POST', '/v1/execute');
    xhr.onload = function () {
        if (xhr.status == 200) {
            result.value = xhr.responseText;
            result.style.color = 'unset';
        } else {
            result.value = xhr.responseText;
            result.style.color = 'red';
        }
    };
    xhr.onerror = function(e) {
        alert('error');
    }
    xhr.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    result.value = "Waiting...";
    xhr.send(JSON.stringify(data));
});

window.addEventListener('load', function(e){
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/v1/examples');
    xhr.onload = function () {
        if (xhr.status == 200) {
            var files = JSON.parse(xhr.responseText);
            var select = document.getElementById('files');
            select.innerHTML = '';
            for(var i = 0; i < files.length; i++) {
                var option = document.createElement('option');
                option.value = files[i];
                option.innerText = files[i];
                select.appendChild(option);
            }
            var event = new Event('change');
            select.dispatchEvent(event);
        }
    };
    xhr.onerror = function(e) {
        alert('error');
    }
    xhr.send();
});

var select = document.getElementById('files');
select.addEventListener('change', function(e){
    var source = form.elements['source'];
    var xhr = new XMLHttpRequest();
    xhr.open('GET', '/v1/examples/'+encodeURI(select.value));
    xhr.onload = function () {
        if (xhr.status == 200) {
            source.value = xhr.responseText;
        }
    };
    xhr.onerror = function(e) {
        alert('error');
    }
    source.value = "Waiting...";
    xhr.send();
});
