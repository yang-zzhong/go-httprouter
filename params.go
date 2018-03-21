package httprouter

type Params struct {
	params map[string]interface{}
}

type eachCall func(string, interface{}) bool

func NewP() *Params {
	p := new(Params)
	p.params = make(map[string]interface{})

	return p
}

func (p *Params) Set(key string, value interface{}) {
	p.params[key] = value
}

func (p *Params) Get(key string) interface{} {
	return p.params[key]
}

func (p *Params) Del(key string) {
	delete(p.params, key)
}

func (p *Params) Each(each eachCall) {
	for key, name := range p.params {
		if !each(key, name) {
			return
		}
	}
}

func (p *Params) Exists(key string) bool {
	_, ok := p.params[key]

	return ok
}
