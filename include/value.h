#pragma once

#include <string>
#include <map>
#include <exception>

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

struct String {
    char* p;
    int n;
};

extern std::map<ValueType::Value, std::string> typeNames;

class FunctionExpression;
class Builtin;
class Constructor;
class KeyGetter;
class Values;

struct IExpression {
    virtual Value* Evaluate(Context* ctx) = 0;
    virtual void Assign(Context* ctx, Value* value) = 0;
};

struct ICallable {
    virtual void Execute(Context* ctx, Values* args) = 0;
};

struct IString {
    virtual void String() = 0;
};

class Value : public IExpression {
public:
    ValueType::Value type;
    union {
        bool b;
        int64_t i;
        String str;
        String var;
        KeyGetter* obj;
        EvaluatedFunctionExpression* func;
        Builtin* bi;
        Constructor* ctor;
    };

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
    static Value* fromString(const char* s) {
        auto v = new Value();
        v->type = ValueType::String;
        // TODO
        return v;
    }
    static Value* fromVariable(const char* s) {
        auto v = new Value();
        v->type = ValueType::Variable;
        // TODO
        return v;
    }
    static Value* fromObject(KeyGetter* getter) {
        auto v = new Value();
        v->type = ValueType::Object;
        v->obj = getter;
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
            throw Exception("wrong use");
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
    String string() {
        checkType(ValueType::String);
        return str;
    }
    String variable() {
        checkType(ValueType::Variable);
        return str;
    }
    KeyGetter* object() {
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

}