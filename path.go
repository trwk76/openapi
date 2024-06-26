package openapi

import (
	"fmt"

	"github.com/trwk76/openapi/spec"
)

func (a *API) NamedPath(name string, setup func(p *Path)) *Path {
	if a.param != nil {
		panic(fmt.Errorf("cannot add a named path since a parameter path is already declared"))
	} else if _, ok := a.named[name]; ok {
		panic(fmt.Errorf("named path '%s' is already declared", name))
	}

	res := newPath(a, "/"+name, "/"+name)
	a.named[name] = res

	if setup != nil {
		setup(res)
	}

	return res
}

type (
	Path struct {
		paths
		api    *API
		gpath  string
		apath  string
		params []spec.ParameterOrRef

		OpPrefix    string
		Summary     string
		Description string
		Tags        []Tag
	}

	paths struct {
		named map[string]*Path
		param *Path
	}
)

func (p *Path) NamedPath(name string, setup func(p *Path)) *Path {
	if p.param != nil {
		panic(fmt.Errorf("cannot add a named path since a parameter path is already declared"))
	} else if _, ok := p.named[name]; ok {
		panic(fmt.Errorf("named path '%s' is already declared", name))
	}

	res := newPath(p.api, p.gpath+"/"+name, p.apath+"/"+name)
	res.params = p.params
	res.OpPrefix = p.OpPrefix
	res.Tags = p.Tags
	p.named[name] = res

	if setup != nil {
		setup(res)
	}

	return res
}

func (p *Path) ParamPath(param spec.ParameterOrRef, setup func(p *Path)) *Path {
	if p.param != nil {
		panic(fmt.Errorf("a parameter path is already declared"))
	} else if len(p.named) > 0 {
		panic(fmt.Errorf("cannot declare a named parameter since named path(s) already declared"))
	}

	pdef := p.api.Param(param)
	if pdef.In != spec.ParameterPath {
		panic(fmt.Errorf("parameter is not a path parameter"))
	}

	res := newPath(p.api, p.gpath+"/:"+pdef.Name, p.apath+fmt.Sprintf("/{%s}", pdef.Name))
	res.params = append(p.params, param)
	res.OpPrefix = p.OpPrefix
	res.Tags = p.Tags
	p.param = res

	if setup != nil {
		setup(res)
	}

	return res
}

func newPath(api *API, gpath string, apath string) *Path {
	return &Path{
		paths: newPaths(),
		api:   api,
		gpath: gpath,
		apath: apath,
	}
}

func newPaths() paths {
	return paths{
		named: make(map[string]*Path),
	}
}
