package liquid

import (
	"strconv"
	"testing"
)

type testValueStruct struct {
	F int
}

func (tv testValueStruct) M1() int                { return 3 }
func (tv testValueStruct) M2(query string) string { return query + "3" }

type testValueStructTwo struct {
	PropertyOne int
	ArrayOne    []string
}

func (tv testValueStructTwo) FunctionOne() int                { return 3 }
func (tv testValueStructTwo) FunctionTwo(query string) string { return query + "3" }
func (tv testValueStructTwo) FunctionThree(query string, val int) string {
	t := strconv.Itoa(val)
	return query + " " + t
}

func TestValueStruct_One(t *testing.T) {
	engine := NewEngine()

	har := testValueStruct{}
	har.F = 9

	template := `{{ m.F }}.{{ m.M1 }}`
	bindings := map[string]interface{}{
		"m": har,
	}

	out, err := engine.ParseAndRenderString(template, bindings)
	if err != nil {
		t.Log(err)
	}
	t.Log(out)

	_ = engine
	_ = har
}

func TestValueStruct_Two(t *testing.T) {
	engine := NewEngine()

	har := testValueStruct{}
	har.F = 9

	template := `{% assign val = "c" %}{{ m.F }}.{{ m.M1 }}.{{ val }}`
	bindings := map[string]interface{}{
		"m": har,
	}

	out, err := engine.ParseAndRenderString(template, bindings)
	if err != nil {
		t.Log(err)
		t.Fail()
	}
	t.Log(out)

	_ = engine
	_ = har
}

func TestValueStruct_Three(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{}
	har.PropertyOne = 9
	har.ArrayOne = append(har.ArrayOne, "one")
	har.ArrayOne = append(har.ArrayOne, "two")

	template := `{{ struct.PropertyOne }}.{{ struct.FunctionOne }} a text token {{ struct.ArrayOne[1] }}`
	bindings := map[string]interface{}{
		"struct": har,
	}

	out, err := engine.ParseAndRenderString(template, bindings)
	if err != nil {
		t.Log(err)
	}
	t.Log(out)

	_ = engine
	_ = har
}

func TestValueStruct_Four(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{}
	har.PropertyOne = 9
	har.ArrayOne = append(har.ArrayOne, "one")
	har.ArrayOne = append(har.ArrayOne, "two")

	template := `{{ struct.ArrayOne[1] }}`
	bindings := map[string]interface{}{
		"struct": har,
	}

	out, err := engine.ParseAndRenderString(template, bindings)
	if err != nil {
		t.Log(err)
	}

	t.Log(out)

	if out != "two" {
		t.Fail()
	}

	_ = engine
	_ = har
}

func TestValueStruct_Five(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{}
	har.PropertyOne = 9

	template := `{{ struct.FunctionTwo("chris") }}`
	bindings := map[string]interface{}{
		"struct": har,
	}

	out, err := engine.ParseAndRenderString(template, bindings)
	if err != nil {
		t.Log(err)
	}
	t.Log(out)

	_ = engine
	_ = har
}

func TestValueStruct_Six(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{}
	har.PropertyOne = 9

	template := `{{ struct.FunctionThree("chris",57) }}`
	bindings := map[string]interface{}{
		"struct": har,
	}

	out, err := engine.ParseAndRenderString(template, bindings)
	if err != nil {
		t.Log(err)
	}
	t.Log(out)

	_ = engine
	_ = har
}


