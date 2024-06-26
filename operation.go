package openapi

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/trwk76/openapi/spec"
)

func (p *Path) GET(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	if res.RequestBody != nil {
		panic(fmt.Errorf("http 'GET' operation does not accept a request body"))
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.GET = res
		p.api.g.GET(p.gpath, hdl)
	})

	return res
}

func (p *Path) PUT(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.PUT = res
		p.api.g.PUT(p.gpath, hdl)
	})

	return res
}

func (p *Path) POST(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.POST = res
		p.api.g.POST(p.gpath, hdl)
	})

	return res
}

func (p *Path) DELETE(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	if res.RequestBody != nil {
		panic(fmt.Errorf("http 'DELETE' operation does not accept a request body"))
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.DELETE = res
		p.api.g.DELETE(p.gpath, hdl)
	})

	return res
}

func (p *Path) OPTIONS(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	if res.RequestBody != nil {
		panic(fmt.Errorf("http 'OPTIONS' operation does not accept a request body"))
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.OPTIONS = res
		p.api.g.OPTIONS(p.gpath, hdl)
	})

	return res
}

func (p *Path) HEAD(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	if res.RequestBody != nil {
		panic(fmt.Errorf("http 'HEAD' operation does not accept a request body"))
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.HEAD = res
		p.api.g.HEAD(p.gpath, hdl)
	})

	return res
}

func (p *Path) PATCH(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.PATCH = res
		p.api.g.PATCH(p.gpath, hdl)
	})

	return res
}

func (p *Path) TRACE(id string, hdl gin.HandlerFunc, setup func(s *spec.Operation)) *spec.Operation {
	tags := make([]string, len(p.Tags))

	for idx, itm := range p.Tags {
		tags[idx] = itm.name
	}

	res := &spec.Operation{OperationID: p.OpPrefix + id, Tags: tags}

	if setup != nil {
		setup(res)
	}

	if res.RequestBody != nil {
		panic(fmt.Errorf("http 'TRACE' operation does not accept a request body"))
	}

	p.ensureItem(func(s *spec.PathItem) {
		s.TRACE = res
		p.api.g.Handle(http.MethodTrace, p.gpath, hdl)
	})

	return res
}

func (p *Path) ensureItem(f func(s *spec.PathItem)) {
	pi, fnd := p.api.s.Paths[p.apath]
	if !fnd {
		pi.Summary = p.Summary
		pi.Description = p.Description
		pi.Parameters = p.params
	}

	f(&pi)

	p.api.s.Paths[p.apath] = pi
}
