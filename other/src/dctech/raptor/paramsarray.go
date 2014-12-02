/*
For copyright/license see header in file "doc.go"
*/

package raptor

import "strconv"

// ParamsArray is a trimmed down base.ScriptArray.
type ParamsArray []*Value

func NewParamsArray(val []*Value) Indexable {
	rtn := ParamsArray(val)
	return &rtn
}

func (this *ParamsArray) Get(index string) *Value {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		panic("Index not a valid number.")
	}
	if val < 0 || val >= int64(len(*this)) {
		panic("Index out of range.")
	}
	return (*this)[val]
}

func (this *ParamsArray) Exists(index string) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(*this)) {
		return false
	}
	return true
}

func (this *ParamsArray) Len() int64 {
	return int64(len(*this))
}

func (this *ParamsArray) Keys() []string {
	rtn := make([]string, 0, len(*this))
	for key := range *this {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}

func (this *ParamsArray) String() string {
	return GenericIndexableToString("ParamsArray", this)
}

func (this *ParamsArray) CodeString() string {
	return GenericIndexableToCodeString("array", this, false)
}
