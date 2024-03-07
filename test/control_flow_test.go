package test

import "testing"

func TestIfStatements(t *testing.T) {
	program := `
var num = 32;
if (num % 2 == 0) {
	num = num / 16;
}
print num;

var prime = 15;
if (prime % 2 == 0) {
	prime = prime / 2;
} else {
	prime = prime / 3;
}

print prime;
`

	testPrograms(t, []testCase{
		{program, "2\n5\n"},
	})
}

func TestWhileStatements(t *testing.T) {
	program1 := `
var i = 0;
while (i < 10) {
	i = i + 1;
}

print i;
`

	program2 := `
var a = 0;
var b = 1;
while (true) {
	if (a > 200) {
		break;
	}

	var c = a + b;
	a = b;
	b = c;
}

print a;
`
	program3 := `
var i = 0;
while (true) {
    i = i + 1;

    if (i >= 5 and i < 10) {
        print "Finished for real";
        continue;
    }

    if (i == 10) {
        print "Finished for real now";
        break;
    }

    while (true) {
        if (i > 5) {
            print "Finished";
            break;
        }
        i = i + 1;
    }
}
`

	testPrograms(t, []testCase{
		{program1, "10\n"},
		{program2, "233\n"},
		{program3, "Finished\nFinished for real\nFinished for real\nFinished for real\nFinished for real now\n"},
	})
}

func TestForStatements(t *testing.T) {
	program1 := `
var i = 1;
for (;;) {
	if (i == 32) {
		break;
	}

	i = i * 2;
}
print i;

for (var i = 1; ;) {
	if (i == 32) {
		break;
	}

	i = i * 2;
}
print i;

for (var i = 1; i < 32;) {
	i = i * 2;
}
print i;

for (var i = 1; i < 32; i = i * 2) {}
print i;
`
	program2 := `
for (var i = 0; i < 2; i = i + 1) {
	for (var j = 0; j < 2; j = j + 1) {
		print "[" + str(i) + "]" + "[" + str(j) + "]";
	}
}
`

	testPrograms(t, []testCase{
		{program1, "32\n32\n32\n32\n"},
		{program2, "[0][0]\n[0][1]\n[1][0]\n[1][1]\n"},
	})
}
