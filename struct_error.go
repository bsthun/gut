package gut

import "errors"

type ErrorBlock struct {
	Message string `json:"message,omitempty"`
	Code    string `json:"code,omitempty"`
	Err     error  `json:"error,omitempty"`
}

type ErrorInstance struct {
	Errors []*ErrorBlock `json:"errors"`
}

func (v *ErrorInstance) Error() string {
	if len(v.Errors) > 0 {
		return v.Errors[len(v.Errors)-1].Message
	}
	return "Empty error"
}

func Err(critical bool, message string, args2 ...any) *ErrorInstance {
	if len(args2) == 1 {
		if code, ok := args2[0].(string); ok {
			return &ErrorInstance{
				Errors: []*ErrorBlock{
					{
						Message: message,
						Code:    code,
						Err:     nil,
					},
				},
			}
		}
		if err, ok := args2[0].(error); ok {
			var code *ErrorInstance
			if errors.As(args2[0].(error), &code) {
				return &ErrorInstance{
					Errors: append(code.Errors, &ErrorBlock{
						Message: message,
						Code:    "",
						Err:     err,
					}),
				}
			}
			return &ErrorInstance{
				Errors: []*ErrorBlock{
					{
						Message: message,
						Code:    "",
						Err:     err,
					},
				},
			}
		}
	}

	if len(args2) == 2 {
		if code, ok := args2[0].(string); ok {
			if err, ok := args2[1].(error); ok {
				return &ErrorInstance{
					Errors: []*ErrorBlock{
						{
							Message: message,
							Code:    code,
							Err:     err,
						},
					},
				}
			}
		}
	}

	return &ErrorInstance{
		Errors: []*ErrorBlock{
			{
				Message: message,
				Err:     nil,
			},
		},
	}
}
