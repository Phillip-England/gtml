package gtml

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

// ##==================================================================
const (
	KeyElementComponent = "_component"
	KeyElementFor       = "_for"
	KeyElementIf        = "_if"
	KeyElementElse      = "_else"
)

// ##==================================================================
type Element interface {
	GetSelection() *goquery.Selection
	GetParam() (Param, error)
	SetHtml(htmlStr string)
	GetHtml() string
	Print()
	GetType() string
	GetAttr() string
	GetAttrParts() []string
	GetProps() []Prop
}

func GetFullElementList() []string {
	childElements := GetChildElementList()
	full := append(childElements, KeyElementComponent)
	return full
}

func GetChildElementList() []string {
	return []string{KeyElementFor, KeyElementIf, KeyElementElse}
}

func NewElement(sel *goquery.Selection) (Element, error) {
	match := gqpp.GetFirstMatchingAttr(sel, GetFullElementList()...)
	switch match {
	case KeyElementComponent:
		elm, err := NewElementComponent(sel)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementFor:
		elm, err := NewElementFor(sel)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementIf:
		elm, err := NewElementIf(sel)
		if err != nil {
			return nil, err
		}
		return elm, nil
	case KeyElementElse:
		elm, err := NewElementElse(sel)
		if err != nil {
			return nil, err
		}
		return elm, nil
	}
	htmlStr, err := gqpp.NewHtmlFromSelection(sel)
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("provided selection is not a valid element: %s", htmlStr)
}

