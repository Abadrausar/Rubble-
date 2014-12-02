package nca5

//import "fmt"

type ScriptMap map[string]*Value

func NewScriptMap() Indexable {
	return make(ScriptMap, 20)
}

func (this ScriptMap) Get(index string) *Value {
	if _, ok := this[index]; ok {
		return this[index]
	}
	return NewValue("")
}

func (this ScriptMap) Set(index string, value *Value) {
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
