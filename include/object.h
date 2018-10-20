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

#define _AddProp(name, func) \
    _props[name] = Value::fromBuiltin(this, name, get_mfn(func))

class Object : public IObject {
public:
    Object(const std::string& typeName)
        : _typeName(typeName)
        {}
    static Object* New() {
        return new Object();
    }
private:
    Object() : Object("Object") {}

public:
    virtual std::string TypeName() override {return _typeName;}
    virtual Value* GetKey(const std::string& key) override;
    virtual void SetKey(const std::string& key, Value* val) override;
    virtual std::string ToString() override;

private:
    std::vector<std::string> _sortedKeys();

protected:
    std::string _typeName;
    std::unordered_map<std::string, Value*> _props;
};

class Array : public Object, public IArray {
public:
    Array() : Object("Array") {}

public:
    virtual int Len() override { return (int)_elems.size(); }
    virtual Value* GetElem(int index) override;
    virtual void SetElem(int index, Value* value) override;
    virtual std::string ToString() override;

private:
    void _checkIndex(int index);

protected:
    std::vector<Value*> _elems;
};

class Global : public  Object {
public:
    Global() : Object("Global") {
        _AddProp("println", &Global::println);
    }

private:
    Value*  println(Context* ctx, Values* args);
};


}
