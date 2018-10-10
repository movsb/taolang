#include <iostream>

#include "tokenizer.h"

int main(int argc, char* argv[]) {
	try{
		auto fp = std::fopen(argv[1], "rb");
		auto tkz = taolang::Tokenizer(fp);
		while((true)) {
			auto tk = tkz.Next();
			std::cout << tk.string() << std::endl;
			if(tk.type == taolang::TokenType::ttEOF) {
				break;
			}
		}
	} catch(const std::exception& e) {
		std::cerr << e.what() << std::endl;
	}
}
