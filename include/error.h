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

struct Exception : public std::exception {
    Exception(){}
    Exception(const std::string& what)
        : _what(what)
    {}
    std::string _what;
    virtual const char* what() const throw() override {
        return _what.c_str();
    }
};

struct SyntaxError : public Exception {
    SyntaxError(const char* format, ...) {
        __fmt_args();
    }
};

struct NameError : public Exception {
    NameError(const char* format, ...) {
        __fmt_args();
    }
};

struct TypeError : public Exception {
    TypeError(const char* format, ...) {
        __fmt_args();
    }
};

struct NotCallableError : public Exception {
    NotCallableError() {

    }
};

struct NotIndexableError : public Exception {
    NotIndexableError() {

    }
};

struct NotAssignableError : public Exception {
    NotAssignableError() {

    }
};

struct RangeError : public Exception {
    RangeError() {

    }
};

struct KeyTypeError : public Exception {
    KeyTypeError() {

    }
};

}
