#pragma once

#include <string>
#include <vector>

#include "value.h"

namespace taolang {

struct Symbol {
    std::string name;
    Value*      value;
};

class Context {
public:
    Context(Context* parent)
        : _parent(parent)
        , _broke(false)
        , _hasRet(false)
        , _retVal(nullptr)
    {}
public:
    Value* FindSymbol(const std::string& name, bool outer);
    Value* MustFind(const std::string& name, bool outer);
    Value* FromGlobal(const std::string& name);
    void AddSymbol(const std::string& name, Value* value);
    void AddObject(const std::string& name, IObject* obj);
    void SetSymbol(const std::string& name, Value* value);
    void SetParent(Context* parent);
    void SetReturn(Value* value);
    void SetBreak();

public:
    Context*            _parent;
    std::vector<Symbol> _symbols;
    bool                _broke;
    bool                _hasRet;
    Value               *_retVal;
};

}
