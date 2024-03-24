package test

import "testing"

func TestTraits(t *testing.T) {
	program1 := `
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

print Scanner;
print Parser;
Interpreter().interpret();
`

	program2 := `
trait Sin {
	class sin30 {
		return 0.50;
	}

	class sin45 {
		return 0.07;
	}

	class sin90 {
		return 1;
	}

	sin(radians) {
		print "Calculating... sin " + str(radians);
		return nil;
	}
}

trait Cos {
	class cos30 {
		return 0.15;
	}

	class cos45 {
		return 0.52;
	}

	class cos90 {
		return 0;
	}

	cos(radians) {
		print "Calculating... cos " + str(radians);
		return nil;
	}
}

class Math {
	fib(n) {
		if (n <= 1) {
			return n;
		}
		return this.fib(n - 1) + this.fib(n - 2);
	}
}

class MyMath < Math <> Sin, Cos {
	factorial(n) {
		if (n == 0) {
			return 1;
		}
		return n * this.factorial(n - 1);
	}
}

print MyMath.sin30;
print MyMath.sin45;
print MyMath.sin90;
print MyMath.cos30;
print MyMath.cos45;
print MyMath.cos90;
MyMath().sin(270);
MyMath().cos(360);
print MyMath().fib(5);
print MyMath().factorial(5);
`
	assertPrograms(t, []testCase{
		{program1, "<trait Scanner>\n<trait Parser>\nParsing... [a, b, c] Done.\nInterpreting...\n"},
		{program2, "0.5\n0.07\n1\n0.15\n0.52\n0\nCalculating... sin 270\nCalculating... cos 360\n5\n120\n"},
	})
}
