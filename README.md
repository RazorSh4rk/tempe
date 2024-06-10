# tempe

`go get github.com/RazorSh4rk/tempe`

Replace string templates with functions

Tempe (tem-pee) is a small library that lets you create 
chainable replacers with dynamic functions embedded in
them to transform text. Originally made as part of a 
[terminal chatGPT client](https://github.com/RazorSh4rk/chatty), but I extracted the code because it could be useful for other apps too.

Note: like 3 minutes after I made this I found [fasttemplate](https://github.com/valyala/fasttemplate) which does kinda the same thing, so if you prefer that syntax or just hate me personally, that's a good alternative.

## Docs

Tempe comes with 2 structs:

```golang
type Sub struct {
	Key      string
	Value    string
	Function func(string, int) string
	Repeat   bool
	Regex    bool
}

Sub{...}.Apply(&yourString)
```
This is the base of the library, it represents one substitution.

- _Key_: string or regex string to replace

- _Value_: string to replace with, or empty

- _Function_: a function that gets called for a replacement. The inputs are the found match, and the index of the found match. Return any string to act as a replacement.

If you are using regex, the entire match will be returned, for example `{{myKey}}`

If you are *NOT* using regex, the function will be called with ("", 0) every time. If you want to have access to these variables, just set `Regex` to `true`, it will still match standard strings.

- _Repeat_: replace the first hit (false) or all of them (true)

- _Regex_: use regex matching (true) or simple string matching (false)

```golang
type Subs struct {
	Subs        []Sub
	FailOnErr   bool
	ErrCallback func(error, int)
}

Subs{...}.ApplyAll(&yourString)
```

This is a chain of replacers that can all be called in order and can be supplied with a callback in case an error happens. Useful if the replacers are user supplied or if you are building a template engine.

- _ErrCallback_: will receive the error and the index of the replacer which received the error.

## Examples

There are a lot of examples in the test file, but here are the basics:

_There is a real world usage example in the /example folder_

#### _Replace a thing with an other thing_
```golang
s := "hello world"
template := tempe.Sub{
	Key:   "hello",
	Value: "bye",
}
template.Apply(&s)
// bye world
```

or, without storing the sub:

```golang
s := "hello world"
(&tempe.Sub{
	Key:   "hello",
	Value: "bye",
}).Apply(&s)
```

#### _replace a string with a function_

```golang
s := "host is: /hname/"
(&tempe.Sub{
	Key: "/hname/",
	Function: func(s string, i int) string {
		name, _ := os.Hostname()
        return name
	},
}).Apply(&s)
// host is: yourhost.name
```

#### _replace the first match of a regex with a function based on the matched value_

```golang
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
// cats cats cat cats
```

#### _replace a regex with the index of where it was found_

```golang
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
// let's count 0 1 2 3 4

// setting the Regex to false here would result in
// let's count 0 0 0 0 0
```

#### _replace multiple regex matches based on their values_

```golang
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
// Joe Swanson and Peter Griffin and Glenn Quagmire
```

#### _chain multiple replaces_

This works with any of the previous replacers too

```golang
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
// the time is 12h 00m
```