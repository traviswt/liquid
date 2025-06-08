package liquid

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"strconv"
	"strings"
	"testing"
)

var emptyBindings = map[string]any{}

// There's a lot more tests in the filters and tags sub-packages.
// This collects a minimal set for testing end-to-end.
var liquidTests = []struct{ in, expected string }{
	{`{{ page.title }}`, "Introduction"},
	{`{% if x %}true{% endif %}`, "true"},
	{`{{ "upper" | upcase }}`, "UPPER"},
}

var testBindings = map[string]any{
	"x":  123,
	"ar": []string{"first", "second", "third"},
	"page": map[string]any{
		"title": "Introduction",
	},
}

func TestEngine_ParseAndRenderString(t *testing.T) {
	engine := NewEngine()
	for i, test := range liquidTests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			out, err := engine.ParseAndRenderString(test.in, testBindings)
			require.NoErrorf(t, err, test.in)
			require.Equalf(t, test.expected, out, test.in)
		})
	}
}

func TestBasicEngine_ParseAndRenderString(t *testing.T) {
	engine := NewBasicEngine()

	t.Run("1", func(t *testing.T) {
		test := liquidTests[0]
		out, err := engine.ParseAndRenderString(test.in, testBindings)
		require.NoErrorf(t, err, test.in)
		require.Equalf(t, test.expected, out, test.in)
	})

	for i, test := range liquidTests[1:] {
		t.Run(strconv.Itoa(i+2), func(t *testing.T) {
			out, err := engine.ParseAndRenderString(test.in, testBindings)
			require.Errorf(t, err, test.in)
			require.Emptyf(t, out, test.in)
		})
	}
}

type capWriter struct {
	bytes.Buffer
}

func (c *capWriter) Write(bs []byte) (int, error) {
	return c.Buffer.Write([]byte(strings.ToUpper(string(bs))))
}

func TestEngine_ParseAndFRender(t *testing.T) {
	engine := NewEngine()
	for i, test := range liquidTests {
		t.Run(strconv.Itoa(i+1), func(t *testing.T) {
			wr := capWriter{}
			err := engine.ParseAndFRender(&wr, []byte(test.in), testBindings)
			require.NoErrorf(t, err, test.in)
			require.Equalf(t, strings.ToUpper(test.expected), wr.String(), test.in)
		})
	}
}

func TestEngine_ParseAndRenderString_ptr_to_hash(t *testing.T) {
	params := map[string]any{
		"message": &map[string]any{
			"Text":       "hello",
			"jsonNumber": json.Number("123"),
		},
	}
	engine := NewEngine()
	template := "{{ message.Text }} {{message.jsonNumber}}"
	str, err := engine.ParseAndRenderString(template, params)
	require.NoError(t, err)
	require.Equal(t, "hello 123", str)
}

type testStruct struct{ Text string }

func TestEngine_ParseAndRenderString_struct(t *testing.T) {
	params := map[string]any{
		"message": testStruct{
			Text: "hello",
		},
	}
	engine := NewEngine()
	template := "{{ message.Text }}"
	str, err := engine.ParseAndRenderString(template, params)
	require.NoError(t, err)
	require.Equal(t, "hello", str)
}

func TestEngine_ParseAndRender_errors(t *testing.T) {
	_, err := NewEngine().ParseAndRenderString("{{ syntax error }}", emptyBindings)
	require.Error(t, err)
	_, err = NewEngine().ParseAndRenderString("{% if %}", emptyBindings)
	require.Error(t, err)
	_, err = NewEngine().ParseAndRenderString("{% undefined_tag %}", emptyBindings)
	require.Error(t, err)
	_, err = NewEngine().ParseAndRenderString("{% a | undefined_filter %}", emptyBindings)
	require.Error(t, err)
}

func BenchmarkEngine_Parse(b *testing.B) {
	engine := NewEngine()
	buf := new(bytes.Buffer)
	for range 1000 {
		_, err := io.WriteString(buf, `if{% if true %}true{% elsif %}elsif{% else %}else{% endif %}`)
		require.NoError(b, err)
		_, err = io.WriteString(buf, `loop{% for item in array %}loop{% break %}{% endfor %}`)
		require.NoError(b, err)
		_, err = io.WriteString(buf, `case{% case value %}{% when a %}{% when b %{% endcase %}`)
		require.NoError(b, err)
		_, err = io.WriteString(buf, `expr{{ a and b }}{{ a add: b }}`)
		require.NoError(b, err)
	}
	s := buf.Bytes()
	b.ResetTimer()
	for range b.N {
		_, err := engine.ParseTemplate(s)
		require.NoError(b, err)
	}
}

