package tags

import (
	"io"
	_ "path/filepath"
	"strconv"
	"strings"

	"github.com/chrisatbd/liquid/render"
)

// "subtemplate/path, var1:val1, var2: val2, var3: val3, ..."

// 1. Any argument value with single or double quotes is interpreted as a string literal
// 2. Any argument value without quotes is first checked as a global context binding and the interpreted as either a string literal, floating point, integer, or boolean
// 3. Any argument value with digits only and within the range of the system's integer max/min size but no decimal points is interpreted as an int32/int64
// 4. Any argument value with digits only and containing a decimal point and within the rnage of the system's floating point max/min size is interpreted and parsed as a floating point

func renderTag(source string) (func(io.Writer, render.Context) error, error) {
	return func(w io.Writer, ctx render.Context) error {

		//NOTECJH:  So the trick lies right here in the evaluate string
		//what do we want to do with the parser so that we can get back the template
		//file name as well as the parameters in some sort of a key value pair.
		//1st pass will be that args come in as one big string and we parse here.
		//then we will look at moving to the actual template parser.  Just not
		//100% sure how deep it is going to cut.
		//args is for some debugging
		args := ctx.TagArgs()
		value, err := ctx.EvaluateString(args)
		if err != nil {
			return err
		}
		rel, ok := value.(string)
		if !ok {
			return ctx.Errorf("render requires a string argument; got %v", value)
		}

		pathAndVars := strings.Split(rel, ",")
		if len(pathAndVars) == 0 {
			return ctx.Errorf("invalidTagArgs")
		}

		localContextBindings := make(map[string]interface{})

		for _, varPair := range pathAndVars[1:] {
			if !strings.Contains(varPair, ":") {
				return ctx.Errorf("invalidTagArgs")
			}
			parts := strings.Split(varPair, ":")
			if len(parts) != 2 {
				return ctx.Errorf("invalidTagArgs")
			}
	
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
	
			localContextBindings[key] = getAndConvertValue(val, ctx)
		}

		file := strings.TrimSpace(pathAndVars[0])

		//filename := filepath.Join(filepath.Dir(ctx.SourceFile()))
		s, err := ctx.RenderFile(file, localContextBindings)
		if err != nil {
			return err
		}
		_, err = io.WriteString(w, s)
		return err
	}, nil
}


// will either interpret the value as a pre-existing variable value, string literal, int/float, bool
func getAndConvertValue(value string, ctx render.Context) interface{} {
	if isVarValue(value, ctx) {
		return ctx.Get(value)
	}

	if isInt(value) {
		val, _ := strconv.Atoi(value)
		return val
	}

	if isFloat(value) {
		val, _ := strconv.ParseFloat(value, 64)
		return val
	}

	if isBool(value) {
		if value == "true" || value == "True" {
			return true
		} else {
			return false
		}
	}

	if isStringLiteral(value) {
		return value[1 : len(value)-1]
	}

	return value
}

// checks if argument is a pre-existing variable value
func isVarValue(value string, ctx render.Context) bool {
	if ctx != nil {
		_, ok := ctx.Bindings()[value]
		return ok
	}

	return false
}

func isInt(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}

func isFloat(value string) bool {
	_, err := strconv.ParseFloat(value, 64)
	return err == nil
}

func isBool(value string) bool {
	return value == "true" || value == "false" || value == "True" || value == "False"
}


func isStringLiteral(value string) bool {
	if len(value) < 2 {
		return false
	}
	firstChar := value[0]
	lastChar := value[len(value)-1]

	return !isInt(value) && !isFloat(value) && (firstChar == '"' && lastChar == '"') || (firstChar == '\'' && lastChar == '\'')
}
