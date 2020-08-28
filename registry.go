package unit

import (
	"fmt"
	"strings"
)

// Registry is an experimental type for recording unit symbols for lookup later.
//
// This type is under development and will likely change.
type Registry struct {
	name   string
	syms   map[string]Maker
	pref   map[string]func(m Maker) Maker
	Parent *Registry
}

func NewRegistry(name string, parent *Registry) *Registry {
	return &Registry{Parent: parent}
}

func (r *Registry) Primitive(symbol string, alias ...string) Maker {
	m := Primitive(symbol)
	return r.register(symbol, m, alias...)
}

func (r *Registry) Derive(symbol string, v Qualified, alias ...string) Maker {
	m := Derive(symbol, v)
	return r.register(symbol, m, alias...)
}

func (r *Registry) Register(m Maker, alias ...string) Maker {
	if u := m.Unit(); u != nil {
		return r.register(u.Symbol(), m, alias...)
	}
	panic("Register can only be used for named units")
}

func (r *Registry) Prefix(symbol string, mult float64, alias ...string) func(m Maker) Maker {
	m := Prefix(symbol, mult)
	if r.pref == nil {
		r.pref = make(map[string]func(Maker) Maker)
	}
	for _, s := range append(alias, symbol) {
		if r.pref[s] != nil {
			panic(fmt.Sprintf("Prefix %q already registered", s))
		}
		r.pref[s] = m
	}
	return m
}

func (r *Registry) register(symbol string, m Maker, alias ...string) Maker {
	if r.syms == nil {
		r.syms = make(map[string]Maker)
	}
	for _, s := range append(alias, symbol) {
		if r.syms[s] != nil {
			panic(fmt.Sprintf("Unit %q already registered", s))
		}
		r.syms[s] = m
	}
	return m
}

func (r *Registry) Find(n string) Maker {
	if r == nil {
		return nil
	}
	if got := r.syms[n]; got != nil {
		return got
	}
	for prefix, m := range r.pref {
		if strings.HasPrefix(n, prefix) {
			if sub := r.Find(n[len(prefix):]); sub != nil {
				return m(sub)
			}
		}
	}
	return r.Parent.Find(n)
}
