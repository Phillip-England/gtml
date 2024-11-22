package gtml

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/phillip-england/fungi"
	"github.com/phillip-england/gqpp"
	"github.com/phillip-england/purse"
)

// ##==================================================================
type Placeholder interface {
	Print()
	GetFoundAs() string
	GetPointingTo() Element
}

func NewPlaceholder(foundAsHtml string, pointingTo Element) (Placeholder, error) {
	if pointingTo.GetType() != KeyElementComponent {
		return nil, fmt.Errorf("a placeholder must point to a valid _component element: %s", pointingTo.GetHtml())
	}
	place, err := NewPlaceholderComponent(foundAsHtml, pointingTo)
	if err != nil {
		return nil, err
	}
	return place, nil
}

// ##==================================================================
type PlaceholderComponent struct {
	Name              string
	NodeName          string
	Html              string
	PointingTo        Element
	Params            []Param
	Attrs             []Attr
	FuncParamSlice    []string
	ComponentFuncCall string
}

func NewPlaceholderComponent(foundAsHtml string, pointingTo Element) (*PlaceholderComponent, error) {
	place := &PlaceholderComponent{
		PointingTo: pointingTo,
	}
	err := fungi.Process(
		func() error { return place.initNodeName(foundAsHtml) },
		func() error { return place.initName() },
		func() error { return place.initHtml(foundAsHtml) },
		func() error { return place.initParamNames() },
		func() error { return place.initAttrs() },
		func() error { return place.initFuncParamSlice() },
		func() error { return place.initComponentFuncCall() },
	)
	if err != nil {
		return nil, err
	}
	for _, attr := range place.Attrs {
		fmt.Println(attr.GetType())
	}
	return place, nil
}

func (place *PlaceholderComponent) initNodeName(foundAsHtml string) error {
	sel, err := gqpp.NewSelectionFromStr(foundAsHtml)
	if err != nil {
		return err
	}
	nodeName := goquery.NodeName(sel)
	place.NodeName = nodeName
	return nil
}

func (place *PlaceholderComponent) initName() error {
	nameAttr := place.PointingTo.GetAttr()
	place.Name = nameAttr
	return nil
}

func (place *PlaceholderComponent) initHtml(foundAsHtml string) error {
	htmlStr := purse.ReplaceFirstInstanceOf(foundAsHtml, place.NodeName, place.Name)
	htmlStr = purse.ReplaceLastInstanceOf(htmlStr, place.NodeName, place.Name)
	place.Html = htmlStr
	return nil
}

func (place *PlaceholderComponent) initParamNames() error {
	params, err := GetElementParams(place.PointingTo)
	if err != nil {
		return err
	}
	for _, param := range params {
		place.Params = append(place.Params, param)
	}
	return nil
}

func (place *PlaceholderComponent) initAttrs() error {
	sel, err := gqpp.NewSelectionFromStr(place.Html)
	if err != nil {
		return err
	}
	for _, node := range sel.Nodes {
		for _, attr := range node.Attr {
			attrType, err := NewAttr(attr.Key, attr.Val)
			if err != nil {
				return err
			}
			place.Attrs = append(place.Attrs, attrType)
		}
	}
	return nil
}

func (place *PlaceholderComponent) initFuncParamSlice() error {
	funcParamSlice := make([]string, 0)

	place.FuncParamSlice = funcParamSlice
	return nil
}

func (place *PlaceholderComponent) initComponentFuncCall() error {
	paramStr := strings.Join(place.FuncParamSlice, ", ")
	call := fmt.Sprintf("%s(%s)", place.Name, paramStr)
	place.ComponentFuncCall = call
	return nil
}
func (place *PlaceholderComponent) Print()                 { fmt.Println(place.Html) }
func (place *PlaceholderComponent) GetFoundAs() string     { return place.Html }
func (place *PlaceholderComponent) GetPointingTo() Element { return place.PointingTo }

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================

// ##==================================================================
