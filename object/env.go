package object

type Env struct {
	store map[string]Object
	top   *Env
}

func NewEnv(envs *Env) *Env {
	return &Env{store: make(map[string]Object), top: envs}
}
func (r *Env) Get(name string) Object {
	v, ok := r.store[name]
	if !ok && r.top != nil {
		return r.top.Get(name)
	}
	//if v, ok_ := table[name]; ok_ {
	//	return &BuiltFun{Fun_: v}
	//}
	if !ok {
		return NULL_
	}
	return v
}
func (r *Env) Set(name string, object Object) {
	r.store[name] = object
}
