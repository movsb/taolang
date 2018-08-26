package main

var source = `
function main() {
	let a;
	let b = nil;
	let c = true;
	let d = 123;
	let e = "str";
	let f = function(x,y,z) {return x+y*z;};
	let g = function() {return "test";}();

	print(a,b,c,d,e,f,g);

	nil = 0;
}
`

func main() {
	tokenizer := NewTokenizer(source)
	parser := NewParser(tokenizer)
	program := parser.Parse()
	program.Execute()
}
