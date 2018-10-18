#include "context.h"
#include "value.h"

namespace taolang {

Value* Context::FindSymbol(const std::string& name, bool outer) {
    for(auto& symbol : _symbols) {
        if(symbol.name == name) {
            return symbol.value;
        }
    }
    if(outer) {
        if(_parent != nullptr) {
            return _parent->FindSymbol(name, outer);
        }
        return FromGlobal(name);
    }
    return nullptr;
}

Value* Context::MustFind(const std::string& name, bool outer) {
    auto val = FindSymbol(name, outer);
    if(val == nullptr) {
        throw NameError("name %s not defined", name.c_str());
    }
    return val;
}

Value* Context::FromGlobal(const std::string& name) {
    auto global = MustFind("global", true);
    if(!global->isObject()) {
        throw TypeError("global is not an object");
    }
    auto obj = global->object();
}

void Context::AddSymbol(const std::string& name, Value* value) {
    if(FindSymbol(name, false) != nullptr) {
        throw NameError("name `%s' redefined");
    }
    _symbols.push_back({name, value});
}

void Context::SetSymbol(const std::string& name, Value* value) {
    for(auto& symbol : _symbols) {
        if(symbol.name == name) {
            symbol.value = value;
            return;
        }
    }
    if(_parent != nullptr) {
        _parent->SetSymbol(name, value);
        return;
    }
    throw NameError("name `%s' not defined", name.c_str());
}

void Context::SetParent(Context* parent) {
    _parent = parent;
}

void Context::SetReturn(Value* value) {
    _hasRet = true;
    _retVal = value;
}

void Context::SetBreak() {
    _broke = true;
}

}
