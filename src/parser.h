#pragma once

#include "ast.h"
#include "tokenizer.h"

/* Àà BNF ·¶Ê½
£¨º¬×óµİ¹é£©
expr
    : expr + term
    | expr - term
    | term
    ;

term
    : term * factor
    | term / factor
    | term % factor
    |: factor
    ;

factor
    : (exp)
    | - exp
    | number
    | fn ( exp )
    ;

£¨²»º¬×óµİ¹é£©
expr
    : term expr1
    ;

expr1
    : + term expr1
    | - term expr1
    | epsilon
    ;

term
    : factor term1
    ;

term1
    : * factor term1
    | / factor term1
    | % factor term1
    | epsilon
    ;

factor
    : ( exp )
    | - exp
    | number
    | fn ( exp )
    ;
*/

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

                if(e1node)
                    return new_ast_node(ast_type::binary_plus, tnode, e1node);
                else
                    return tnode;
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
                    if(e1node)
                        root = new_ast_node(ast_type::binary_plus, e1node, tnode);
                    else
                        root = tnode;
                    break;
                case tokenizer::type_t::minus:
                    tnode = term();
                    e1node = expr1();
                    if(e1node)
                        root = new_ast_node(ast_type::binary_minus, e1node, tnode);
                    else
                        root = new_ast_node(ast_type::unary_minus, tnode);
                    break;
                }

                if(root != nullptr)
                    return root;

                if(tk.type == tokenizer::type_t::eof) {
                    _tkr.reuse();
                    return nullptr;
                }

                _tkr.reuse();
                return nullptr;
                return new_ast_node(0);
            }

            ast_node* term() {
                auto fnode = factor();
                auto t1node = term1();

                if(t1node)
                    return new_ast_node(ast_type::binary_mul, fnode, t1node);
                else
                    return fnode;
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
                    if(t1node)
                        root = new_ast_node(ast_type::binary_mul, t1node, fnode);
                    else
                        root = fnode;
                    break;
                case tokenizer::type_t::div:
                    fnode = factor();
                    t1node = term1();
                    if(t1node)
                        root = new_ast_node(ast_type::binary_div, t1node, fnode);
                    else
                        root = fnode;
                    break;
                case tokenizer::type_t::mod:
                    fnode = factor();
                    t1node = term1();
                    if(t1node)
                        root = new_ast_node(ast_type::binary_mod, t1node, fnode);
                    else
                        root = fnode;
                    break;
                }

                if(root != nullptr)
                    return root;

                if(tk.type == tokenizer::type_t::eof) {
                    _tkr.reuse();
                    return nullptr;
                }

                _tkr.reuse();
                return nullptr;
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
                case tokenizer::type_t::fn:
                    node = new_ast_node(tk.fn);
                    _tkr.expect(tokenizer::type_t::left_parenthesis);
                    node = new_ast_node(ast_type::call, node, factor());
                    _tkr.expect(tokenizer::type_t::right_parenthesis);
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

            ast_node* new_ast_node(const std::string& fn) {
                ast_node* node = new ast_node;
                node->type = ast_type::call;
                node->fn = fn;
                return node;
            }
        };
    }
}
