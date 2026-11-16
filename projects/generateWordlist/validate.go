package main

import "regexp"

type expression struct {
	re        *regexp.Regexp
	sample    []byte
	minLength int
}

type vStatus struct {
	matched       bool
	result        bool
	toBeContinued bool
}

// create a new []byte for testing
func prepare(bs []byte, e expression) []byte {
	testing := make([]byte, len(e.sample))
	var i int
	for i < len(bs) {
		testing[i] = bs[i]
		i++
	}
	for i < len(e.sample) {
		testing[i] = e.sample[i]
		i++
	}
	return testing

}
func validate(bs []byte, ex expression) vStatus {

	testing := prepare(bs, ex)
	minLength := len(ex.sample)
	maxLength := len(ex.sample)
	if ex.minLength > 0 {
		minLength = ex.minLength
	}

	if ex.re.Match(testing) {
		matched := true
		toBeContinued := len(bs) < maxLength
		result := len(bs) >= minLength
		return vStatus{
			matched:       matched,
			toBeContinued: toBeContinued,
			result:        result,
		}
	} else {
		return vStatus{}
	}
}
