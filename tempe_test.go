package tempe_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/razorsh4rk/tempe"
)

// replace a string with an other string
func TestOne(t *testing.T) {
	s := "hello world"

	template := tempe.Sub{
		Key:   "hello",
		Value: "bye",
	}

	template.Apply(&s)

	if s != "bye world" {
		t.Fail()
	}
}

// replace a string with a static function
func TestFunc(t *testing.T) {
	s := "host is: /hname/"
	name, _ := os.Hostname()

	(&tempe.Sub{
		Key: "/hname/",
		Function: func(s string, i int) string {
			return name
		},
	}).Apply(&s)

	if s != "host is: "+name {
		t.Fail()
	}
}

// replace multiple strings with a string
func TestRepeating(t *testing.T) {
	s := "I love rust because rust is awesome"

	(&tempe.Sub{
		Key:    "rust",
		Value:  "go",
		Repeat: true,
	}).Apply(&s)

	if s != "I love go because go is awesome" {
		t.Fail()
	}
}

// replace a regex string with a dynamic function
func TestRegex(t *testing.T) {
	s := "let's count num num num num num"

	template := tempe.Sub{
		Key:    "num",
		Regex:  true,
		Repeat: true,
		Function: func(s string, i int) string {
			return fmt.Sprint(i)
		},
	}

	template.Apply(&s)

	if s != "let's count 0 1 2 3 4" {
		fmt.Println(s)
		t.Fail()
	}
}

// replace the first match of a regex with a string
func TestRegexSingle(t *testing.T) {
	s := "cats cats dogs cats"

	(&tempe.Sub{
		Key:    "dog(s?)",
		Value:  "cats",
		Regex:  true,
		Repeat: false,
	}).Apply(&s)

	fmt.Println(s)
	if s != "cats cats cats cats" {
		t.Fail()
	}
}

// replace the first match of a regex with a dynamic function
func TestRegexSingleFn(t *testing.T) {
	s := "cats cats dog cats"

	(&tempe.Sub{
		Key:    "dog(s?)",
		Regex:  true,
		Repeat: false,
		Function: func(s string, i int) string {
			switch s {
			case "dog":
				return "cat"
			case "dogs":
				return "cats"
			}

			return ""
		},
	}).Apply(&s)

	fmt.Println(s)
	if s != "cats cats cat cats" {
		t.Fail()
	}
}

// replace a string with a dynamic function
func TestNoRegex(t *testing.T) {
	s := "let's count num num num num num"

	template := tempe.Sub{
		Key:    "num",
		Regex:  false,
		Repeat: true,
		Function: func(s string, i int) string {
			return fmt.Sprint(i)
		},
	}

	template.Apply(&s)

	// if regex is false, the function gets called with ("", 0)
	if s != "let's count 0 0 0 0 0" {
		t.Fail()
	}
}

// replace windows slashes with unix slashes
func TestPath(t *testing.T) {
	s := "\\home\\user\\app"

	template := tempe.Sub{
		Key:    "\\",
		Regex:  true,
		Repeat: true,
		Value:  "/",
	}

	template.Apply(&s)

	if s != "/home/user/app" {
		t.Fail()
	}
}

// replace a regex with dynamic values based on the matches
func TestRegexDiffering(t *testing.T) {
	s := "{name1} and {name2} and {name3}"

	template := tempe.Sub{
		Key:    "{name[1-3]}",
		Regex:  true,
		Repeat: true,
		Function: func(s string, i int) string {
			var name string
			switch s {
			case "{name1}":
				name = "Joe Swanson"
			case "{name2}":
				name = "Peter Griffin"
			case "{name3}":
				name = "Glenn Quagmire"
			}

			return name
		},
	}

	template.Apply(&s)

	if s != "Joe Swanson and Peter Griffin and Glenn Quagmire" {
		t.Fail()
	}
}

// chain multiple replacers
func TestMultiple(t *testing.T) {
	s := "the time is {{hour}}h {{minute}}m"

	templates := tempe.Subs{
		Subs: []tempe.Sub{
			{
				Key:   "{{hour}}",
				Value: "12",
			},
			{
				Key:   "{{minute}}",
				Value: "00",
			},
		},
	}

	templates.ApplyAll(&s)

	if s != "the time is 12h 00m" {
		t.Fail()
	}
}

// chain multiple replacers and handle possible errors
// useful for building apps where the user supplies
// the replacer keys
func TestMultipleErr(t *testing.T) {
	s := "{{value}} {{error}}"

	templates := tempe.Subs{
		Subs: []tempe.Sub{
			{
				Key:   "{{value}}",
				Value: "57",
			},
			{},
		},
		FailOnErr: false,
		ErrCallback: func(err error, i int) {
			//log.Printf("Failure at index %d: %s", i, err.Error())
		},
	}

	templates.ApplyAll(&s)

	if s != "57 {{error}}" {
		t.Fail()
	}
}
