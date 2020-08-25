package unit

import (
	"fmt"
	"strings"
)

type Registry struct {
	syms   map[string]Maker
	pref   map[string]Maker
	Parent *Registry
}

func (r *Registry) Primitive(symbol string, alias ...string) Maker {
	m := Primitive(symbol)
	return r.register(&r.syms, symbol, m, alias...)
}

func (r *Registry) Derive(symbol string, v Qualified, alias ...string) Maker {
	m := Derive(symbol, v)
	return r.register(&r.syms, symbol, m, alias...)
}

func (r *Registry) Register(m Maker, alias ...string) Maker {
	if u := m.Unit(); u != nil {
		return r.register(&r.syms, u.Symbol(), m, alias...)
	}
	panic("Register can only be used for named units")
}

func (r *Registry) Prefix(symbol string, m float64, alias ...string) func(m Maker) Maker {
	mk := r.register(&r.pref, symbol, Scalar(m), alias...)
	return mk.Mul
}

func (r *Registry) register(where *map[string]Maker, symbol string, m Maker, alias ...string) Maker {
	if *where == nil {
		*where = make(map[string]Maker)
	}
	for _, s := range append(alias, symbol) {
		if (*where)[s] != nil {
			panic(fmt.Sprintf("Unit %q already registered", s))
		}
		(*where)[s] = m
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
				return sub.Mul(m)
			}
		}
	}
	return r.Parent.Find(n)
}
