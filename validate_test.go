package validator

import (
	"log"
	"testing"
)

type Test1 struct {
	FirstName     string `validation:"len:5,25"`
	LastName      string `validation:"len:2,50"`
	Age           int    `validation:"val:18,150"`
	Price         int    `validation:"val:0,9999"`
	PostCode      string `validation:"" validation_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         string `validation:"email"`
	BelowZero     int    `validation:"val:-6,-2"`
	DiscountPrice int    `validation:"val:0,8000"`
	Country       string `validation_regexp:"^[A-Z][A-Z]$"`
	County        string `validation:"len:,40"`
}

type Test2 struct {
	FirstName     string `mytag:"req len:5,25"`
	LastName      string `mytag:"req len:2,50"`
	Age           int    `mytag:"req val:18,150"`
	Price         int    `mytag:"req val:0,9999"`
	PostCode      string `mytag:"req" mytag_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         string `mytag:"req email"`
	BelowZero     int    `mytag:"val:-6,-2"`
	DiscountPrice int    `mytag:"val:0,8000"`
	Country       string `mytag_regexp:"^[A-Z][A-Z]$"`
	County        string `mytag:"len:,40"`
}

type Test3 struct {
	FirstName     *string `mytag:"req len:5,25"`
	LastName      *string `mytag:"len:2,50"`
	Age           *int    `mytag:"req val:18,150"`
	Price         *int    `mytag:"val:0,9999"`
	PostCode      *string `mytag:"req" mytag_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         *string `mytag:"req email"`
	BelowZero     *int    `mytag:"val:-6,-2"`
	DiscountPrice *int    `mytag:"val:0,8000"`
	Country       *string `mytag:"^[A-Z][A-Z]$"`
	County        *string `mytag:"len:40"`
}

func TestWithDefaultValues(t *testing.T) {
	s := Test1{}

	expectedViolations := map[string]uint64{
		"FirstName": FailLenMin,
		"LastName":  FailLenMin,
		"Age":       FailValMin,
		"PostCode":  FailRegExp,
		"Email":     FailEmail,
		"Country":   FailRegExp,
		"BelowZero": FailValMax,
	}

	ok, violations, _ := Validate(s, &ValidateOptions{})
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithInvalidValues(t *testing.T) {
	s := Test1{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
	}

	expectedViolations := map[string]uint64{
		"FirstName":     FailLenMax,
		"LastName":      FailLenMin,
		"Age":           FailValMin,
		"PostCode":      FailRegExp,
		"Email":         FailEmail,
		"BelowZero":     FailValMax,
		"DiscountPrice": FailValMax,
		"Country":       FailRegExp,
	}

	ok, violations, _ := Validate(s, &ValidateOptions{})
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithValidValues(t *testing.T) {
	s := Test1{
		FirstName:     "Johnny",
		LastName:      "Smith",
		Age:           35,
		Price:         0,
		PostCode:      "43-155",
		Email:         "john@example.com",
		BelowZero:     -4,
		DiscountPrice: 8000,
		Country:       "GB",
		County:        "Enfield",
	}

	expectedViolations := map[string]uint64{}

	ok, violations, _ := Validate(s, &ValidateOptions{})
	if !ok {
		t.Fatalf("validation should have succeeded")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithInvalidValuesAndFieldRestriction(t *testing.T) {
	s := Test1{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
	}

	expectedViolations := map[string]uint64{
		"FirstName": FailLenMax,
		"LastName":  FailLenMin,
	}
	opts := &ValidateOptions{
		RestrictFields: map[string]bool{
			"FirstName": true,
			"LastName":  true,
		},
	}

	ok, violations, _ := Validate(s, opts)
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithInvalidValuesAndOverwrittenTagName(t *testing.T) {
	s := Test2{
		FirstName:     "123456789012345678901234567890",
		LastName:      "b",
		Age:           15,
		Price:         0,
		PostCode:      "AA123",
		Email:         "invalidEmail",
		BelowZero:     8,
		DiscountPrice: 9999,
		Country:       "Tokelau",
		County:        "",
	}
	expectedViolations := map[string]uint64{
		"FirstName":     FailLenMax,
		"LastName":      FailLenMin,
		"Age":           FailValMin,
		"PostCode":      FailRegExp,
		"Email":         FailEmail,
		"BelowZero":     FailValMax,
		"DiscountPrice": FailValMax,
		"Country":       FailRegExp,
	}
	opts := &ValidateOptions{
		TagName: "mytag",
	}
	ok, violations, _ := Validate(s, opts)
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithAllInvalidValuesAndPointerFields(t *testing.T) {
	s := Test3{}

	expectedViolations := map[string]uint64{
		"FirstName": FailReq,
		"Age":       FailReq,
		"PostCode":  FailReq,
		"Email":     FailReq,
	}
	opts := &ValidateOptions{
		TagName: "mytag",
	}
	ok, violations, _ := Validate(s, opts)
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithInvalidValuesAndPointerFields(t *testing.T) {
	firstName := "a"
	age := 3
	postCode := "a123"

	s := Test3{
		FirstName: &firstName,
		Age:       &age,
		PostCode:  &postCode,
	}

	expectedViolations := map[string]uint64{
		"FirstName": FailLenMin,
		"Age":       FailValMin,
		"PostCode":  FailRegExp,
		"Email":     FailReq,
	}
	opts := &ValidateOptions{
		TagName: "mytag",
	}
	ok, violations, _ := Validate(s, opts)
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithOverwriteValues(t *testing.T) {
	s := Test2{
		FirstName: "First",
		Age:       40,
		PostCode:  "11-111",
		Email:     "email",
	}

	expectedViolations := map[string]uint64{
		"Age":      FailValMax,
		"PostCode": FailRegExp,
		"Email":    FailType,
	}
	opts := &ValidateOptions{
		TagName: "mytag",
		RestrictFields: map[string]bool{
			"Age":      true,
			"PostCode": true,
			"Email":    true,
		},
		OverwriteValues: map[string]interface{}{
			"Age":      400,
			"PostCode": "123-456",
			"Email":    44,
		},
	}
	ok, violations, _ := Validate(s, opts)
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func TestWithOverwriteValuesAndPointerFields(t *testing.T) {
	age := 300
	postCode := "a123"
	email := 44

	s := Test3{}

	expectedViolations := map[string]uint64{
		"Age":      FailValMax,
		"PostCode": FailRegExp,
		"Email":    FailType,
	}
	opts := &ValidateOptions{
		TagName: "mytag",
		RestrictFields: map[string]bool{
			"Age":      true,
			"PostCode": true,
			"Email":    true,
		},
		OverwriteValues: map[string]interface{}{
			"Age":      &age,
			"PostCode": &postCode,
			"Email":    &email,
		},
	}
	ok, violations, _ := Validate(s, opts)
	if ok {
		t.Fatalf("validation should have failed")
	}

	compareViolations(violations, expectedViolations, t)
}

func compareViolations(violations map[string]uint64, expectedViolations map[string]uint64, t *testing.T) {
	if len(violations) != len(expectedViolations) {
		log.Printf("Violations:\n")
		for k, v := range violations {
			log.Printf("%s %d", k, v)
		}
		t.Fatalf("Validate returned invalid number of failed fields %d where it should be %d", len(violations), len(expectedViolations))
	}
	for k, v := range expectedViolations {
		if violations[k] != v {
			t.Fatalf("Validate returned invalid failure flag of %d where it should be %d for %s", violations[k], v, k)
		}
	}
}
