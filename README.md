# Glox Programming Language

Glox is an enhanced superset of the Lox programming language, introducing a variety of new features designed to improve functionality, ease of use, and flexibility.

## New Features

### 1. Multiline Comments
Multiline comments can be written between `/*` and `*/`, allowing for more flexible commentary within your code.

```lox
/*
This is a multiline comment.
You can add as many lines as you wish.
*/
var bits = 64;
```

### 2. Modulus Operator
The modulus operator % is introduced for performing arithmetic modulus operations.
```lox
var remainder = 10 % 3; // 1
```

### 3. Ternary Operators
Ternary operators offer a concise way to perform if-else operations in a single line.
```lox
var result = condition ? "True" : "False";
```

### 4. Built-in str() Function
The "str()" function converts non-string values into strings, useful for string concatenation.
```lox
var myNumber = 10;
print("My number is " + str(myNumber));
```

### 5. Break and Continue Statements
"break" and "continue" statements control the flow of loops more precisely.
```lox
for (var i = 0; i < 10; i = i + 1) {
    if (i == 5) continue; // Skip the rest of the loop when i is 5.
    if (i == 8) break; // Exit the loop when i is 8.
    print(i);
}
```

### 6. Anonymous Functions
Anonymous functions can be created with "fun() {...}", allowing for functions without names.
```lox
var printer = fun() { print("Hello, Glox!"); };
printer();
```

### 7. Unused Local Variable Warning
Glox warns you when a local variable is declared but not used, aiding in code cleanliness and optimization.
```lox
{
    var f = fun() { return 32; };
}
// [line 2] Warning: Unused variable 'f'.
```

### 8. Static Class Methods
Static methods within classes don't require an instance of the class to be called.
```lox
class Toolbox {
    class sayHello() { print("Hello from Toolbox!"); }
}
Toolbox.sayHello();
```
### 9. Getter Methods
Getter methods allow for computed properties within classes.
```lox
class Circle {
    init(radius) {
        this.radius = radius;
    }

    area { return 3.14 * this.radius * this.radius; }
}
```

### 10. Static Getter Methods
Static getter methods provide class-level computed properties.
```lox
class Constants {
    class PI { return 3.14; }
}
print(Constants.PI); // 3.14
```

### 11. Traits
Glox introduces trait for defining behavior and implements for applying those traits to classes.
```lox
trait Scanner {
    scanned {
        return ["a", "b", "c"];
    }
}

trait Parser {
    parse(tokens) {
        print "Parsing... " + str(tokens) + " Done.";
    }
}

class Interpreter <> Scanner, Parser {
    interpret() {
        this.parse(this.scanned);
        print "Interpreting...";
    }
}
```

You can also inherit classes and implement traits at the same time.
```lox
trait Sin {
    sin90 {
        return 0;
    }
}

trait Cos {
    cos90 {
        return 1;
    }
}

class Math {
    class pi {
        return 3.1415
    }
}

class MyMath < Math <> Sin, Cos {

}
```

### 12. Arrays
Glox adds support for arrays to store a sequence of values.
```lox
var byte = [0, 0, 0, 0, 0, 0, 0, 0];
byte[0] = 1;
byte[4] = 1;
byte[7] = 1;
print(byte); //prints [1, 0, 0, 0, 1, 0, 0, 1]
```# Glox Programming Language

Glox is an enhanced superset of the Lox programming language, introducing a variety of new features designed to improve functionality, ease of use, and flexibility.

## New Features

### 1. Multiline Comments
Multiline comments can be written between `/*` and `*/`, allowing for more flexible commentary within your code.

```lox
/*
This is a multiline comment.
You can add as many lines as you wish.
*/
var bits = 64;
```

### 2. Modulus Operator
The modulus operator % is introduced for performing arithmetic modulus operations.
```lox
var remainder = 10 % 3; // 1
```

### 3. Ternary Operators
Ternary operators offer a concise way to perform if-else operations in a single line.
```lox
var result = condition ? "True" : "False";
```

