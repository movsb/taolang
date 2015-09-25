#pragma once

#include "ast.h"

namespace taolang {
    namespace eval {
        class eval_error : public std::exception {
        public:
            eval_error(const char* msg)
                : std::exception(msg)
            {

            }
        };

        using namespace ast;

        class evaluator_t {
        public:
            double eval(ast_node* node) {
                return _eval(node);
            }

        protected:
            double _eval(ast_node* node) {
                if(node->type == ast_type::value) {
                    return node->value;
                }
                else if(node->type == ast_type::unary_minus) {
                    return -_eval(node->left);
                }
                else if (node->type == ast_type::call) {
                    auto& fn = node->left->fn;
                    auto args = _eval(node->right);
                    std::cout << "调用函数: " << fn << ", 参数: " << args << std::endl;
                    return 0;
                }
                else {
                    auto v1 = _eval(node->left);
                    auto v2 = _eval(node->right);

                    switch(node->type)
                    {
                    case ast_type::binary_plus:
                        return v1 + v2;
                    case ast_type::binary_minus:
                        return v1 - v2;
                    case ast_type::binary_mul:
                        return v1 * v2;
                    case ast_type::binary_div:
                        return v1 / v2;
                    case ast_type::binary_mod:
                        return int(v1) % int(v2);
                    }
                }

                throw eval_error("error while eval-ing.");
            }
        };
    }
}
