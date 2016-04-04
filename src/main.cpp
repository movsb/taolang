#include <iostream>
#include <string>

#include "ast.h"
#include "tokenizer.h"
#include "parser.h"
#include "eval.h"

bool run(const char* syntax) {
    taolang::parser::parser_t   parser;

    try {
        auto ast = parser.parse(syntax);
        try {
            taolang::eval::evaluator_t  evaler;
            auto val = evaler.eval(ast);

            std::cout << syntax << " = " << val << std::endl;
        }
        catch(taolang::eval::eval_error& e) {
            std::cout << "eval_error: " << e.what() << std::endl;
        }
    }
    catch(taolang::parser::parser_error& e) {
        std::cout << "parser_error: " << e.what() << std::endl;
    }

    return true;
}

int main()
{
    auto test = []() {
        run("1+2+3+4");
        run("1*2*3*4");
        run("1-2-3-4");
        run("1/2/3/4");
        run("1*2+3*4");
        run("1+2*3+4");
        run("(1+2)*(3+4)");
        run("1+(2*3)*(4+5)");
        run("1+(2*3)/4+5");
        run("5/(4+3)/2");
        run("1 + 2.5");
        run("125");
        run("-1");
        run("-1+(-2)");
        run("-1+(-2.0)");
        run("   1*2,5");
        run("   1*2.5e2");
        run("M1 + 2.5");
        run("1 + 2&5");
        run("1 * 2.5.6");
        run("2 * 2.5");
        run("1 ** 2.5");
        run("*1 / 2.5");
    };

    test();

    std::string program;
    do {
        std::cout << "input program: ";
        std::getline(std::cin, program, '\n');
    } while(run(program.c_str()));

    return 0;
}
