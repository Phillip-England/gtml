package gtml

import (
	"fmt"
	"strings"

	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyPropForType     = "PROPFORTYPE"
	KeyPropForStr      = "PROPFORSTR"
	KeyPropStr         = "PROPSTR"
	KeyPropPlaceholder = "PROPPLACEHOLDER"
)

// ##==================================================================
type Prop interface {
	GetRaw() string
	GetValue() string
	GetType() string
	Print()
}

func NewProp(str string, compNames []string) (Prop, error) {
	val := purse.Squeeze(str)
	val = purse.RemoveAllSubStr(val, "{{", "}}")
	if val == "" {
		return nil, fmt.Errorf("empty prop tag provided: %s", str)
	}
	if strings.Contains(val, "(") && strings.Contains(val, ")") {
		openIndex := strings.Index(val, "(")
		potentialName := val[:openIndex]
		if purse.MustEqualOneOf(potentialName, compNames...) {
			prop, err := NewPropPaceholder(str, val)
			if err != nil {
				return nil, err
			}
			return prop, nil
		}
	}
	if len(val) > 1 && val[0] == '.' {
		val = strings.Replace(val, ".", "", 1)
		prop, err := NewPropForStr(str, val)
		if err != nil {
			return nil, err
		}
		return prop, nil
	}
	if strings.Count(val, ".") == 1 && len(strings.Split(val, ".")) == 2 {
		prop, err := NewPropForType(str, val)
		if err != nil {
			return nil, err
		}
		return prop, nil
	}
	prop, err := NewPropStr(str, val)
	if err != nil {
		return nil, err
	}
	return prop, nil
}

func PropAsWriteStringCall(prop Prop, builderName string) string {
	return fmt.Sprintf("%s.WriteString(%s)", builderName, prop.GetValue())
}

// ##==================================================================
type PropForType struct {
	Raw   string
	Value string
	Type  string
}

func NewPropForType(raw string, val string) (*PropForType, error) {
	prop := &PropForType{
		Raw:   raw,
		Value: val,
		Type:  KeyPropForType,
	}
	return prop, nil
}

func (prop *PropForType) GetRaw() string   { return prop.Raw }
func (prop *PropForType) GetValue() string { return prop.Value }
func (prop *PropForType) GetType() string  { return prop.Type }
func (prop *PropForType) Print() {
	fmt.Println(fmt.Sprintf("raw: %s\nvalue: %s\ntype: %s", prop.Raw, prop.Value, prop.Type))
}

// ##==================================================================
type PropForStr struct {
	Raw   string
	Value string
	Type  string
}

func NewPropForStr(raw string, val string) (*PropForStr, error) {
	prop := &PropForStr{
		Raw:   raw,
		Value: val,
		Type:  KeyPropForStr,
	}
	return prop, nil
}

func (prop *PropForStr) GetRaw() string   { return prop.Raw }
func (prop *PropForStr) GetValue() string { return prop.Value }
func (prop *PropForStr) GetType() string  { return prop.Type }
func (prop *PropForStr) Print() {
	fmt.Println(fmt.Sprintf("raw: %s\nvalue: %s\ntype: %s", prop.Raw, prop.Value, prop.Type))
}

// ##==================================================================
type PropStr struct {
	Raw   string
	Value string
	Type  string
}

func NewPropStr(raw string, val string) (*PropStr, error) {
	prop := &PropStr{
		Raw:   raw,
		Value: val,
		Type:  KeyPropStr,
	}
	return prop, nil
}

func (prop *PropStr) GetRaw() string   { return prop.Raw }
func (prop *PropStr) GetValue() string { return prop.Value }
func (prop *PropStr) GetType() string  { return prop.Type }
func (prop *PropStr) Print() {
	fmt.Println(fmt.Sprintf("raw: %s\nvalue: %s\ntype: %s", prop.Raw, prop.Value, prop.Type))
}

// ##==================================================================
type PropPlaceholder struct {
	Raw   string
	Value string
	Type  string
}

func NewPropPaceholder(raw string, val string) (*PropPlaceholder, error) {
	prop := &PropPlaceholder{
		Raw:   raw,
		Value: val,
		Type:  KeyPropPlaceholder,
	}
	return prop, nil
}

func (prop *PropPlaceholder) GetRaw() string   { return prop.Raw }
func (prop *PropPlaceholder) GetValue() string { return prop.Value }
func (prop *PropPlaceholder) GetType() string  { return prop.Type }
func (prop *PropPlaceholder) Print() {
	fmt.Println(fmt.Sprintf("raw: %s\nvalue: %s\ntype: %s", prop.Raw, prop.Value, prop.Type))
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
