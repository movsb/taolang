function init() {
	const examples_ = examples();
	console.table(examples_);

	var select = document.getElementById('files');
	select.addEventListener('change', function(e){
		var source = form.elements['source'];
		source.value = examples_[select.value];
	});

	var select = document.getElementById('files');
	select.innerHTML = '';
	let keys = Object.keys(examples_).sort();
	for(var i = 0; i < keys.length; i++) {
		var option = document.createElement('option');
		option.value = keys[i];
		option.innerText = keys[i];
		select.appendChild(option);
	}
	
	var event = new Event('change');
	select.dispatchEvent(event);
}

var form = document.getElementById('form');
var source = form.elements['source'];
source.value = '初始化中...';

const go = new Go();
WebAssembly.instantiateStreaming(
    fetch("main.wasm"),
    go.importObject
).then(async (result) => {
    go.run(result.instance);
	setTimeout(init, 500);
});

form.addEventListener('submit', function(e) {
    e.preventDefault();
	if (typeof execute != 'function') {
		alert('WASM 运行时尚未加载成功。');
		return;
	}
    var source = form.elements['source'].value;
    var result = document.getElementById('result');
    result.value = execute(source);
});