func TestEngine_ParseTemplateAndCache(t *testing.T) {
	// Given two templates...
	templateA := []byte("Foo")
	templateB := []byte(`{% include "template_a.html" %}, Bar`)

	// Cache the first
	eng := NewEngine()
	_, err := eng.ParseTemplateAndCache(templateA, "template_a.html", 1)
	require.NoError(t, err)

	// ...and execute the second.
	result, err := eng.ParseAndRender(templateB, Bindings{})
	require.NoError(t, err)
	require.Equal(t, "Foo, Bar", string(result))
}

func TestEngine_ListFilters(t *testing.T) {
	eng := NewEngine()
	spew.Dump(eng.ListFilters())
}

func TestEngine_ListTags(t *testing.T) {
	eng := NewEngine()
	spew.Dump(eng.ListTags())
}

type MockTemplateStore struct{}

func (tl *MockTemplateStore) ReadTemplate(filename string) ([]byte, error) {
	template := []byte(fmt.Sprintf("Message Text: {{ message.Text }} from: %v.", filename))
	return template, nil
}

func Test_template_store(t *testing.T) {
	template := []byte(`{% include "template.liquid" %}`)
	mockstore := &MockTemplateStore{}
	params := map[string]any{
		"message": testStruct{
			Text: "filename",
		},
	}
	engine := NewEngine()
	engine.RegisterTemplateStore(mockstore)
	out, _ := engine.ParseAndRenderString(string(template), params)
	require.Equal(t, "Message Text: filename from: template.liquid.", out)
}

func Test_Base64RoundTrip(t *testing.T) {
	a := assert.New(t)
	template := `

<h1>{{ page.title | base64_encode }}</h1>
<p>{{ description | base64_decode }}</p>

`
	bindings := map[string]any{
		"page": map[string]string{
			"title": "The Best Page Ever!",
		},
		"description": "dGhpcyBwYWdlIGlzIHRoZSBiZXN0IHBhZ2UgZXZlciE=",
	}
	mockstore := &MockTemplateStore{}
	engine := NewEngine()
	engine.RegisterTemplateStore(mockstore)
	result, err := engine.ParseAndRenderString(template, bindings)
	if !a.NoError(err) {
		return
	}
	if !a.NotEmpty(result) {
		return
	}
	spew.Dump(result)
	if !a.Equal("\n\n<h1>VGhlIEJlc3QgUGFnZSBFdmVyIQ==</h1>\n<p>this page is the best page ever!</p>\n\n", result) {
		return
	}
}

func Test_Jq(t *testing.T) {
	a := assert.New(t)
	template := `

<h1>{{ page.title | base64_encode }}</h1>
<p>{{ description | base64_decode | jq: '.values.testb' }}</p>

`
	bindings := map[string]any{
		"page": map[string]string{
			"title": "The Best Page Ever!",
		},
		"description": "eyJ2YWx1ZXMiOnsidGVzdGEiOiJhbHBoYSIsInRlc3RiIjoiYmV0YSJ9fQ==",
	}
	mockstore := &MockTemplateStore{}
	engine := NewEngine()
	engine.RegisterTemplateStore(mockstore)
	result, err := engine.ParseAndRenderString(template, bindings)
	if !a.NoError(err) {
		return
	}
	if !a.NotEmpty(result) {
		return
	}
	if !a.Equal("\n\n<h1>VGhlIEJlc3QgUGFnZSBFdmVyIQ==</h1>\n<p>beta</p>\n\n", result) {
		return
	}
	spew.Dump(result)
}

