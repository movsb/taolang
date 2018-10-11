#include <string>
#include <vector>

#include "value.h"

namespace taolang {

struct Symbol {
    std::string name;
    Value       value;
};

class Context {
public:
    Value* FindSymbol(const std::string& name, bool outer);
    Value* MustFind(const std::string& name, bool outer);
    Value* FromGlobal(const std::string& name);
    void AddSymbol(const std::string& name, Value* value);
    void SetSymbol(const std::string& name, Value* value);
    void SetParent(Context* parent);
    void SetReturn(Value* value);
    void SetBreak();

public:
    Context*            _parent;
    std::vector<Symbol> _symbols;
    bool                _broke;
    bool                _hasRet;
    Value               _retVal;

};

}
