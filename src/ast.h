#pragma once

#include <string>

namespace taolang {
    namespace ast {
        enum class ast_type {
            undefined,
            only_node,
            binary_plus,
            binary_minus,
            binary_mul,
            binary_div,
            binary_mod,
            unary_plus,
            unary_minus,
            value,
            call,
        };

        class ast_node {
        public:
            ast_type    type;
            double      value;
            std::string fn;
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