func Test_ReadFile(t *testing.T) {
	a := assert.New(t)
	template := `

<h1>{{ page.title }}</h1>
<p>{{ description | base64_decode | jq: '.values.testb' }}</p>
<p>{{ './LICENSE' | load_file | base64_encode }}</p>

`
	bindings := map[string]any{
		"page": map[string]string{
			"title": "The Best Page Ever!",
		},
		"description": "eyJ2YWx1ZXMiOnsidGVzdGEiOiJhbHBoYSIsInRlc3RiIjoiYmV0YSJ9fQ==",
	}
	mockstore := &MockTemplateStore{}
	engine := NewEngine()
	engine.RegisterTemplateStore(mockstore)
	result, err := engine.ParseAndRenderString(template, bindings)
	if !a.NoError(err) {
		return
	}
	if !a.NotEmpty(result) {
		return
	}
	if !a.Equal("\n\n<h1>The Best Page Ever!</h1>\n<p>beta</p>\n<p>TUlUIExpY2Vuc2UKCkNvcHlyaWdodCAoYykgMjAxNyBPbGl2ZXIgU3RlZWxlCgpQZXJtaXNzaW9uIGlzIGhlcmVieSBncmFudGVkLCBmcmVlIG9mIGNoYXJnZSwgdG8gYW55IHBlcnNvbiBvYnRhaW5pbmcgYSBjb3B5Cm9mIHRoaXMgc29mdHdhcmUgYW5kIGFzc29jaWF0ZWQgZG9jdW1lbnRhdGlvbiBmaWxlcyAodGhlICJTb2Z0d2FyZSIpLCB0byBkZWFsCmluIHRoZSBTb2Z0d2FyZSB3aXRob3V0IHJlc3RyaWN0aW9uLCBpbmNsdWRpbmcgd2l0aG91dCBsaW1pdGF0aW9uIHRoZSByaWdodHMKdG8gdXNlLCBjb3B5LCBtb2RpZnksIG1lcmdlLCBwdWJsaXNoLCBkaXN0cmlidXRlLCBzdWJsaWNlbnNlLCBhbmQvb3Igc2VsbApjb3BpZXMgb2YgdGhlIFNvZnR3YXJlLCBhbmQgdG8gcGVybWl0IHBlcnNvbnMgdG8gd2hvbSB0aGUgU29mdHdhcmUgaXMKZnVybmlzaGVkIHRvIGRvIHNvLCBzdWJqZWN0IHRvIHRoZSBmb2xsb3dpbmcgY29uZGl0aW9uczoKClRoZSBhYm92ZSBjb3B5cmlnaHQgbm90aWNlIGFuZCB0aGlzIHBlcm1pc3Npb24gbm90aWNlIHNoYWxsIGJlIGluY2x1ZGVkIGluIGFsbApjb3BpZXMgb3Igc3Vic3RhbnRpYWwgcG9ydGlvbnMgb2YgdGhlIFNvZnR3YXJlLgoKVEhFIFNPRlRXQVJFIElTIFBST1ZJREVEICJBUyBJUyIsIFdJVEhPVVQgV0FSUkFOVFkgT0YgQU5ZIEtJTkQsIEVYUFJFU1MgT1IKSU1QTElFRCwgSU5DTFVESU5HIEJVVCBOT1QgTElNSVRFRCBUTyBUSEUgV0FSUkFOVElFUyBPRiBNRVJDSEFOVEFCSUxJVFksCkZJVE5FU1MgRk9SIEEgUEFSVElDVUxBUiBQVVJQT1NFIEFORCBOT05JTkZSSU5HRU1FTlQuIElOIE5PIEVWRU5UIFNIQUxMIFRIRQpBVVRIT1JTIE9SIENPUFlSSUdIVCBIT0xERVJTIEJFIExJQUJMRSBGT1IgQU5ZIENMQUlNLCBEQU1BR0VTIE9SIE9USEVSCkxJQUJJTElUWSwgV0hFVEhFUiBJTiBBTiBBQ1RJT04gT0YgQ09OVFJBQ1QsIFRPUlQgT1IgT1RIRVJXSVNFLCBBUklTSU5HIEZST00sCk9VVCBPRiBPUiBJTiBDT05ORUNUSU9OIFdJVEggVEhFIFNPRlRXQVJFIE9SIFRIRSBVU0UgT1IgT1RIRVIgREVBTElOR1MgSU4gVEhFClNPRlRXQVJFLgo=</p>\n\n", result) {
		return
	}
	spew.Dump(result)
}