func WalkElementChildren(elm Element, fn func(child Element) error) error {
	var potErr error
	elm.GetSelection().Find("*").Each(func(i int, inner *goquery.Selection) {
		child, err := NewElement(inner)
		if err != nil {
			// skip elements which are not a valid Element
		} else {
			err = fn(child)
			if err != nil {
				potErr = err
				return
			}
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func GetElementParams(elm Element) ([]Param, error) {
	params := make([]Param, 0)
	elementSpecificParams := make([]Param, 0)
	err := WalkElementChildren(elm, func(child Element) error {
		param, err := child.GetParam()
		if err != nil {
			return err
		}
		if !slices.Contains(elementSpecificParams, param) && param != nil {
			elementSpecificParams = append(elementSpecificParams, param)
		}
		return nil
	})
	if err != nil {
		return params, err
	}
	strParams := make([]Param, 0)
	for _, prop := range elm.GetProps() {
		param, err := NewParam(prop.GetValue(), "string")
		if err != nil {
			return params, err
		}
		if !slices.Contains(strParams, param) && prop.GetType() == KeyPropStr {
			strParams = append(strParams, param)
		}
	}
	params = append(strParams, elementSpecificParams...)
	return params, nil
}

func WalkElementDirectChildren(elm Element, fn func(child Element) error) error {
	err := WalkElementChildren(elm, func(child Element) error {
		if !gqpp.HasParentWithAttrs(child.GetSelection(), elm.GetSelection(), GetChildElementList()...) {
			err := fn(child)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func WalkElementProps(elm Element, fn func(prop Prop) error) error {
	allProps := make([]Prop, 0)
	for _, prop := range elm.GetProps() {
		allProps = append(allProps, prop)
	}
	err := WalkElementChildren(elm, func(child Element) error {
		for _, prop := range child.GetProps() {
			allProps = append(allProps, prop)
		}
		return nil
	})
	if err != nil {
		return err
	}
	for _, prop := range allProps {
		err := fn(prop)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetElementProps(elm Element) ([]Prop, error) {
	props := make([]Prop, 0)
	elmHtml := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		elmHtml = strings.Replace(elmHtml, childHtml, "", 1)
		return nil
	})
	if err != nil {
		return props, err
	}
	strProps := purse.ScanBetweenSubStrs(elmHtml, "{{", "}}")
	for _, strProp := range strProps {
		prop, err := NewProp(strProp)
		if err != nil {
			return props, err
		}
		props = append(props, prop)
	}
	return props, nil
}

func GetElementVars(elm Element) ([]Var, error) {
	vars := make([]Var, 0)
	err := WalkElementDirectChildren(elm, func(child Element) error {
		innerVar, err := NewVar(child)
		if err != nil {
			return nil
		}
		vars = append(vars, innerVar)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return vars, nil
}

func GetElementAsBuilderSeries(elm Element, builderName string) (string, error) {
	clay := elm.GetHtml()
	err := WalkElementDirectChildren(elm, func(child Element) error {
		childHtml := child.GetHtml()
		newVar, err := NewVar(child)
		if err != nil {
			return err
		}
		varType := newVar.GetType()
		if purse.MustEqualOneOf(varType, KeyVarGoElse, KeyVarGoFor, KeyVarGoIf) {
			call := fmt.Sprintf("%s.WriteString(%s)", builderName, newVar.GetVarName())
			clay = strings.Replace(clay, childHtml, call, 1)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	err = WalkElementProps(elm, func(prop Prop) error {
		call := fmt.Sprintf("%s.WriteString(%s)", builderName, prop.GetValue())
		clay = strings.Replace(clay, prop.GetRaw(), call, 1)
		return nil
	})
	if err != nil {
		return "", err
	}
	if strings.Index(clay, builderName) == -1 {
		singleCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, clay)
		return singleCall, nil
	}
	series := ""
	for {
		builderIndex := strings.Index(clay, builderName)
		if builderIndex == -1 {
			break
		}
		htmlPart := clay[:builderIndex]
		if htmlPart != "" {
			htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, htmlPart)
			series += htmlCall + "\n"
			clay = strings.Replace(clay, htmlPart, "", 1)
		}
		endBuilderIndex := strings.Index(clay, ")")
		builderPart := clay[:endBuilderIndex+1]
		series += builderPart + "\n"
		clay = strings.Replace(clay, builderPart, "", 1)
	}
	if len(clay) > 0 {
		htmlCall := fmt.Sprintf("%s.WriteString(`%s`)", builderName, clay)
		series += htmlCall + "\n"
	}
	return series, nil
}

func WalkAllElementNodes(elm Element, fn func(sel *goquery.Selection) error) error {
	var potErr error
	elm.GetSelection().Find("*").Each(func(i int, s *goquery.Selection) {
		err := fn(s)
		if err != nil {
			potErr = err
			return
		}
	})
	if potErr != nil {
		return potErr
	}
	return nil
}

func GetElementPlaceholders(elm Element, allElements []Element) ([]Placeholder, error) {
	placeholders := make([]Placeholder, 0)
	for _, sibling := range allElements {
		err := WalkAllElementNodes(elm, func(sel *goquery.Selection) error {
			nodeName := goquery.NodeName(sel)
			nodeHtml, err := gqpp.NewHtmlFromSelection(sel)
			if err != nil {
				return err
			}
			if nodeName == strings.ToLower(sibling.GetAttr()) {
				place, err := NewPlaceholder(nodeHtml, sibling)
				if err != nil {
					return err
				}
				placeholders = append(placeholders, place)
			}
			return nil
		})
		if err != nil {
			panic(err)
		}
	}
	return placeholders, nil
}

func ReadComponentElementsFromFile(path string) ([]Element, error) {
	elms := make([]Element, 0)
	f, err := os.ReadFile(path)
	if err != nil {
		return elms, err
	}
	fStr := string(f)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(fStr))
	if err != nil {
		return elms, err
	}
	var potErr error
	doc.Find("*").Each(func(i int, sel *goquery.Selection) {
		_, exists := sel.Attr(KeyElementComponent)
		if exists {
			elm, err := NewElement(sel)
			if err != nil {
				potErr = err
				return
			}
			elms = append(elms, elm)
		}
	})
	if potErr != nil {
		return elms, potErr
	}
	return elms, nil
}

// ##==================================================================
type ElementComponent struct {
	Selection    *goquery.Selection
	Html         string
	Type         string
	Attr         string
	AttrParts    []string
	Name         string
	Props        []Prop
	Placeholders []Placeholder
}

func NewElementComponent(sel *goquery.Selection) (*ElementComponent, error) {
	elm := &ElementComponent{}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initName() },
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementComponent) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementComponent) GetParam() (Param, error) {
	return nil, nil
}
func (elm *ElementComponent) GetHtml() string        { return elm.Html }
func (elm *ElementComponent) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementComponent) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementComponent) GetType() string        { return elm.Type }
func (elm *ElementComponent) GetAttr() string        { return elm.Attr }
func (elm *ElementComponent) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementComponent) GetName() string        { return elm.Name }
func (elm *ElementComponent) GetProps() []Prop       { return elm.Props }
func (elm *ElementComponent) SetPlaceholders(placeholders []Placeholder) {
	elm.Placeholders = placeholders
}

func (elm *ElementComponent) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementComponent) initType() error {
	elm.Type = KeyElementComponent
	return nil
}

func (elm *ElementComponent) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementComponent) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementComponent)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementComponent, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementComponent) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementComponent) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================
type ElementFor struct {
	Selection    *goquery.Selection
	Html         string
	Type         string
	Attr         string
	AttrParts    []string
	Name         string
	Props        []Prop
	Placeholders []Placeholder
}

