package csv

import "dctech/nca6"
import "strconv"

type CSVFile []*nca6.Value

func NewCSVFile() nca6.Indexable {
	tmp := make(CSVFile, 0, 20)
	return &tmp
}

func (this *CSVFile) Get(index string) *nca6.Value {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		panic("Index not a valid number.")
	}
	if val < 0 || val >= int64(len(*this)) {
		panic("Index out of range.")
	}
	return (*this)[val]
}

func (this *CSVFile) Set(index string, value *nca6.Value) {
	panic("CSV files are read only.")
}

func (this *CSVFile) Exists(index string) bool {
	val, err := strconv.ParseInt(index, 0, 64)
	if err != nil {
		return false
	}
	if val < 0 || val >= int64(len(*this)) {
		return false
	}
	return true
}

func (this *CSVFile) Len() int64 {
	return int64(len(*this))
}

func (this *CSVFile) Keys() []string {
	rtn := make([]string, 0, len(*this))
	for key := range *this {
		rtn = append(rtn, strconv.FormatInt(int64(key), 10))
	}
	return rtn
}
