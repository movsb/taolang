#pragma once

#include <vector>
#include <unordered_map>
#include <string>

#include "value.h"

namespace taolang {

template<typename T>
BuiltinFunction get_mfn(T t) {
    union {
        T t;
        BuiltinFunction f;
    } u;
    u.t = t;
    return u.f;
}

class Object : public IObject {
public:
    Object()
        : _array(false)
    {}

public:
    virtual Value* GetKey(const std::string& key) override;
    virtual void SetKey(const std::string& key, Value* val) override;

protected:
    std::vector<Value*> _elems;
    std::unordered_map<std::string, Value*> _props;
    bool _array;
};

class Global : public  Object {
public:
    Global() {
        _props["println"] = Value::fromBuiltin(this, "println", get_mfn(&Global::println));
    }

private:
    Value* __attribute__((stdcall)) println(Context* ctx, Values* args) {

    }
};


}
