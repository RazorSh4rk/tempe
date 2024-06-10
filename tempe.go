package tempe

import (
	"errors"
	"regexp"
	"strings"
)

type Sub struct {
	Key      string
	Value    string
	Function func(string, int) string
	Repeat   bool
	Regex    bool
}

type Subs struct {
	Subs        []Sub
	FailOnErr   bool
	ErrCallback func(error, int)
}

func (s *Sub) Apply(template *string) (err error) {
	if s.Key == "" {
		return errors.New("need a key to replace")
	}

	var processed string
	regex, err := regexp.Compile(s.Key)

	if err != nil || !s.Regex {
		// use simple string matching
		if s.Function == nil {
			// replace with value
			if s.Repeat {
				// replace multiple
				processed = rA(*template, s.Key, s.Value)
			} else {
				// replace one
				processed = r(*template, s.Key, s.Value)
			}
		} else {
			// replace with fn
			val := s.Function("", 0)
			if s.Repeat {
				processed = rA(*template, s.Key, val)
			} else {
				processed = r(*template, s.Key, val)
			}
		}
	} else {
		// use regex
		if s.Function == nil {
			// replace with value
			if s.Repeat {
				// replace multiple
				processed = rRA(*template, regex, func(found string, idx int) string {
					return s.Value
				})
			} else {
				// replace one
				processed = rR(*template, regex, func(found string, idx int) string {
					return s.Value
				})
			}
		} else {
			// replace with fn
			if s.Repeat {
				processed = rRA(*template, regex, s.Function)
			} else {
				processed = rR(*template, regex, s.Function)
			}
		}
	}

	*template = processed

	return
}

func (s Subs) ApplyAll(template *string) (err error) {
	for i, sub := range s.Subs {
		err = sub.Apply(template)
		if err != nil {
			if s.FailOnErr {
				return err
			} else {
				if s.ErrCallback != nil {
					s.ErrCallback(err, i)
				}
			}
		}
	}

	return nil
}

func r(s string, key string, value string) string {
	return strings.Replace(s, key, value, 1)
}

func rA(s string, key string, value string) string {
	return strings.ReplaceAll(s, key, value)
}

func rR(s string, regex *regexp.Regexp, callback func(string, int) string) string {
	first := regex.Find([]byte(s))
	replacement := callback(string(first), 0)

	return regex.ReplaceAllString(s, replacement)
}

func rRA(s string, regex *regexp.Regexp, callback func(string, int) string) string {
	all := regex.FindAllString(s, -1)
	for idx, found := range all {
		s = r(s, found, callback(found, idx))
	}
	return s
}
