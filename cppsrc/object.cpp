#include "object.h"

namespace taolang {

Value* Object::GetKey(const std::string& key) {
    return _props[key];
}

void Object::SetKey(const std::string& key, Value* val) {
    _props[key] = val;
}

}
