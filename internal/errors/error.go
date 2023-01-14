package errors

import "errors"

// TODO
/*type InternalError struct {
	Code int
	Err  error
}

func (i *InternalError) Error() string {
	return i.Err.Error()
}*/

var (
	ErrGithubUserNotFound          = errors.New("github user with provided username not found")
	ErrInvalidToken                = errors.New("invalid token")
	ErrAccessDenied                = errors.New("access denied")
	ErrUserNotFound                = errors.New("user not found")
	ErrProjectNotFound             = errors.New("project not found")
	ErrEmailAlreadyExists          = errors.New("email already exists")
	ErrUsernameAlreadyExists       = errors.New("username already exists")
	ErrGithubUsernameAlreadyExists = errors.New("github username already exists")
	ErrInvalidRole                 = errors.New("invalid role")
	ErrProjectNameAlreadyExists    = errors.New("project name already exists")
	ErrProjectDateIsNotValid       = errors.New("project date is not valid")
)

// package ierr

// import (
// 	"fmt"
// 	"net/http"
// 	"sort"
// 	"strings"

// type Cancel struct {
// 	msg   string
// 	code  int
// 	props map[string]interface{}

// 	next *Error
// }

// const (
// 	ID     = "id"
// 	Amount = "amount"
// 	Field  = "field"
// 	Fields = "fields"
// 	Value  = "value"
// )

// var (
// 	ErrKYCTypeUnknown                = New("unknown KYC type").BadRequest()
// 	ErrKYCIDNotProvided              = New("KYC ID is not provided").BadRequest()
// 	ErrInvalidAmount                 = New("invalid amount format").BadRequest()
// 	ErrKYCEmptyField                 = New("a KYC field mustn't be empty").BadRequest()
// 	ErrKYCEmptyFields                = New("KYC fields mustn't be empty").BadRequest()
// 	ErrKYCOnlyTypeProvided           = New("only 'type' field was provided").BadRequest()
// 	ErrCustomerNameTooLong           = New("a customer name should be at maximum 50 characters long").BadRequest()
// 	ErrCustomerAddressTooLong        = New("a customer address should be at maximum 100 characters long").BadRequest()
// 	ErrCustomerTypeNotSpecified      = New("query parameter 'type' is not specified").BadRequest()
// 	ErrKYCUnexpectedCountryCode      = New("unexpected KYC country code, it must be in ISO Alpha-3 format").BadRequest()
// 	ErrWrongCancelTransactionStatus  = New("transaction status must be 'pending_external' to cancel transaction").BadRequest()
// 	ErrWrongCancelTransactionType    = New("transaction type must be 'cash' to cancel transaction").BadRequest()
// 	ErrUnexpectedTransactionBankCode = New("unexpected transaction bank code").BadRequest()
// 	ErrUnexpectedTransactionType     = New("unexpected transaction type").BadRequest()
// 	ErrWrongTransactionBankAccount   = New("wrong transaction bank account number").BadRequest()

// 	ErrUnauthorized = New("use the JWT from /auth endpoint in 'Authorization: Bearer <JWT>' for authorization").Unauthorized()

// 	ErrInvalidToken = New("invalid token").Forbidden()

// 	ErrCustomerNotFound    = New("customer not found").NotFound()
// 	ErrTransactionNotFound = New("transaction not found").NotFound()
// )

// func Internal(err error) *Error {
// 	return &Error{
// 		msg:   err.Error(),
// 		code:  0,
// 		props: map[string]interface{}{},
// 		next:  ErrUnexpectedTransactionBankCode,
// 	}
// }

// func New(msg string) *Error {
// 	return &Error{msg: msg}
// }

// func Get(err error) *Error {
// 	if ierr, ok := err.(*Error); ok {
// 		return ierr
// 	}

// 	return Internal(err)
// }

// func Wrap(err error, msg string) *Error {
// 	return Get(err).Wrap(msg)
// }

// func (e *Error) copy() *Error {
// 	return &Error{
// 		msg:  e.msg,
// 		code: e.code,
// 	}
// }

// func (e *Error) error() string {
// 	msg := e.msg
// 	if e.props != nil {
// 		props := make([]string, 0, len(e.props))
// 		for prop, val := range e.props {
// 			switch v := val.(type) {
// 			case []string:
// 				sort.Strings(v)
// 			}
// 			props = append(props, fmt.Sprintf("%s: %s", prop, val))
// 		}
// 		sort.Strings(props)
// 		msg += fmt.Sprintf(" { %s }", strings.Join(props, ", "))
// 	}

// 	return msg
// }

// func (e *Error) Error() string {
// 	msgs := []string{e.error()}
// 	for ierr := e; ierr.next != nil; {
// 		ierr = ierr.next
// 		msgs = append(msgs, ierr.error())
// 	}
// 	sort.Strings(msgs)

// 	return strings.Join(msgs, "; ")
// }

// func (e *Error) Code() int {
// 	if e.code == 0 {
// 		return http.StatusInternalServerError
// 	}
// 	return e.code
// }

// func (e *Error) Is(err error) bool {
// 	return e.msg == err.Error()
// }

// func (e *Error) Wrap(msg string) *Error {
// 	err := e.copy()
// 	err.msg = msg + ": " + e.Error()
// 	return err
// }

// func (e *Error) Add(err error) *Error {
// 	if e == nil {
// 		return Get(err)
// 	}

// 	var ierr *Error
// 	for ierr = e; ierr.next != nil; ierr = ierr.next {
// 	}

// 	ierr.next = Get(err)
// 	return e
// }

// func (e Error) WithProperty(prop string, val interface{}) *Error {
// 	if e.props == nil {
// 		e.props = make(map[string]interface{})
// 	}
// 	e.props[prop] = val
// 	return &e
// }

// func (e Error) WithProperties(props ...interface{}) *Error {
// 	if e.props == nil {
// 		e.props = make(map[string]interface{}, len(props)/2)
// 	}

// 	if len(props)%2 == 1 {
// 		props = append(props, "")
// 	}

// 	for i := 0; i < len(props); i += 2 {
// 		e.props[props[i].(string)] = props[i+1]
// 	}

// 	return &e
// }

// func (e Error) BadRequest() *Error {
// 	return e.updateCode(http.StatusBadRequest)
// }
// func (e Error) Unauthorized() *Error {
// 	return e.updateCode(http.StatusUnauthorized)
// }
// func (e Error) Forbidden() *Error {
// 	return e.updateCode(http.StatusForbidden)
// }
// func (e Error) NotFound() *Error {
// 	return
// 	e.updateCode(http.StatusNotFound)
// }
// func (e Error) UnprocessableEntity() *Error {
// 	return e.updateCode(http.StatusUnprocessableEntity)
// }

// func (e *Error) isInternal() bool {
// 	return e.code == 0
// }

// func (e *Error) updateCode(code int) *Error {
// 	if e.isInternal() {
// 		e.code = code
// 	}
// 	return e
// }
