#pragma once

#include <vector>
#include <unordered_map>
#include <string>
#include <iostream>

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
    Object() {}

public:
    virtual Value* GetKey(const std::string& key) override;
    virtual void SetKey(const std::string& key, Value* val) override;

protected:
    std::unordered_map<std::string, Value*> _props;
};

class Array : public Object, public IArray {
public:
    virtual int Len() override { return (int)_elems.size(); }
    virtual Value* GetElem(int index) override;
    virtual void SetElem(int index, Value* value) override;

private:
    void _checkIndex(int index);

protected:
    std::vector<Value*> _elems;
};

class Global : public  Object {
public:
    Global() {
        _props["println"] = Value::fromBuiltin(this, "println", get_mfn(&Global::println));
    }

private:
    Value*  println(Context* ctx, Values* args) {
        std::cout << "println called";
        return Value::fromNil();
    }
};


}
