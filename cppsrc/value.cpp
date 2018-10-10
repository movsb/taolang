#include "value.h"

namespace taolang {

std::map<ValueType::Value, std::string> typeNames = {
    {ValueType::Nil, "nil"},
    {ValueType::Boolean, "boolean"},
    {ValueType::Number, "number"},
    {ValueType::String, "string"},
    {ValueType::Variable, "variable"},
    {ValueType::Object, "object"},
    {ValueType::Function, "function"},
    {ValueType::Builtin, "builtin"},
    {ValueType::Class, "class"},
};

}
