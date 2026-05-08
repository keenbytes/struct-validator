# struct-validator

Verify the values of struct fields using tags

### Example code

See the below code snippets:
```
type Test1 struct {
	FirstName     string `validation:"req len:5,25"`
	LastName      string `validation:"req len:2,50"`
	Age           int    `validation:"req val:18,150"`
	Price         int    `validation:"req val:0,9999"`
	PostCode      string `validation:"req" validation_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         string `validation:"req email"`
	BelowZero     int    `validation:"val:-6,-2"`
	DiscountPrice int    `validation:"val:0,8000"`
	Country       string `validation_regexp:"^[A-Z][A-Z]$"`
	County        string `validation:"len:,40"`
}

s := &Test1{
	FirstName: "Name that is way too long and certainly not valid",
	...
}

o := validator.&ValidationOptions{
	RestrictFields: map[string]bool{
		"FirstName": true,
		"LastName":  true,
		...
	},
	...
}

isValid, fieldViolations, err := validator.Validate(s, &o)
```
