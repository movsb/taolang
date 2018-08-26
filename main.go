package main

var source = `
function f1(a,b,c) {
	a = 233;
	print(a,b,c);
}

function main() {
	function f2(a,b,c) {
		return a+b*c;
	}
	var a = 133;
	a = 333;
	f1(123, 456, 789);
	print("string", 123, true, nil, f2(1,2,3), a, print);
}
`

func main() {
	tokenizer := NewTokenizer(source)
	parser := NewParser(tokenizer)
	program := parser.Parse()
	program.Execute()
}
