package test

import "testing"

func TestArrays(t *testing.T) {
	program := `
var empty = [];
print empty;

var primes = [2, 3, 5, 7, 11];
print primes[2];
print primes;

var byte = [0, 0, 0, 0, 0, 0, 0, 0];
byte[7] = 1;
byte[4] = 1;
byte[0] = 1;


var result = 0;
for (var i = 0; i < 8; i = i + 1) {
	if (byte[(7 - i)] == 0) {
		continue;
	}

	var power = 1;
	for (var j = 0; j < i; j = j + 1) {
		power = power * 2;
	}
	result = result + power;
}

print "Binary" + str(byte) + " == " + str(result);
`
	assertPrograms(t, []testCase{
		{program, "[]\n5\n[2, 3, 5, 7, 11]\nBinary[1, 0, 0, 0, 1, 0, 0, 1] == 137\n"},
	})
}
