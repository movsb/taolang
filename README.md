# taolang

A Javascript-like dynamic language.

## Grammar

## Syntax

### 变量定义

```js
let a;          // 变量定义，值为 nil
let b = nil;    // nil 是空类型
let c = true;   // 布尔变量
let d = 123;    // 数值类型，内部类型为 int
let e = "str";  // 字符串类型（原始字符串）
let f = function(x,y,z) {return x+y*z;};        // 函数类型，可以直接当表达式使用
let g = function() {return "test"+"str";}();    // 函数作为表达式，定义后可以直接调用
```

### 函数定义

全局函数或具名函数：

```js
function name(x,y,z) {

}
```

函数作为表达式赋值：

```js
function name() {
    let f = function() {

    };
}
```

函数作为回调函数：

```js
function sync(callback,a,b,c) {
    return callback(a,b,c);
}

function main() {
    print(
        sync(
            function(x,y,z) {
                return x+y*z;
            },
            2,3,4
        )
    );
}
```

### 代码块 & 作用域

仅支持词法作用域，不支持 Javascript 的函数作用域。

```js
function main() {
    let a = 1;
    {
        let a = 2;
        println("inner a: ", a); // 2
        // return a;
    }
    println("outer a: ", a); // 1
}
```
