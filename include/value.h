#pragma once

#include <string>
#include <map>

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

class Value {
public:
    ValueType::Value type;
    union {
        bool b;
        int64_t i;
        String str;
        String var;
        KeyGetter* obj;
        FunctionExpression* func;
        Builtin* bi;
        Constructor* ctor;
    };

public:
    static Value* FromNil() {
        auto v = new Value();
        v->type = ValueType::Nil;
        return v;
    }
    static Value* FromBoolean(bool b) {
        auto v = new Value();
        v->type = ValueType::Boolean;
        v->b = b;
        return v;
    }
    static Value* FromNumber(int64_t i) {
        auto v = new Value();
        v->type = ValueType::Number;
        v->i = i;
        return v;
    }
    static Value* FromString(const char* s) {
        auto v = new Value();
        v->type = ValueType::String;
        // TODO
        return v;
    }
    static Value* FromVariable(const char* s) {
        auto v = new Value();
        v->type = ValueType::Variable;
        // TODO
        return v;
    }
    static Value* FromObject(KeyGetter* getter) {
        auto v = new Value();
        v->type = ValueType::Object;
        v->obj = getter;
        return v;
    }
};

}
