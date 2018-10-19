#include "program.h"
#include "context.h"
#include "expression.h"

namespace taolang {

void Program::Execute() {
    auto globalContext = new Context(nullptr);
    for(auto stmt : _stmts) {
        stmt->Execute(globalContext);
    }
    auto main = Value::fromVariable("main");
    CallFunc(globalContext, main, nullptr);
}

}
