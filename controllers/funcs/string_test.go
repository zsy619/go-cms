package funcs

import "testing"

func TestSubStr(t *testing.T) {
	tstr := "河南省教育网有限公司"
	t.Log(SubStr(tstr, 0, 10))
	t.Log(SubStr(tstr, -1, 0))
	t.Log(SubStr(tstr, 0, 22))
}
