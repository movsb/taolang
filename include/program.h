#pragma once

#include <vector>
#include "value.h"

namespace taolang {

class Program {
public:
    void Execute();
    std::vector<IStatement*> _stmts;
};

}