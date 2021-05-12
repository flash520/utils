/**
 * @Author: koulei
 * @Description: TODO
 * @File:  mqtt_test
 * @Version: 1.0.0
 * @Date: 2021/5/12 00:57
 */

package mqtt

import "testing"

func TestMqtt(t *testing.T) {
	mc := CreateMqttClient("mqtt.dev.cdqidi.cn:1883", "test", "test").GetConn()
	if token := mc.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	token := mc.Publish("test/abc", 0, false, "test")
	token.Wait()
}
