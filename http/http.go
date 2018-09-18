package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/skylark"
	"github.com/google/skylark/skylarkstruct"
	"github.com/qri-io/starlib/util"
)

// ModuleName defines the expected name for this Module when used
// in skylark's load() function, eg: load('http.sky', 'http')
const ModuleName = "http.sky"

var (
	// Client is the http client used to create the http module. override with
	// a custom client before calling LoadModule
	Client = http.DefaultClient
	// Guard is a global RequestGuard used in LoadModule. override with a custom
	// implementation before calling LoadModule
	Guard RequestGuard
)

// LoadModule creates an http Module
func LoadModule() (skylark.StringDict, error) {
	var m = &Module{cli: Client}
	if Guard != nil {
		m.rg = Guard
	}
	ns := skylark.StringDict{
		"http": m.Struct(),
	}
	return ns, nil
}

// RequestGuard controls access to http by checking before making requests
// if Allowed returns an error the request will be denied
type RequestGuard interface {
	Allowed(req *http.Request) error
	// RequestCompleted(res *http.Response)
}

// Module joins http tools to a dataset, allowing dataset
// to follow along with http requests
type Module struct {
	cli *http.Client
	rg  RequestGuard
}

// Struct returns this module's methods as a skylark Struct
func (m *Module) Struct() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(skylarkstruct.Default, m.StringDict())
}

// StringDict returns all module methods in a skylark.StringDict
func (m *Module) StringDict() skylark.StringDict {
	return skylark.StringDict{
		"get":     skylark.NewBuiltin("get", m.reqMethod("get")),
		"put":     skylark.NewBuiltin("put", m.reqMethod("put")),
		"post":    skylark.NewBuiltin("post", m.reqMethod("post")),
		"delete":  skylark.NewBuiltin("delete", m.reqMethod("delete")),
		"patch":   skylark.NewBuiltin("patch", m.reqMethod("patch")),
		"options": skylark.NewBuiltin("options", m.reqMethod("options")),
	}
}

// reqMethod is a factory function for generating skylark builtin functions for different http request methods
func (m *Module) reqMethod(method string) func(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	return func(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
		var (
			urlv     skylark.String
			params   = &skylark.Dict{}
			headers  = &skylark.Dict{}
			data     = &skylark.Dict{}
			auth     skylark.Tuple
			jsondata skylark.Value
		)

		if err := skylark.UnpackArgs(method, args, kwargs, "url", &urlv, "params?", &params, "headers", &headers, "data", &data, "json", &jsondata, "auth", &auth); err != nil {
			return nil, err
		}

		rawurl, err := util.AsString(urlv)
		if err != nil {
			return nil, err
		}
		if err = setQueryParams(&rawurl, params); err != nil {
			return nil, err
		}

		req, err := http.NewRequest(strings.ToUpper(method), rawurl, nil)
		if err != nil {
			return nil, err
		}
		if m.rg != nil {
			if err := m.rg.Allowed(req); err != nil {
				return nil, err
			}
		}

		if err = setHeaders(req, headers); err != nil {
			return nil, err
		}
		if err = setAuth(req, auth); err != nil {
			return nil, err
		}
		if err = setBody(req, data, jsondata); err != nil {
			return nil, err
		}

		res, err := m.cli.Do(req)
		if err != nil {
			return nil, err
		}

		r := &Response{*res}
		return r.Struct(), nil
	}
}

func setQueryParams(rawurl *string, params *skylark.Dict) error {
	keys := params.Keys()
	if len(keys) == 0 {
		return nil
	}

	u, err := url.Parse(*rawurl)
	if err != nil {
		return err
	}

	q := u.Query()
	for _, key := range keys {
		keystr, err := util.AsString(key)
		if err != nil {
			return err
		}

		val, _, err := params.Get(key)
		if err != nil {
			return err
		}
		if val.Type() != "string" {
			return fmt.Errorf("expected param value for key '%s' to be a string. got: '%s'", key, val.Type())
		}
		valstr, err := util.AsString(val)
		if err != nil {
			return err
		}

		q.Set(keystr, valstr)
	}

	u.RawQuery = q.Encode()
	*rawurl = u.String()
	return nil
}

