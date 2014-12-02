package base

import "dctech/nca6"

type ScriptMap map[string]*nca6.Value

func NewScriptMap() nca6.Indexable {
	return make(ScriptMap, 20)
}

func (this ScriptMap) Get(index string) *nca6.Value {
	if _, ok := this[index]; ok {
		return this[index]
	}
	return nca6.NewValueInt64(0)
}

func (this ScriptMap) Set(index string, value *nca6.Value) {
	this[index] = value
}

func (this ScriptMap) Exists(index string) bool {
	_, ok := this[index]
	return ok
}

func (this ScriptMap) Len() int64 {
	return int64(len(this))
}

func (this ScriptMap) Keys() []string {
	rtn := make([]string, 0, len(this))
	for key := range this {
		rtn = append(rtn, key)
	}
	return rtn
}

func (this ScriptMap) ReadOnly() bool {
	return false
}
