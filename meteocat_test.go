package meteocat

import (
	"testing"
)

// TestValidDataUnit tests whether or not ValidDataUnit provides
// the correct assertion on provided data unit.
func TestValidCodisEstacions(t *testing.T) {
	for u := range CodisEstacions {
		if !ValidCodiEstacio(u) {
			t.Error("False positive on data unit")
		}
	}

	if ValidCodiEstacio("anything") {
		t.Error("Invalid data unit")
	}
}

// TestValidCodiVariable tests whether or not ValidCodiVariable provides
// the correct assertion on provided data unit.
func TestValidCodiVariable(t *testing.T) {
	for _, s := range CodisVariables {
		if !ValidCodiVariable(s) {
			t.Error("False positive on data unit symbol")
		}
	}

	if ValidCodiVariable("X") {
		t.Error("Invalid data unit symbol")
	}
}

// TestCheckAPIKeyExists tests whether or not CheckAPIKeyExists provides
// the correct assertion on provided data unit.
func TestCheckAPIKeyExists(t *testing.T) {
		apiKey := "asdf1234"

	if !CheckAPIKeyExists(apiKey) {
		t.Error("Key not set")
	}
}

// TestSetOptionsWithEmpty tests setOptions function will do nothing
// when options are empty.
func TestSetOptionsWithEmpty(t *testing.T) {
	s := NewSettings()
	err := setOptions(s, nil)
	if err != nil {
		t.Error(err)
	}
}