func NewElementFor(sel *goquery.Selection) (*ElementFor, error) {
	elm := &ElementFor{}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initName() },
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementFor) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementFor) GetParam() (Param, error) {
	parts := elm.GetAttrParts()
	iterItems := parts[2]
	if strings.Contains(iterItems, ".") {
		return nil, nil
	}
	iterType := parts[3]
	param, err := NewParam(iterItems, iterType)
	if err != nil {
		return nil, err
	}
	return param, nil
}
func (elm *ElementFor) GetHtml() string        { return elm.Html }
func (elm *ElementFor) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementFor) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementFor) GetType() string        { return elm.Type }
func (elm *ElementFor) GetAttr() string        { return elm.Attr }
func (elm *ElementFor) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementFor) GetName() string        { return elm.Name }
func (elm *ElementFor) GetProps() []Prop       { return elm.Props }
func (elm *ElementFor) SetPlaceholders(placeholders []Placeholder) {
	elm.Placeholders = placeholders
}

func (elm *ElementFor) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementFor) initType() error {
	elm.Type = KeyElementFor
	return nil
}

func (elm *ElementFor) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementFor) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementFor)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementFor, 4)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementFor) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementFor) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================
type ElementIf struct {
	Selection    *goquery.Selection
	Html         string
	Type         string
	Attr         string
	AttrParts    []string
	Name         string
	Props        []Prop
	Placeholders []Placeholder
}

func NewElementIf(sel *goquery.Selection) (*ElementIf, error) {
	elm := &ElementIf{}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initName() },
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementIf) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementIf) GetParam() (Param, error) {
	param, err := NewParam(elm.Attr, "bool")
	if err != nil {
		return nil, err
	}
	return param, nil
}
func (elm *ElementIf) GetHtml() string        { return elm.Html }
func (elm *ElementIf) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementIf) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementIf) GetType() string        { return elm.Type }
func (elm *ElementIf) GetAttr() string        { return elm.Attr }
func (elm *ElementIf) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementIf) GetName() string        { return elm.Name }
func (elm *ElementIf) GetProps() []Prop       { return elm.Props }
func (elm *ElementIf) SetPlaceholders(placeholders []Placeholder) {
	elm.Placeholders = placeholders
}

func (elm *ElementIf) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementIf) initType() error {
	elm.Type = KeyElementIf
	return nil
}

func (elm *ElementIf) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementIf) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementIf)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementIf, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementIf) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementIf) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================
type ElementElse struct {
	Selection *goquery.Selection
	Html      string
	Type      string
	Attr      string
	AttrParts []string
	Name      string
	Props     []Prop
}

func NewElementElse(sel *goquery.Selection) (*ElementElse, error) {
	elm := &ElementElse{}
	err := fungi.Process(
		func() error { return elm.initSelection(sel) },
		func() error { return elm.initType() },
		func() error { return elm.initHtml() },
		func() error { return elm.initAttr() },
		func() error { return elm.initName() },
		func() error { return elm.initProps() },
	)
	if err != nil {
		return nil, err
	}
	return elm, nil
}

func (elm *ElementElse) GetSelection() *goquery.Selection { return elm.Selection }
func (elm *ElementElse) GetParam() (Param, error) {
	param, err := NewParam(elm.Attr, "bool")
	if err != nil {
		return nil, err
	}
	return param, nil
}
func (elm *ElementElse) GetHtml() string        { return elm.Html }
func (elm *ElementElse) SetHtml(htmlStr string) { elm.Html = htmlStr }
func (elm *ElementElse) Print()                 { fmt.Println(elm.Html) }
func (elm *ElementElse) GetType() string        { return elm.Type }
func (elm *ElementElse) GetAttr() string        { return elm.Attr }
func (elm *ElementElse) GetAttrParts() []string { return elm.AttrParts }
func (elm *ElementElse) GetName() string        { return elm.Name }
func (elm *ElementElse) GetProps() []Prop       { return elm.Props }

func (elm *ElementElse) initSelection(sel *goquery.Selection) error {
	elm.Selection = sel
	return nil
}

func (elm *ElementElse) initType() error {
	elm.Type = KeyElementElse
	return nil
}

func (elm *ElementElse) initHtml() error {
	htmlStr, err := gqpp.NewHtmlFromSelection(elm.GetSelection())
	if err != nil {
		return err
	}
	elm.Html = htmlStr
	return nil
}

func (elm *ElementElse) initAttr() error {
	attr, err := gqpp.ForceElementAttr(elm.GetSelection(), KeyElementElse)
	if err != nil {
		return err
	}
	parts, err := gqpp.ForceElementAttrParts(elm.GetSelection(), KeyElementElse, 1)
	if err != nil {
		return err
	}
	elm.Attr = attr
	elm.AttrParts = parts
	return nil
}

func (elm *ElementElse) initName() error {
	elm.Name = fmt.Sprintf("%s:%s", elm.GetType(), elm.GetAttr())
	return nil
}

func (elm *ElementElse) initProps() error {
	props, err := GetElementProps(elm)
	if err != nil {
		return err
	}
	elm.Props = props
	return nil
}

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
