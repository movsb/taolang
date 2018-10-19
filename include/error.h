#pragma once

#include <cstdarg>

#include <string>
#include <exception>

namespace taolang {

#define __fmt_args() \
    do {                        \
        static char _buf[1024]; \
        va_list va;             \
        va_start(va, format);   \
        std::vsnprintf(&_buf[0], sizeof(_buf)/sizeof(_buf[0]), format, va); \
        va_end(va);     \
        _what = _buf;   \
    } while((0))

struct Error : public std::exception {
    Error(){}
    Error(const std::string& what)
        : _what(what)
    {}
    std::string _what;
    virtual const char* what() const throw() override {
        return _what.c_str();
    }
};

struct SyntaxError : public Error {
    SyntaxError(const char* format, ...) {
        __fmt_args();
    }
};

struct NameError : public Error {
    NameError(const char* format, ...) {
        __fmt_args();
    }
};

struct TypeError : public Error {
    TypeError(const char* format, ...) {
        __fmt_args();
    }
};

struct NotCallableError : public Error {
    NotCallableError() {

    }
};

struct NotIndexableError : public Error {
    NotIndexableError() {

    }
};

struct NotAssignableError : public Error {
    NotAssignableError() {

    }
};

struct RangeError : public Error {
    RangeError(const std::string& e) {
        _what = e;
    }
};

struct KeyTypeError : public Error {
    KeyTypeError() {

    }
};

}
