#pragma once

#include <cctype>
#include <cstring>

#include <string>
#include <memory>
#include <exception>

namespace taolang {
    namespace tokenizer {
        class token_error : public std::exception {
        public:
            token_error(const char* msg)
                : std::exception(msg)
            {

            }
        };

        enum class type_t {
            error,
            plus,
            minus,
            mul,
            div,
            mod,
            eof,
            left_parenthesis,
            right_parenthesis,
            number,
            fn,
        };

        struct token_t {
            type_t          type;
            double          value;
            char            symbol;
            std::string     fn;

            token_t()
                : type(type_t::error)
                , value(0)
                , symbol(0)
            {}
        };

        class tokenizer_t {
        protected:
            const char*     _t;
            const char*     _p;
            token_t         _tk;
            bool            _reuse;
        public:
            tokenizer_t()
                : _t(nullptr)
                , _p(nullptr)
                , _reuse(false)
            {

            }

        public:
            void feed(const char* t) {
                _t = t;
                _p = _t;
            }

            token_t& cur() {
                return _tk;
            }

            void reuse() {
                _reuse = true;
            }

            token_t& next() {
                if(_reuse) {
                    _reuse = false;
                    return _tk;
                }

                _skip_ws();

                _tk.type = type_t::error;

                if(_input() == '\0') {
                    _tk.type = type_t::eof;
                    goto _exit;
                }
                else if(::isdigit(_input())) {
                    _tk.type = type_t::number;
                    _tk.value = _number();
                    goto _exit;
                } 
                else if (::isalpha(_input())) {
                    _tk.type = type_t::fn;
                    _tk.fn = _ident();
                    goto _exit;
                }

                bool is_symbol = true;
                switch(_input()) {
                case '+':
                    _tk.type = type_t::plus;
                    break;
                case '-':
                    _tk.type = type_t::minus;
                    break;
                case '*':
                    _tk.type = type_t::mul;
                    break;
                case '/':
                    _tk.type = type_t::div;
                    break;
                case '(':
                    _tk.type = type_t::left_parenthesis;
                    break;
                case ')':
                    _tk.type = type_t::right_parenthesis;
                    break;
                case '%':
                    _tk.type = type_t::mod;
                    break;
                default:
                    is_symbol = false;
                }

                if(is_symbol) {
                    _char_next();
                    goto _exit;
                }

                throw token_error("unexpected token.");

            _exit:
                return _tk;
            }

            void expect(type_t type) {
                if(next().type == type)
                    return;

                throw token_error("unexpected expectation.");
            }

        protected: // aux
            void _char_next() {
                ++_p;
            }

            int _input() {
                return *_p;
            }

            void _skip_ws() {
                while(::isspace(_input()))
                    _char_next();
            }

            double _number() {
                //_skip_ws();

                auto p = _p;

                while(::isdigit(_input()))
                    _char_next();
                if(_input() == '.')
                    _char_next();
                while(::isdigit(_input()))
                    _char_next();

                std::string buf(p, _p - p);

                return ::atof(buf.c_str());
            }

            std::string _ident() {
                auto p = _p;
                while (::isalpha(_input()))
                    _char_next();

                return std::string(p, _p - p);
            }
        };
    }
}
