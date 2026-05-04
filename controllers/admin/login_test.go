package admin

import (
	"fmt"
	"testing"

	"github.com/zsy619/tools/xcrypto"

	"haedu.gov.cn/cms/global"
)

func Test_CreatePassword(t *testing.T) {
	key := global.ReverseLowerString("admin@2023")
	password, err := xcrypto.Sm4Encrypt("Admin$2023xyz", key)
	if err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(password)
	}
}

func Test_Password(t *testing.T) {
	key := global.ReverseLowerString("admin@2023")
	password, err := xcrypto.Sm4Decrypt("c9431549e347910a35fcf4e4a04bcb75", key)
	if err != nil {
		t.Fatal(err)
	} else {
		// Admin@2023
		fmt.Println(password)
	}
}
