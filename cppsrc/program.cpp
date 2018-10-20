#include "program.h"
#include "context.h"
#include "expression.h"
#include "object.h"

namespace taolang {

void Program::Execute() {
    auto globalContext = new Context(nullptr);
    auto globalObject = new Global();

    globalContext->AddObject("global", globalObject);

    for(auto stmt : _stmts) {
        stmt->Execute(globalContext);
    }

    auto main = Value::fromVariable("main");
    CallFunc(globalContext, main, nullptr);
}

}
