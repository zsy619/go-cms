//golang 模拟登陆微信公众平台，突破微信群发每日一条限制

package simulate

import (
	"fmt"
	"log"
	"testing"
)

func TestNewWebWeChat(t *testing.T) {
	wechat := NewWebWeChat("49660513@qq.com", "yy04081018")

	if wechat.Login() == true {
		log.Println(wechat.GetFakeId())
		tofakeid := "oVpLw5hVbF0KuMmJTWE7BrtO4GXE" //my fakeid for test
		wechat.SendTextMsg(tofakeid, "Hello Phil.")
	} else {
		fmt.Println("wechat login failed.")
	}
}
