package policy

type Validator interface {
	Validate(string) error
}

type MethodValidator struct{}

func (*MethodValidator) Validate(string) error { return nil }

type Rules struct {
	Methods   map[string]bool
	Validator Validator
}

func LoadDefault() *Rules {
	var validator *MethodValidator
	return &Rules{Validator: validator}
}
