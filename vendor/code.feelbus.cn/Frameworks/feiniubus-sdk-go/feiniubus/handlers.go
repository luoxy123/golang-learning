package feiniubus

import (
	"fmt"
	"strings"
)

// Handlers provides a collection of request handlers
type Handlers struct {
	Validate         HandlerList
	Build            HandlerList
	Send             HandlerList
	ValidateResponse HandlerList
	Unmarshal        HandlerList
	UnmarshalError   HandlerList
}

// Copy returns of this handler's lists.
func (h *Handlers) Copy() Handlers {
	return Handlers{
		Validate:         h.Validate.copy(),
		Build:            h.Build.copy(),
		Send:             h.Send.copy(),
		ValidateResponse: h.ValidateResponse.copy(),
		Unmarshal:        h.Unmarshal.copy(),
		UnmarshalError:   h.UnmarshalError.copy(),
	}
}

// HandlerListRunItem represents an entry in the HandlerList which
// is being run
type HandlerListRunItem struct {
	Index   int
	Handler NamedHandler
	Request *Request
}

// HandlerList manages zero or more handlers in a list
type HandlerList struct {
	list        []NamedHandler
	AfterEachFn func(item HandlerListRunItem) bool
}

// NamedHandler is a struct that contains a name and function callback
type NamedHandler struct {
	Name string
	Fn   func(*Request)
}

// PushBack pushes handler f to the back of the handler list.
func (l *HandlerList) PushBack(f func(*Request)) {
	l.list = append(l.list, NamedHandler{"__anonymous", f})
}

// PushFront pushes handler f  to the front of the handler list.
func (l *HandlerList) PushFront(f func(*Request)) {
	l.list = append(l.list, NamedHandler{"__anonymous", f})
}

// PushBackNamed pushes named handler f to the back of the handler list
func (l *HandlerList) PushBackNamed(n NamedHandler) {
	l.list = append(l.list, n)
}

// PushFrontNamed pushes named handler f to the front of the handler list.
func (l *HandlerList) PushFrontNamed(n NamedHandler) {
	l.list = append([]NamedHandler{n}, l.list...)
}

// Remove removes a NamedHandler n
func (l *HandlerList) Remove(n NamedHandler) {
	newlist := []NamedHandler{}
	for _, m := range l.list {
		if m.Name != n.Name {
			newlist = append(newlist, m)
		}
	}
	l.list = newlist
}

// Run executes all handlers in the list
func (l *HandlerList) Run(r *Request) {
	for i, h := range l.list {
		h.Fn(r)
		item := HandlerListRunItem{
			Index: i, Handler: h, Request: r,
		}
		if l.AfterEachFn != nil && !l.AfterEachFn(item) {
			return
		}
	}
}

func (l *HandlerList) copy() HandlerList {
	n := HandlerList{
		AfterEachFn: l.AfterEachFn,
	}
	n.list = append([]NamedHandler{}, l.list...)
	return n
}

// HandlerListStopOnError returns false to stop the HandlerList iterating
func HandlerListStopOnError(item HandlerListRunItem) bool {
	return item.Request.Error == nil
}

// MakeAddToUserAgentHandler will add the name/version pair to the User-Agent
func MakeAddToUserAgentHandler(name, version string, extra ...string) func(*Request) {
	ua := fmt.Sprintf("%s/%s", name, version)
	if len(extra) > 0 {
		ua += fmt.Sprintf(" (%s)", strings.Join(extra, "; "))
	}
	return func(r *Request) {
		AddToUserAgent(r, ua)
	}
}
