#include <iostream>
#include <sstream>
#include <algorithm>

#include "object.h"
#include "error.h"

namespace taolang {

Value* Object::GetKey(const std::string& key) {
    return _props[key];
}

void Object::SetKey(const std::string& key, Value* val) {
    _props[key] = val;
}

std::string Object::ToString() {
    std::ostringstream oss;
    auto keys = _sortedKeys();
    auto i=0, n = (int)keys.size();
    oss << "{";
    for(const auto& key : keys) {
        oss << key;
        oss << _props[key]->ToString();
        if(i != n-1) {
            oss << ",";
        }
        i++;
    }
    oss << "}";
    return std::move(oss.str());
}

std::vector<std::string> Object::_sortedKeys() {
    std::vector<std::string> keys;
    for(auto& prop : _props) {
        keys.emplace_back(prop.first);
    }
    std::sort(keys.begin(), keys.end());
    return std::move(keys);
}

Value* Array::GetElem(int index) {
    _checkIndex(index);
    return _elems[index];
}

void Array::SetElem(int index, Value* value) {
    _checkIndex(index);
    _elems[index] = value;
}

void Array::_checkIndex(int index) {
    if(index < 0 || index > Len()) {
        throw RangeError("array index out of range");
    }
}

std::string Array::ToString() {
    std::ostringstream oss;
    auto n = Len();
    oss << "[";
    for(auto i=0; i<n; i++) {
        oss << _elems[i]->ToString();
        if(i != n-1) {
            oss << ",";
        }
    }
    oss << "}";
    return std::move(oss.str());
}

Value*  Global::println(Context* ctx, Values* args) {
    std::ostringstream oss;
    auto n = args->Size();
    for(auto i = 0; i < n; i++) {
        oss << args->Get(i)->ToString();
        if(i != n-1) {
            oss << " ";
        }
    }
    oss << "\n";
    auto str = oss.str();
    std::cout << str;
    return Value::fromNil();
}

}
