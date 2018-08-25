package main

var source = `
function f1(a,b,c) {

}

function main() {
	function f1(a,b,c) {

	}
	var a = 133;
	print("string", 123, true, nil, f1, a, print);
}
`

func main() {
	tokenizer := NewTokenizer(source)
	parser := NewParser(tokenizer)
	program := parser.Parse()
	program.Execute()
}
