#pragma once

#include <string>
#include <map>
#include <exception>
#include <vector>

#include "error.h"

namespace taolang {

struct ValueType {
    enum Value {
        Nil,
        Boolean,
        Number,
        String,
        Variable,
        Object,
        Function,
        Builtin,
        Class,
    };
};

extern std::map<ValueType::Value, std::string> typeNames;

class FunctionExpression;
class EvaluatedFunctionExpression;
class Builtin;
class Constructor;
class KeyGetter;
class Value;
class Values;
class Context;

struct IExpression {
    virtual Value* Evaluate(Context* ctx) = 0;
    virtual void Assign(Context* ctx, Value* value) = 0;
};

struct IStatement {
    virtual void Execute(Context* ctx) = 0;
};

struct ICallable {
    virtual Value* Execute(Context* ctx, Values* args) = 0;
};

struct IObject {
    virtual Value* GetKey(const std::string& key) = 0;
    virtual void SetKey(const std::string& key, Value* val) = 0;
};

struct IArray {
    virtual int Len() = 0;
    virtual Value* GetElem(int pos) = 0;
    virtual void SetElem(int pos, Value* val) = 0;
};

struct IString {
    virtual void String() = 0;
};

typedef Value* (*BuiltinFunction)(void* that, Context* ctx, Values* args);

class Builtin : public ICallable {
public:
    void* _that;
    std::string _name;
    BuiltinFunction _func;
public:
    virtual Value* Execute(Context* ctx, Values* args) override;
};

class Value : public IExpression {
public:
    Value() {}
public:
    ValueType::Value type;
    union {
        bool b;
        int64_t i;
        IObject* obj;
        EvaluatedFunctionExpression* func;
        Builtin* bi;
        Constructor* ctor;
    };
    std::string str;
    std::string var;

public:
    static Value* fromNil() {
        auto v = new Value();
        v->type = ValueType::Nil;
        return v;
    }
    static Value* fromBoolean(bool b) {
        auto v = new Value();
        v->type = ValueType::Boolean;
        v->b = b;
        return v;
    }
    static Value* fromNumber(int64_t i) {
        auto v = new Value();
        v->type = ValueType::Number;
        v->i = i;
        return v;
    }
    template<typename T>
    static Value* fromString(const T& s) {
        auto v = new Value();
        v->type = ValueType::String;
        v->str = s;
        return v;
    }
    template<typename T>
    static Value* fromVariable(const T& s) {
        auto v = new Value();
        v->type = ValueType::Variable;
        v->var = s;
        return v;
    }
    static Value* fromObject(IObject* obj) {
        auto v = new Value();
        v->type = ValueType::Object;
        v->obj = obj;
        return v;
    }
    static Value* fromFunction(FunctionExpression* func, Context* closure);
    static Value* fromBuiltin(void* that, const std::string& name, BuiltinFunction func) {
        auto v = new Value();
        v->type = ValueType::Builtin;
        auto bi = new Builtin();
        bi->_that = that;
        bi->_name = name;
        bi->_func = func;
        v->bi = bi;
        return v;
    }

public:
    bool isNil() {
        return type == ValueType::Nil;
    }
    bool isBoolean() {
        return type == ValueType::Boolean;
    }
    bool isNumber() {
        return type == ValueType::Number;
    }
    bool isString() {
        return type == ValueType::String;
    }
    bool isObject() {
        return type == ValueType::Object;
    }
    bool isVariable() {
        return type == ValueType::Variable;
    }
    bool isFunction() {
        return type == ValueType::Function;
    }
    bool isBuiltin() {
        return type == ValueType::Builtin;
    }
    bool isConstructor() {
        return type == ValueType::Class;
    }
    bool isCallable() {
        return callable() != nullptr;
    }
    void checkType(ValueType::Value type) {
        if(this->type != type) {
            throw Error("wrong use");
        }
    }

public:
    bool boolean() {
        checkType(ValueType::Boolean);
        return b;
    }
    int64_t number() {
        checkType(ValueType::Number);
        return i;
    }
    std::string string() {
        checkType(ValueType::String);
        return str;
    }
    std::string variable() {
        checkType(ValueType::Variable);
        return str;
    }
    IObject* object() {
        checkType(ValueType::Object);
        return obj;
    }
    EvaluatedFunctionExpression* function() {
        checkType(ValueType::Function);
        return func;
    }
    Builtin* builtin() {
        checkType(ValueType::Builtin);
        return bi;
    }
    Constructor* constructor() {
        checkType(ValueType::Class);
        return ctor;
    }
    ICallable* callable() {
        // TODO
        throw NotCallableError();
    }

public:
    virtual Value* Evaluate(Context* ctx) override;
    virtual void Assign(Context* ctx, Value* value) override;

public:
    std::string TypeName() {
        return typeNames[type];
    }

    std::string String();
    bool truth(Context* ctx);
};

class Values {
public:
    int Size() {
        return (int)_values.size();
    }
    void Push(Value* value) {
        _values.push_back(value);
    }
    Value* Get(int i) {
        if(i<0 && i > Size() - 1) {
            throw RangeError("Values' index out of range");
        }
        return _values[i];
    }

private:
    std::vector<Value*> _values;
};

}
