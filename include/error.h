#pragma once

#include <cstdarg>

#include <string>
#include <exception>

namespace taolang {

extern char _err_buf[];
extern const int _err_buf_size;

#define __fmt_args(prefix) \
    do {                          \
        va_list va;               \
        va_start(va, format);     \
        std::vsnprintf(&_err_buf[0], _err_buf_size, format, va); \
        va_end(va);     \
        _what = std::string(prefix) + _err_buf; \
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
        __fmt_args("SyntaxError: ");
    }
};

struct NameError : public Error {
    NameError(const char* format, ...) {
        __fmt_args("NameError: ");
    }
};

struct TypeError : public Error {
    TypeError(const char* format, ...) {
        __fmt_args("TypeError: ");
    }
};

struct NotCallableError : public Error {
    NotCallableError(const char* format, ...) {
        __fmt_args("NotCallableError: ");
    }
};

struct NotIndexableError : public Error {
    NotIndexableError(const char* format, ...) {
        __fmt_args("NotIndexableError: ");
    }
};

struct NotAssignableError : public Error {
    NotAssignableError(const char* format, ...) {
        __fmt_args("NotAssignableError: ");
    }
};

struct RangeError : public Error {
    RangeError(const char* format, ...) {
        __fmt_args("RangeError: ");
    }
};

}
