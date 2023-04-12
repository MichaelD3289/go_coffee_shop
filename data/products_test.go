package data

import "testing"

func TestChecksValidation(t *testing.T) {
	p := &Product{
		Name: "Tea",
		Price: 1.00,
		SKU: "abc-abc-abcd",
	}

	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}