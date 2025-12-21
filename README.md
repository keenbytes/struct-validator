# struct-validator

Verify the values of struct fields using tags

### Example code

Use the package with the following URL:
```
import "github.com/keenbytes/struct-validator"
```

And see the below code snippet:
```
type Test1 struct {
	FirstName     string `validation:"req lenmin:5 lenmax:25"`
	LastName      string `validation:"req lenmin:2 lenmax:50"`
	Age           int    `validation:"req valmin:18 valmax:150"`
	Price         int    `validation:"req valmin:0 valmax:9999"`
	PostCode      string `validation:"req" validation_regexp:"^[0-9][0-9]-[0-9][0-9][0-9]$"`
	Email         string `validation:"req email"`
	BelowZero     int    `validation:"valmin:-6 valmax:-2"`
	DiscountPrice int    `validation:"valmin:0 valmax:8000"`
	Country       string `validation_regexp:"^[A-Z][A-Z]$"`
	County        string `validation:"lenmax:40"`
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
