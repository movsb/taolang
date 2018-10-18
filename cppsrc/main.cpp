#include <iostream>

#include "tokenizer.h"
#include "parser.h"
#include "program.h"

int main(int argc, char* argv[]) {
	try{
		auto fp = std::fopen(argv[1], "rb");
		auto tkz = taolang::Tokenizer(fp);
        auto parser = taolang::Parser(&tkz);
        auto program = parser.Parse();
        program->Execute();
	} catch(const std::exception& e) {
		std::cerr << e.what() << std::endl;
	}
}