### 4. Built-in str() Function
The "str()" function converts non-string values into strings, useful for string concatenation.
```lox
var myNumber = 10;
print("My number is " + str(myNumber));
```

### 5. Break and Continue Statements
"break" and "continue" statements control the flow of loops more precisely.
```lox
for (var i = 0; i < 10; i = i + 1) {
    if (i == 5) continue; // Skip the rest of the loop when i is 5.
    if (i == 8) break; // Exit the loop when i is 8.
    print(i);
}
```

### 6. Anonymous Functions
Anonymous functions can be created with "fun() {...}", allowing for functions without names.
```lox
var printer = fun() { print("Hello, Glox!"); };
printer();
```

### 7. Unused Local Variable Warning
Glox warns you when a local variable is declared but not used, aiding in code cleanliness and optimization.
```lox
{
    var f = fun() { return 32; };
}
// [line 2] Warning: Unused variable 'f'.
```

### 8. Static Class Methods
Static methods within classes don't require an instance of the class to be called.
```lox
class Toolbox {
    class sayHello() { print("Hello from Toolbox!"); }
}
Toolbox.sayHello();
```
### 9. Getter Methods
Getter methods allow for computed properties within classes.
```lox
class Circle {
    init(radius) {
        this.radius = radius;
    }

    area { return 3.14 * this.radius * this.radius; }
}
```

### 10. Static Getter Methods
Static getter methods provide class-level computed properties.
```lox
class Constants {
    class PI { return 3.14; }
}
print(Constants.PI); // 3.14
```

### 11. Traits
Glox introduces trait for defining behavior and implements for applying those traits to classes.
```lox
trait Scanner {
    scanned {
        return ["a", "b", "c"];
    }
}

trait Parser {
    parse(tokens) {
        print "Parsing... " + str(tokens) + " Done.";
    }
}

class Interpreter <> Scanner, Parser {
    interpret() {
        this.parse(this.scanned);
        print "Interpreting...";
    }
}
```

You can also inherit classes and also implement traits at the same time.
```lox
trait Sin {
    sin90 {
        return 0;
    }
}

trait Cos {
    cos90 {
        return 1;
    }
}

class Math {
    class pi {
        return 3.1415
    }
}

class MyMath < Math <> Sin, Cos {

}
```

### 12. Arrays
Glox adds support for arrays to store a sequence of values.
```lox
var byte = [0, 0, 0, 0, 0, 0, 0, 0];
byte[0] = 1;
byte[4] = 1;
byte[7] = 1;
print(byte); //prints [1, 0, 0, 0, 1, 0, 0, 1]
```

Built-in "append" function is also provided to append new values to an existing array.
```lox
var primes = [];
primes = append(primes, 2);
primes = append(primes, 3);
primes = append(primes, 5);
print(primes); //prints [2, 3, 5]
```

### 13. REPL Support
Glox enhances the development experience by introducing a REPL environment, allowing for interactive coding sessions. This feature enables you to write and test Glox code in real-time.

To start the REPL, simply run:
```bash
glox
```
Once in the REPL, you can type Glox code and see the output or result immediately.
```lox
> 10 % 3
1
> var name = "Glox";
> print("Hello, " + name);
Hello, Glox
```

Built-in "append" function is also provided to append new values to an existing array.
```lox
var primes = [];
primes = append(primes, 2);
primes = append(primes, 3);
primes = append(primes, 5);
print(primes); //prints [2, 3, 5]
```

### 13. REPL Support
Glox enhances the development experience by introducing a REPL environment, allowing for interactive coding sessions. This feature enables you to write and test Glox code in real-time.

To start the REPL, simply run:
```bash
glox
```
Once in the REPL, you can type Glox code and see the output or result immediately.
```lox
> 10 % 3
1
> var name = "Glox";
> print("Hello, " + name);
Hello, Glox
```

## Building
Glox only requires `go >= 1.22` and does not require third-party dependencies, so building it should be a breeze. In the root directory of a project:
```bash
go build -o glox
```
