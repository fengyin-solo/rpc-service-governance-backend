package admission

import "rpcgate/internal/policy"

type Gate struct{ rules *policy.Rules }

func NewGate(rules *policy.Rules) *Gate { return &Gate{rules: rules} }

func (g *Gate) Allow(method string) error {
	if g.rules.Validator != nil {
		if err := g.rules.Validator.Validate(method); err != nil {
			return err
		}
	}
	g.rules.Methods[method] = true
	return nil
}