func setAuth(req *http.Request, auth skylark.Tuple) error {
	if len(auth) == 0 {
		return nil
	} else if len(auth) == 2 {
		username, err := util.AsString(auth[0])
		if err != nil {
			return fmt.Errorf("parsing auth username string: %s", err.Error())
		}
		password, err := util.AsString(auth[1])
		if err != nil {
			return fmt.Errorf("parsing auth password string: %s", err.Error())
		}
		req.SetBasicAuth(username, password)
		return nil
	}
	return fmt.Errorf("expected two values for auth params tuple")
}

func setHeaders(req *http.Request, headers *skylark.Dict) error {
	keys := headers.Keys()
	if len(keys) == 0 {
		return nil
	}

	for _, key := range keys {
		keystr, err := util.AsString(key)
		if err != nil {
			return err
		}

		val, _, err := headers.Get(key)
		if err != nil {
			return err
		}
		if val.Type() != "string" {
			return fmt.Errorf("expected param value for key '%s' to be a string. got: '%s'", key, val.Type())
		}
		valstr, err := util.AsString(val)
		if err != nil {
			return err
		}

		req.Header.Add(keystr, valstr)
	}

	return nil
}

func setBody(req *http.Request, data *skylark.Dict, jsondata skylark.Value) error {
	if jsondata != nil && jsondata.String() != "" {
		req.Header.Set("Content-Type", "application/json")

		v, err := util.Unmarshal(jsondata)
		if err != nil {
			return err
		}
		data, err := json.Marshal(v)
		if err != nil {
			return err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(data))
	}

	if data.Len() > 0 {
		req.Header.Set("Content-Type", "multipart/form-data")

		if req.Form == nil {
			req.Form = url.Values{}
		}
		for _, key := range data.Keys() {
			keystr, err := util.AsString(key)
			if err != nil {
				return err
			}

			val, _, err := data.Get(key)
			if err != nil {
				return err
			}
			if val.Type() != "string" {
				return fmt.Errorf("expected param value for key '%s' to be a string. got: '%s'", key, val.Type())
			}
			valstr, err := util.AsString(val)
			if err != nil {
				return err
			}

			req.Form.Add(keystr, valstr)
		}
	}

	return nil
}

// Response represents an HTTP response, wrapping a go http.Response with
// skylark methods
type Response struct {
	http.Response
}

// Struct turns a response into a *skylark.Struct
func (r *Response) Struct() *skylarkstruct.Struct {
	return skylarkstruct.FromStringDict(skylarkstruct.Default, skylark.StringDict{
		"url":         skylark.String(r.Request.URL.String()),
		"status_code": skylark.MakeInt(r.StatusCode),
		"headers":     r.HeadersDict(),
		"encoding":    skylark.String(strings.Join(r.TransferEncoding, ",")),

		"text":    skylark.NewBuiltin("text", r.Text),
		"content": skylark.NewBuiltin("content", r.Text),

		"json": skylark.NewBuiltin("json", r.JSON),
	})
}

// HeadersDict flops
func (r *Response) HeadersDict() *skylark.Dict {
	d := new(skylark.Dict)
	for key, vals := range r.Header {
		d.Set(skylark.String(key), skylark.String(strings.Join(vals, ",")))
	}
	return d
}

// Text returns the raw data as a string
func (r *Response) Text(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	// reset reader to allow multiple calls
	r.Body = ioutil.NopCloser(bytes.NewReader(data))

	return skylark.String(string(data)), nil
}

// JSON attempts to parse the response body as JSON
func (r *Response) JSON(thread *skylark.Thread, _ *skylark.Builtin, args skylark.Tuple, kwargs []skylark.Tuple) (skylark.Value, error) {
	var data interface{}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	r.Body.Close()
	// reset reader to allow multiple calls
	r.Body = ioutil.NopCloser(bytes.NewReader(body))
	return util.Marshal(data)
}
