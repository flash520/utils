/**
 * @Author: koulei
 * @Description: TODO
 * @File:  namingClient
 * @Version: 1.0.0
 * @Date: 2021/5/11 22:46
 */

package discovery

import (
	"errors"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	log "github.com/sirupsen/logrus"
)

func RegisterServiceInstance(param vo.RegisterInstanceParam) {
	success, _ := nacosClient.RegisterInstance(param)
	log.Infof("RegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func DeRegisterServiceInstance(param vo.DeregisterInstanceParam) {
	success, _ := nacosClient.DeregisterInstance(param)
	log.Infof("DeRegisterServiceInstance,param:%+v,result:%+v \n\n", param, success)
}

func (n Nacos) GetInstance(serviceName string) (string, error) {
	params := vo.SelectOneHealthInstanceParam{
		Clusters:    nil,
		ServiceName: serviceName,
	}
	instance := n.SelectOneHealthyInstance(params)
	if instance == nil {
		return "", errors.New("服务不存在")
	}
	url := "http://" + instance.Ip + ":" + strconv.FormatUint(instance.Port, 10)
	//fmt.Printf("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
	return url, nil
}

func (n *Nacos) GetService(param vo.GetServiceParam) {
	service, _ := nacosClient.GetService(param)
	log.Infof("GetService,param:%+v, result:%+v \n\n", param, service)
}

func (n *Nacos) SelectAllInstances(param vo.SelectAllInstancesParam) {
	instances, _ := nacosClient.SelectAllInstances(param)
	log.Infof("SelectAllInstance,param:%+v, result:%+v \n\n", param, instances)
}

func (n *Nacos) SelectInstances(param vo.SelectInstancesParam) {
	instances, _ := nacosClient.SelectInstances(param)
	log.Infof("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
}

func (n *Nacos) SelectOneHealthyInstance(param vo.SelectOneHealthInstanceParam) *model.Instance {
	instances, _ := nacosClient.SelectOneHealthyInstance(param)
	log.Infof("SelectInstances,param:%+v, result:%+v \n\n", param, instances)
	return instances
}

func (n *Nacos) Subscribe(param *vo.SubscribeParam) {
	_ = nacosClient.Subscribe(param)
}

func (n *Nacos) UnSubscribe(param *vo.SubscribeParam) {
	_ = nacosClient.Unsubscribe(param)
}

func (n *Nacos) GetAllService(param vo.GetAllServiceInfoParam) model.ServiceList {
	service, _ := nacosClient.GetAllServicesInfo(param)
	log.Infof("GetAllService,param:%+v, result:%+v \n\n", param, service)
	return service
}
