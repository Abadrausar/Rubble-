package nca6

import "strconv"

// ParamsArray is a VERY trimed down base.ScriptArray.
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

func (this *ParamsArray) Set(index string, value *Value) {
	panic("You may not write to a ParamsArray.")
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

func (this *ParamsArray) ReadOnly() bool {
	return true
}
