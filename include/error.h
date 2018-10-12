#pragma once

#include <cstdarg>

#include <string>
#include <exception>

namespace taolang {

struct SyntaxError : public std::exception {
    std::string _what;
    SyntaxError(const char* format, ...)
    {
        // TODO must be single thread
        static char _buf[1024];
        va_list va;
        va_start(va, format);
        std::vsnprintf(&_buf[0], sizeof(_buf)/sizeof(_buf[0]), format, va);
        va_end(va);
        _what = _buf;
    }
    virtual const char* what() const throw() override {
        return _what.c_str();
    }
};

}
