#pragma once

#include "ast.h"
#include "tokenizer.h"

namespace taolang {
    namespace parser {
        class parser_error : public std::exception {
        public:
            parser_error(const char* msg)
                : std::exception(msg)
            {

            }
        };

        using namespace ast;

        class parser_t {
        protected:
            tokenizer::tokenizer_t  _tkr;

        public:
            ast_node* parse(const char* syntax) {
                _tkr.feed(syntax);

                try {
                    return expr();
                }
                catch(tokenizer::token_error& e) {
                    throw parser_error((std::string("invalid token: ") + e.what()).c_str());
                }
            }

        protected:
            ast_node* expr() {
                auto tnode = term();
                auto e1node = expr1();

                return new_ast_node(ast_type::binary_plus, tnode, e1node);
            }

            ast_node* expr1() {
                ast_node* root = nullptr;
                ast_node* tnode;
                ast_node* e1node;

                auto& tk = _tkr.next();
                switch(tk.type)
                {
                case tokenizer::type_t::plus:
                    tnode = term();
                    e1node = expr1();
                    root = new_ast_node(ast_type::binary_plus, e1node, tnode);
                    break;
                case tokenizer::type_t::minus:
                    tnode = term();
                    e1node = expr1();
                    root = new_ast_node(ast_type::binary_minus, e1node, tnode);
                    break;
                }

                if(root != nullptr)
                    return root;

                _tkr.reuse();
                return new_ast_node(0);
            }

            ast_node* term() {
                auto fnode = factor();
                auto t1node = term1();

                return new_ast_node(ast_type::binary_mul, fnode, t1node);
            }

            ast_node* term1() {
                ast_node* root = nullptr;
                ast_node* fnode;
                ast_node* t1node;

                auto& tk = _tkr.next();
                switch(tk.type)
                {
                case tokenizer::type_t::mul:
                    fnode = factor();
                    t1node = term1();
                    root = new_ast_node(ast_type::binary_mul, t1node, fnode);
                    break;
                case tokenizer::type_t::div:
                    fnode = factor();
                    t1node = term1();
                    root = new_ast_node(ast_type::binary_div, t1node, fnode);
                    break;
                }

                if(root != nullptr)
                    return root;

                _tkr.reuse();
                return new_ast_node(1);
            }

            ast_node* factor() {
                ast_node* node = nullptr;
                auto& tk = _tkr.next();

                switch(tk.type)
                {
                case tokenizer::type_t::left_parenthesis:
                    node = expr();
                    _tkr.expect(tokenizer::type_t::right_parenthesis);
                    break;
                case tokenizer::type_t::minus:
                    node = factor();
                    node = new_ast_node(ast_type::unary_minus, node);
                    break;
                case tokenizer::type_t::number:
                    node = new_ast_node(tk.value);
                    break;
                }

                if(node) {
                    return node;
                }

                throw parser_error("syntax error.");
            }
        protected:
            ast_node* new_ast_node(ast_type type, ast_node* left, ast_node* right = nullptr) {
                ast_node* node = new ast_node;
                node->type  = type;
                node->left  = left;
                node->right = right;
                return node;
            }

            ast_node* new_ast_node(double value) {
                ast_node* node = new ast_node;
                node->type  = ast_type::value;
                node->value = value;
                return node;
            }
        };
    }
}
