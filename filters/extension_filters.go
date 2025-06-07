package filters

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/itchyny/gojq"
)

func AddExtensionFilters(fd FilterDictionary) {
	fd.AddFilter("base64_encode", base64Encode)
	fd.AddFilter("base64_decode", base64Decode)
	fd.AddFilter("jq", jq)
}

func base64Encode(str string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(str)), nil
}

func base64Decode(str string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(str)
	return string(b), err
}

func jq(spec interface{}, filter string) (interface{}, error) {
	q, err := gojq.Parse(filter)
	if err != nil {
		return nil, fmt.Errorf("the jq filter string '%s' is not valid: %s", filter, err.Error())
	}
	if q != nil {
		code, err := gojq.Compile(q)
		if err != nil {
			return nil, fmt.Errorf("the jq filter string '%s' is not valid: %s", filter, err.Error())
		}
		res, err := applyJqFilters(spec, []*gojq.Code{code})
		if err != nil {
			return nil, err
		}
		if res == nil || len(res) == 0 {
			return nil, nil
		}
		if len(res) == 1 {
			return res[0], nil
		}
		return res, nil
	}
	return nil, nil
}

func applyJqFilters(spec interface{}, filters []*gojq.Code) ([]any, error) {
	if filters == nil || len(filters) == 0 {
		return nil, nil
	}
	input, err := toMap(spec)
	if err != nil {
		return nil, err
	}
	var values []any
	for _, code := range filters {
		if code != nil {
			it := code.Run(input)
			if it != nil {
				for {
					value, ok := it.Next()
					if !ok {
						break
					}
					if value != nil {
						if _, ok := value.(error); ok {
							break
						}
						values = append(values, value)
					}
				}
			}
		}
	}
	return values, nil
}

func toMap(v interface{}) (map[string]interface{}, error) {
	if v == nil {
		return nil, nil
	}
	var (
		inrec []byte
		err   error
	)
	if gojq.TypeOf(v) == "string" {
		inrec = []byte(v.(string))
	} else if gojq.TypeOf(v) == "null" {
		return nil, nil
	} else {
		inrec, err = json.Marshal(v)
		if err != nil {
			return nil, err
		}
	}
	var out map[string]interface{}
	if err := json.Unmarshal(inrec, &out); err != nil {
		return nil, err
	}
	return out, nil
}
