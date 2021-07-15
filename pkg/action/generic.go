package action

type HTTPMethod string

const (
	HTTPMethodGet     HTTPMethod = "GET"
	HTTPMethodPost    HTTPMethod = "POST"
	HTTPMethodPut     HTTPMethod = "PUT"
	HTTPMethodPatch   HTTPMethod = "PATCH"
	HTTPMethodOptions HTTPMethod = "OPTIONS"
	HTTPMethodDelete  HTTPMethod = "DELETE"
)

func (m HTTPMethod) IsValid() error {
	switch m {
	case HTTPMethodGet, HTTPMethodPost, HTTPMethodPut, HTTPMethodPatch, HTTPMethodOptions, HTTPMethodDelete:
		return nil
	}
	return ErrHTTPMethodNotValid
}

func (m HTTPMethod) String() string {
	return string(m)
}
