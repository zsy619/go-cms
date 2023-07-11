package admin

import (
	"fmt"
	"testing"

	"haedu.gov.cn/cms/global"
	"haedu.gov.cn/tools/xcrypto"
)

func Test_CreatePasswo(t *testing.T) {
	key := global.ReverseLowerString("admin@2023")
	password, err := xcrypto.Sm4Encrypt("Admin$2023xyz", key)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(password)
	}
}
