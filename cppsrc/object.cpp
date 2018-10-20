#include "object.h"
#include "error.h"

namespace taolang {

Value* Object::GetKey(const std::string& key) {
    return _props[key];
}

void Object::SetKey(const std::string& key, Value* val) {
    _props[key] = val;
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

}
