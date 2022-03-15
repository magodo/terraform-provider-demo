package types

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type JsonType struct{}

var _ attr.Type = JsonType{}

// TerraformType returns the tftypes.Type that should be used to
// represent this type. This constrains what user input will be
// accepted and what kind of data can be set in state. The framework
// will use this to translate the Type to something Terraform can
// understand.
func (j JsonType) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// ValueFromTerraform returns a Value given a tftypes.Value. This is
// meant to convert the tftypes.Value into a more convenient Go type
// for the provider to consume the data with.
func (j JsonType) ValueFromTerraform(_ context.Context, in tftypes.Value) (attr.Value, error) {
	// Following is copied from the impl of types.primitive
	if !in.IsKnown() {
		return types.String{Unknown: true}, nil
	}
	if in.IsNull() {
		return types.String{Null: true}, nil
	}
	var s string
	err := in.As(&s)
	if err != nil {
		return nil, err
	}
	return types.String{Value: s}, nil
}

// Equal must return true if the Type is considered semantically equal
// to the Type passed as an argument.
func (j JsonType) Equal(in attr.Type) bool {
	_, ok := in.(JsonType)
	return ok
}

// String should return a human-friendly version of the Type.
func (j JsonType) String() string {
	return "types.JsonType"
}

// Return the attribute or element the AttributePathStep is referring
// to, or an error if the AttributePathStep is referring to an
// attribute or element that doesn't exist.
func (j JsonType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, j.String())
}

func (t JsonType) Validate(ctx context.Context, tfValue tftypes.Value, path *tftypes.AttributePath) diag.Diagnostics {
	var diags diag.Diagnostics

	if !tfValue.Type().Equal(tftypes.String) {
		diags.AddAttributeError(
			path,
			"JSON Type Validation Error",
			fmt.Sprintf("Expected String value, received %T with value: %s", tfValue.Type(), tfValue),
		)
		return diags
	}

	if !tfValue.IsKnown() || tfValue.IsNull() {
		return diags
	}

	var value string
	err := tfValue.As(&value)

	if err != nil {
		diags.AddAttributeError(
			path,
			"JSON Type Validation Error",
			fmt.Sprintf("Cannot convert value to string: %s", err),
		)
		return diags
	}

	var j interface{}
	if err := json.Unmarshal([]byte(value), &j); err != nil {
		diags.AddAttributeError(
			path,
			"JSON Type Validation Error",
			fmt.Sprintf("Invalid JSON: %s", err),
		)
		return diags
	}

	return diags
}
