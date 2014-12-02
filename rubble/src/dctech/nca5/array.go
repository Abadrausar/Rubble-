package nca5

//import "fmt"
import "strconv"

type ScriptArray []*Value

func NewScriptArray() Indexable {
	return make(ScriptArray, 0, 20)
}

func NewScriptArraySized(size int) Indexable {
	return make(ScriptArray, size)
}

func (this ScriptArray) Get(index string) *Value {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		panic("Index not a valid number.")
	}
	if val < 0 || val >= int64(len(this)) {
		panic("Index out of range.")
	}
	return this[val]
}

func (this ScriptArray) Set(index string, value *Value) {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		panic("Index not a valid number.")
	}
	if val < 0 {
		panic("Index too small.")
	}
	if val < int64(len(this)) {
		this[val] = value
		return
	}
	if val == int64(len(this)) {
		this = append(this, value)
	}
	panic("Index too large, you may only extend the array by one at a time.")
}

func (this ScriptArray) Exists(index string) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(this)) {
		return false
	}
	return true
}

func (this ScriptArray) Len() int64 {
	return int64(len(this))
}

func (this ScriptArray) Keys() []string {
	rtn := make([]string, 0, len(this))
	for key := range this {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}
