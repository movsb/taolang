#pragma once

#include <string>

namespace taolang {
    namespace ast {
        enum class ast_type {
            undefined,
            binary_plus,
            binary_minus,
            binary_mul,
            binary_div,
            unary_plus,
            unary_minus,
            value,
        };

        class ast_node {
        public:
            ast_type    type;
            double      value;
            ast_node*   left;
            ast_node*   right;

        public:
            ast_node()
                : type(ast_type::undefined)
                , value(0)
                , left(nullptr)
                , right(nullptr)
            {}

            ~ast_node() {
                delete left;
                delete right;
            }
        };
    }
}
