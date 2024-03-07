package test

import "testing"

func TestGlobalVariables(t *testing.T) {
	program1 := `
var v = "global";

print v;
{
	print v;
	{
		print v;
	}
}
print v;
`

	program2 := `
var v = "global";
print v;
{
	print v;
	v = "inner";
	{
		print v;
	}
}
print v;
`

	testPrograms(t, []testCase{
		{program1, "global\nglobal\nglobal\nglobal\n"},
		{program2, "global\nglobal\ninner\ninner\n"},
	})
}

func TestAssignment(t *testing.T) {
	program := `
var v = 2;
print v;
v = v * 2;
print v;
v = v * 2;
print v;
`

	testPrograms(t, []testCase{
		{program, "2\n4\n8\n"},
	})
}

func TestScope(t *testing.T) {
	program1 := `
var volume = 11;

volume = 0;

{
  var volume = 3 * 4 * 5;
  print volume;
}
print volume;
`

	program2 := `
var a = "global a";
var b = "global b";
var c = "global c";
{
  var a = "outer a";
  var b = "outer b";
  {
    var a = "inner a";
    print a;
    print b;
    print c;
  }
  print a;
  print b;
  print c;
}
print a;
print b;
print c;
`

	testPrograms(t, []testCase{
		{program1, "60\n0\n"},
		{program2, "inner a\nouter b\nglobal c\nouter a\nouter b\nglobal c\nglobal a\nglobal b\nglobal c\n"},
	})
}
