package demo

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// KeyCaseInsensitiveMapType compares map keys case-insensitively for semantic equality.
type KeyCaseInsensitiveMapType struct {
	basetypes.MapType
}

var _ basetypes.MapTypable = KeyCaseInsensitiveMapType{}

func (t KeyCaseInsensitiveMapType) String() string {
	return "customtypes.KeyCaseInsensitiveMapType"
}

func (t KeyCaseInsensitiveMapType) ValueType(ctx context.Context) attr.Value {
	return KeyCaseInsensitiveMapValue{MapValue: basetypes.NewMapNull(t.ElementType())}
}

func (t KeyCaseInsensitiveMapType) Equal(other attr.Type) bool {
	o, ok := other.(KeyCaseInsensitiveMapType)
	return ok && t.MapType.Equal(o.MapType)
}

func (t KeyCaseInsensitiveMapType) WithElementType(elementType attr.Type) attr.TypeWithElementType {
	return KeyCaseInsensitiveMapType{MapType: basetypes.MapType{ElemType: elementType}}
}

func (t KeyCaseInsensitiveMapType) ValueFromMap(ctx context.Context, value basetypes.MapValue) (basetypes.MapValuable, diag.Diagnostics) {
	return KeyCaseInsensitiveMapValue{MapValue: value}, nil
}

func (t KeyCaseInsensitiveMapType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	value, err := t.MapType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	mapValue, ok := value.(basetypes.MapValue)
	if !ok {
		return nil, fmt.Errorf("unexpected value type %T", value)
	}
	return KeyCaseInsensitiveMapValue{MapValue: mapValue}, nil
}

type KeyCaseInsensitiveMapValue struct {
	basetypes.MapValue
}

var _ basetypes.MapValuableWithSemanticEquals = KeyCaseInsensitiveMapValue{}

func (v KeyCaseInsensitiveMapValue) Type(ctx context.Context) attr.Type {
	return KeyCaseInsensitiveMapType{MapType: basetypes.MapType{ElemType: v.ElementType(ctx)}}
}

func (v KeyCaseInsensitiveMapValue) Equal(other attr.Value) bool {
	o, ok := other.(KeyCaseInsensitiveMapValue)
	return ok && v.MapValue.Equal(o.MapValue)
}

func (v KeyCaseInsensitiveMapValue) KeysEqual(other KeyCaseInsensitiveMapValue) bool {
	return v.semanticEquals(other, true)
}

func (v KeyCaseInsensitiveMapValue) MapSemanticEquals(ctx context.Context, other basetypes.MapValuable) (bool, diag.Diagnostics) {
	o, ok := other.(KeyCaseInsensitiveMapValue)
	if !ok {
		var diags diag.Diagnostics
		diags.AddError("Semantic Equality Check Error",
			fmt.Sprintf("Expected KeyCaseInsensitiveMapValue, received %T. Please report this to the provider developers.", other))
		return false, diags
	}
	return v.semanticEquals(o, true), nil
}

func (v KeyCaseInsensitiveMapValue) semanticEquals(o KeyCaseInsensitiveMapValue, keysOnly bool) bool {
	if v.IsUnknown() || o.IsUnknown() {
		return false
	}
	if v.IsNull() || o.IsNull() {
		return v.MapValue.Equal(o.MapValue)
	}
	if v.MapValue.Equal(o.MapValue) {
		return true
	}
	elements, otherElements := v.Elements(), o.Elements()
	if len(elements) != len(otherElements) {
		return false
	}
	normalized, ok := normalizeMapKeys(elements)
	if !ok {
		return false
	}
	otherNormalized, ok := normalizeMapKeys(otherElements)
	if !ok {
		return false
	}
	for key, value := range normalized {
		otherValue, ok := otherNormalized[key]
		if !ok || (!keysOnly && !value.Equal(otherValue)) {
			return false
		}
	}
	return true
}

func normalizeMapKeys(elements map[string]attr.Value) (map[string]attr.Value, bool) {
	normalized := make(map[string]attr.Value, len(elements))
	for key, value := range elements {
		key = strings.ToLower(key)
		if _, exists := normalized[key]; exists {
			return nil, false
		}
		normalized[key] = value
	}
	return normalized, true
}
