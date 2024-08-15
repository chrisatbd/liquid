package tags

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/chrisatbd/liquid/parser"
	"github.com/chrisatbd/liquid/render"
	"github.com/stretchr/testify/require"
)

var renderTestBindings = map[string]interface{}{
	"test": true,
	"var":  "value",
}

func TestRenderTag(t *testing.T) {
	config := render.NewConfig()
	loc := parser.SourceLoc{Pathname: "testdata/include_source.html", LineNo: 1}
	AddStandardTags(config)

	// basic functionality
	//root, err := config.Compile(`{% render "include_target.html", product: "literal product", product2: actualProduct %}`, loc)
	root, err := config.Compile(`{% render "testdata/include_target.html" %}`, loc)
	require.NoError(t, err)
	buf := new(bytes.Buffer)
	err = render.Render(root, buf, renderTestBindings, config)
	require.NoError(t, err)
	require.Equal(t, "include target", strings.TrimSpace(buf.String()))

	// tag and variable
	root, err = config.Compile(`{% include "include_target_2.html" %}`, loc)
	require.NoError(t, err)
	buf = new(bytes.Buffer)
	err = render.Render(root, buf, renderTestBindings, config)
	require.NoError(t, err)
	require.Equal(t, "test value", strings.TrimSpace(buf.String()))

	// errors
	root, err = config.Compile(`{% include 10 %}`, loc)
	require.NoError(t, err)
	err = render.Render(root, io.Discard, renderTestBindings, config)
	require.Error(t, err)
	require.Contains(t, err.Error(), "requires a string")
}

func TestRenderTag_file_not_found_error(t *testing.T) {
	config := render.NewConfig()
	loc := parser.SourceLoc{Pathname: "testdata/include_source.html", LineNo: 1}
	AddStandardTags(config)

	// See the comment in TestIncludeTag_file_not_found_error.
	root, err := config.Compile(`{% render "missing_file.html" %}`, loc)
	require.NoError(t, err)
	err = render.Render(root, io.Discard, renderTestBindings, config)
	require.Error(t, err)
	require.True(t, os.IsNotExist(err.Cause()))
}

func TestRenderTag_cached_value_handling(t *testing.T) {
	// skip for now and then come back at some future point
	t.Skip()
	config := render.NewConfig()
	// missing-file.html does not exist in the testdata directory.
	config.Cache["testdata/missing-file.html"] = []byte("include-content")
	config.Cache["testdata\\missing-file.html"] = []byte("include-content")
	loc := parser.SourceLoc{Pathname: "testdata/include_source.html", LineNo: 1}
	AddStandardTags(config)

	root, err := config.Compile(`{% render "missing-file.html" %}`, loc)
	require.NoError(t, err)
	buf := new(bytes.Buffer)
	err = render.Render(root, buf, includeTestBindings, config)
	require.NoError(t, err)
	require.Equal(t, "include-content", strings.TrimSpace(buf.String()))
}

func TestRenderTagWithStringLiteralParameter(t *testing.T) {
	config := render.NewConfig()
	loc := parser.SourceLoc{Pathname: "testdata/render_with_parameters.liquid", LineNo: 1}
	AddStandardTags(config)

	// basic functionality
	root, err := config.Compile(`{% render "testdata/render_with_parameters.liquid, product: 'literal product'" %}`, loc)
	require.NoError(t, err)
	buf := new(bytes.Buffer)
	err = render.Render(root, buf, renderTestBindings, config)
	require.NoError(t, err)
	require.Equal(t, "include target literal product", strings.TrimSpace(buf.String()))
}
