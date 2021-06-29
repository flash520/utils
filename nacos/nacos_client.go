/**
 * @Author: koulei
 * @Description: TODO
 * @File:  namingClient
 * @Version: 1.0.0
 * @Date: 2021/5/11 22:46
 */

package nacos

import (
	"fmt"

	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func RegisterServiceInstance(param vo.RegisterInstanceParam) {
	success, _ := namingClient.RegisterInstance(param)
	fmt.Printf("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func DeRegisterServiceInstance(param vo.DeregisterInstanceParam) {
	success, _ := namingClient.DeregisterInstance(param)
	fmt.Printf("DeRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func GetService(param vo.GetServiceParam) {
	service, _ := namingClient.GetService(param)
	fmt.Printf("GetService,param:%+v, result:%+v \n\n", param, service)
}

func SelectAllInstances(param vo.SelectAllInstancesParam) {
	instances, _ := namingClient.SelectAllInstances(param)
	fmt.Printf("SelectAllInstance,param:%+v, result:%+v \n\n", param, instances)
}

func SelectInstances(param vo.SelectInstancesParam) {
	instances, _ := namingClient.SelectInstances(param)
	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
}

func SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) *model.Instance {
	instances, _ := namingClient.SelectOneHealthyInstance(param)
	fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
	return instances
}

func Subscribe(param *vo.SubscribeParam) {
	_ = namingClient.Subscribe(param)
}

func UnSubscribe(param *vo.SubscribeParam) {
	_ = namingClient.Unsubscribe(param)
}

func GetAllService(param vo.GetAllServiceInfoParam) model.ServiceList {
	service, _ := namingClient.GetAllServicesInfo(param)
	fmt.Printf("GetAllService,param:%+v, result:%+v \n\n", param, service)
	return service
}
