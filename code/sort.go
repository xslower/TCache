package main

type Uint32Slice []uint32

func (u Uint32Slice) Len() int {
	return len(u)
}

func (u Uint32Slice) Less(i, j int) bool {
	return u[i] < u[j]
}

func (u Uint32Slice) Swap(i, j int) {
	u[i], u[j] = u[j], u[i]
}

type Int64Slice []int64

//Len
func (s Int64Slice) Len() int {
	return len(s)
}

//Less():成绩将有低到高排序
func (s Int64Slice) Less(i, j int) bool {
	return s[i] < s[j]
}

//Swap()
func (s Int64Slice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type MarkSlice []Unit

//Len
func (s MarkSlice) Len() int {
	return len(s)
}

//Less():成绩将有低到高排序
func (s MarkSlice) Less(i, j int) bool {
	return s[i].Score < s[j].Score
}

//Swap()
func (s MarkSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
