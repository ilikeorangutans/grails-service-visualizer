package main

import (
	"bufio"
	"bytes"
	"testing"
)

func TestDependencyScanner(t *testing.T) {

	reader := bufio.NewReader(bytes.NewBufferString("package foo.bar\nclass FooBar {\n DummyService dummyService\n \tdef typelessService\n def semicolonService; \nString somethingElse}"))
	deps := ScanForDependencies("foobar.groovy", reader)

	if len(deps) != 3 {
		t.Error("Expected 3 dependencies")
	}

	// TODO: assert deps details
}

func TestRegexp(t *testing.T) {

	inputs := map[string]string{
		"def fooService":        "fooService",
		"String name":           "",
		"FooService fooService": "fooService",
		"  	SpaceService  	  spaceService      ": "spaceService",
		"// Commented out service": "",
		"def someService  ;    ":   "someService",
		"def someService;":         "someService",
		"    def userService\r\n":  "userService",
	}

	for in, out := range inputs {
		// If the expected result is empty yet we do match the input string then something's wrong with the regex.
		if len(out) == 0 && serviceRegex.MatchString(in) {
			t.Errorf("'%s' matched but was not supposed to", in)
			continue
		} else if len(out) == 0 {
			continue
		}

		submatch := serviceRegex.FindStringSubmatch(in)
		if submatch[1] != out {
			t.Errorf("Got '%s' but expected '%s'", out, submatch[1])
		}

	}

}
