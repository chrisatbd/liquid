package liquid

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
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
	return query + ":" + t
}

type SimpleBindingStruct struct {
	Name string
	Id   int
}

func TestValueStruct_One(t *testing.T) {
	engine := NewEngine()

	har := testValueStruct{}
	har.F = 9

	template := `{{ m.F }}.{{ m.M1 }}`
	bindings := map[string]interface{}{
		"m": har,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "9.3", value)
}

func TestValueStruct_Two(t *testing.T) {
	engine := NewEngine()

	har := testValueStruct{
		F: 9,
	}

	template := `{% assign val = "c" %}{{ m.F }}.{{ m.M1 }}.{{ val }}`
	bindings := map[string]interface{}{
		"m": har,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "9.3.c", value)
}

func TestValueStruct_Three(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{
		PropertyOne: 9,
	}
	har.ArrayOne = append(har.ArrayOne, "one")
	har.ArrayOne = append(har.ArrayOne, "two")

	template := `{{ struct.PropertyOne }}.{{ struct.FunctionOne }} a text token {{ struct.ArrayOne[1] }}`
	bindings := map[string]interface{}{
		"struct": har,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "9.3 a text token two", value)
}

func TestValueStruct_Four(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{
		PropertyOne: 9,
	}
	har.ArrayOne = append(har.ArrayOne, "one")
	har.ArrayOne = append(har.ArrayOne, "two")

	template := `{{ struct.ArrayOne[1] }}`
	bindings := map[string]interface{}{
		"struct": har,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "two", value)
}

func TestValueStruct_Five(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{
		PropertyOne: 9,
	}

	template := `{{ struct.FunctionTwo("chris") }}`
	bindings := map[string]interface{}{
		"struct": har,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "chris3", value)
}

func TestValueStruct_Six(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{
		PropertyOne: 9,
	}

	template := `{{ struct.FunctionThree("chris",57) }} | {{ myStringProp }} | {{ myIntProp }}`

	bindings := map[string]interface{}{
		"struct":       har,
		"myStringProp": "ag",
		"myIntProp":    1,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "chris:57 | ag | 1", value)
}

func TestValueStruct_Seven(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{
		PropertyOne: 9,
	}

	template := `{{ struct.FunctionThree("chris",57) }} | {{ myStringProp }} | {{ myIntProp }} | {{ struct.FunctionThree(myStringProp,myIntProp) }}`

	bindings := map[string]interface{}{
		"struct":       har,
		"myStringProp": "ag",
		"myIntProp":    1,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "chris:57 | ag | 1 | ag:1", value)
}

func TestValueStruct_Eight(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{
		PropertyOne: 9,
	}

	template := `{{ struct.FunctionThree("chris",57) }} | {{ struct.FunctionThree(myStringProp,myIntProp) }} | {{ struct.FunctionThree(sbs.Name,sbs.Id) }}`

	sbs := SimpleBindingStruct{"joe jones", 72}

	bindings := map[string]interface{}{
		"struct":       har,
		"myStringProp": "ag",
		"myIntProp":    1,
		"sbs":          sbs,
	}

	value, err := engine.ParseAndRenderString(template, bindings)

	require.NoError(t, err)
	require.Equal(t, "chris:57 | ag:1 | joe jones:72", value)

}

func TestValueStruct_TypeConversionFail(t *testing.T) {

	engine := NewEngine()

	har := testValueStructTwo{
		PropertyOne: 9,
	}

	// FunctionThree expects a string and an int
	template := `empty: {{ struct.FunctionThree(22,57) }}. has value: {{ struct.FunctionThree("22",57) }}`

	bindings := map[string]interface{}{
		"struct": har,
	}

	//ok, should this be an 'empty' or should we put up an error ?
	value, err := engine.ParseAndRenderString(template, bindings)

	//require.Error(t, err)
	require.NoError(t, err)
	require.Equal(t, "empty: . has value: 22:57", value)
}
