package test

import "testing"

func TestResolver(t *testing.T) {
	program := `
var a = "global";
{
  fun showA() {
    print a;
  }

  showA();
  var a = "block";
  showA();
  print a;
}
`
	testPrograms(t, []testCase{
		{program, "global\nglobal\nblock\n"},
	})
}
