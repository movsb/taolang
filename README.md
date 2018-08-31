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
let h = {       // 定义一个对象
    a: nil,
    b: 1,
    c: "cc",
    d: true,
    e: function() {print("e");},
    f: {
        xxx: "this is xxx",
    },
    g: function () {
        return "d";
    },
    h: function() {
        return {
            what: "what",
            when: "when",
            "who?": "who?",
        };
    },
};
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

### 闭包

```js
function createCounter() {
    let a = 0;
    return function() {
        a = a + 1;
        return a;
    };
}

function main() {
    let c1 = createCounter();
    let c2 = createCounter();
    println(c1()); // 1
    println(c1()); // 2
    println(c2()); // 1
    println(c2()); // 2
}
```

### 对象定义与访问

```js
function main() {
    let x = {
        a: nil,
        b: 1,
        c: "cc",
        d: true,
        e: function() {print("e");},
        f: {
            xxx: "this is xxx",
        },
        g: function () {
            return "d";
        },
        h: function() {
            return {
                what: "what",
                when: "when",
                "who?": "who?",
            };
        },
    };
    println(x);
    println(x.a);
    println(x.b);
    println(x["c"]);
    println(x[x.g()]);
    println(x.e);
    println(x.f.xxx);
    println(x.h().what);
    println(x.h()["who?"]);
    println(x.y);
}
```

### 数组定义与访问

```js
function main() {
    let a = [1,true,nil,"str",{},1+2*3,];
    println(a);
    println(a.length);
    println(a[3]);
    println([9][0]);
}

function main2() {
    let a = [1,2,3,4,5,6,7,8];
    let n = a.length-1;
    while n >= 0 {
        println(a[n]);
        n = n - 1;
    }
}
```

### Lambda表达式与函数式编程

```js
function each() {
    let a = [1,3,5,7,9];
    a.each(e=>println(e));
}

function map() {
    let a = [1,3,5,7,9];
    let b = a.map(x=>x*x);
    println(b);
}

function reduce() {
    let a = [1,2,3];
    let b = a.reduce((memo,num)=>memo+num,0);
    println(b);
}

function find() {
    let a = [1,3,5,7,9];
    let b = a.find(x=>x>5);
    println(b);
}

function filter() {
    let a = [1,2,3,4,5,6,7,8,9];
    let b = a.filter(x=>x%2==0);
    println(b);
}

function where() {
    let a = [
        {a:1,b:3},
        {a:2,b:2},
        {a:3,b:1},
    ];
    let b = a.where(x=>x.a==1||x.b==1);
    println(b);
}
```

### 控制语句

#### while 控制语句

表达式部分不用 `()` 括起来。

```js
function While() {
    let n = 10;
    while n > 0 {
        print(n);
        n = n - 1;
    }
}

function Break() {
    let n = 10;
    while n > 0 {
        print(n);
        break;
    }
}

function Return() {
    let n = 10;
    while n > 0 {
        print(n);
        return nil;
    }
}

function main() {
    While();
    Break();
    Return();
}
```

### if-else 控制语句

```js
function If() {
    if 1 > 0 {
        println("1 > 0");
    }

    if 1 > 2 {
        println("1 > 2");
    } else {
        println("else 1 > 2");
    }

    if 1 > 2 {
        println("1 > 2");
    } else if 2 > 3 {
        println("2 > 3");
    } else {
        println("else");
    }
}

function Break() {
    let a = 10;
    while a > 0 {
        if a == 5 {
            break;
        }
        a= a-5;
    }
    print(a);
}

function Return() {
    let a = 10;
    while a > 0 {
        if a == 8 {
            return a;
        }
        a = a- 1;
    }
}

function main() {
    If();
    Break();
    println(Return());
}
```
