package test

import "testing"

func TestClassDeclaration(t *testing.T) {
	program := `
class Scanner {
  scan() {
    return "scanning...";
  }
}

print Scanner;
`
	testPrograms(t, []testCase{
		{program, "<class Scanner>\n"},
	})
}

func TestClassInstances(t *testing.T) {
	program := `
class Parser {}
var parser = Parser();
print parser;
`
	testPrograms(t, []testCase{
		{program, "<Parser instance>\n"},
	})
}

func TestClassProperties(t *testing.T) {
	program := `
class Parser {
	checkGrammar() {
		return true;
	}
	parse() {
		print "parsing... Done";
	}
}

var parser = Parser();
parser.tokens = "{NUMBER}{PLUS}{NUMBER}";
parser.tokens_eof = parser.tokens + "{EOF}";

print parser.checkGrammar();
parser.parse();
print parser.tokens;
print parser.tokens_eof;

`
	testPrograms(t, []testCase{
		{program, "true\nparsing... Done\n{NUMBER}{PLUS}{NUMBER}\n{NUMBER}{PLUS}{NUMBER}{EOF}\n"},
	})
}

func TestClassStaticProperties(t *testing.T) {
	program := `
class Math {
    class square(n) {
        return n * n;
    }

    class cube(n) {
        return Math.square(n) * n;
    }

	class circleArea(r) {
		return Math.PI * Math.square(r); 
	}
}

Math.PI = 3.14159265359;
Math.superCube = fun(n) {
    return n * Math.cube(n);
};

print Math.square(9);
print Math.cube(9);
print Math.superCube(9);
print Math.circleArea(19);
`
	testPrograms(t, []testCase{
		{program, "81\n729\n6561\n1134.11494794599\n"},
	})
}

func TestClassGetters(t *testing.T) {
	program := `
class Circle {
	init(radius) {
		this.radius = radius;
	}

	square {
		print "Impossible";
		return nil;
	}

	class pi {
		return 3.141592653;
	}

 	area {
		return Circle.pi * this.radius * this.radius;
	}
}

var circle = Circle(4);
print Circle.pi;
print circle.area;
print circle.square;

`
	testPrograms(t, []testCase{
		{program, "3.141592653\n50.265482448\nImpossible\nnil\n"},
	})
}

func TestClassThis(t *testing.T) {
	program1 := `
class Person {
	sayName() {
		print this.name;
  	}

  	saySurname() {
    	print this.surname;
  	}
}

var hank = Person();
hank.name = "Hank";
hank.surname = "Schrader";

var refSurname = hank.saySurname;
hank.sayName();
refSurname();
`

	program2 := `
class Person {
	sayName() {
		print this.name;
  	}

  	saySurname() {
    	print this.surname;
  	}
}

var jane = Person();
jane.name = "Jane";
jane.sayName();

var bill = Person();
bill.name = "Bill";
bill.sayName();

bill.sayName = jane.sayName;
bill.sayName();
`
	testPrograms(t, []testCase{
		{program1, "Hank\nSchrader\n"},
		{program2, "Jane\nBill\nJane\n"},
	})
}

func TestClassInitializers(t *testing.T) {
	program1 := `
class Rectangle {
	init(width, height) {
		this.width = width;
		this.height = height;
	}

	area() {
		return this.width * this.height;
	}
}

var rectangle = Rectangle(7, 8);
var square = Rectangle(9, 9);

print rectangle.area();
print square.area();
`

	program2 := `
class Foo {
	init() {
		return;
		print 1 / 0;
	}
}

class Bar {
  	init(number) {
		this.number = number;
		print this;
	}
}

var foo = Foo();
print foo;

var bar = Bar(11);
var bar2 = bar.init(12);
print bar.number * bar2.number;
`

	testPrograms(t, []testCase{
		{program1, "56\n81\n"},
		{program2, "<Foo instance>\n<Bar instance>\n<Bar instance>\n144\n"},
	})
}

func TestClasses(t *testing.T) {
	program := `
class Car {
    init(make, model, speed, color) {
        this.make = make;
        this.model = model;
        this.speed = speed;
        this.color = color;
    }

    description() {
        return "{" + this.make + ", " + this.model + ", " + str(this.speed) + ", " + this.color + "}";
    }

    drive() {
        print "Vroom Vroooom";
    }
}

var carrera = Car("Porsche", "Carrera GT", 334, "Silver");
print carrera.description();
carrera.drive();

var f40 = Car("Ferrari", "F40", 367, "Red");
print f40.description();
f40.drive();

if (f40.speed > carrera.speed) {
	print "F40 is faster";
}
`
	testPrograms(t, []testCase{
		{program, "{Porsche, Carrera GT, 334, Silver}\nVroom Vroooom\n{Ferrari, F40, 367, Red}\nVroom Vroooom\nF40 is faster\n"},
	})
}
