package main

import (
	"fmt"
	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listerAppsV1 "k8s.io/client-go/listers/apps/v1"
	listerCoreV1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"strconv"
	"time"
)

var err error
var config *rest.Config
var restartNumber string

func getKubeConfig() *kubernetes.Clientset {
	kubeconfig := fmt.Sprintf("%s%s", os.Getenv("HOME"), "/.kube/config")
	if config, err = rest.InClusterConfig(); err != nil {
		fmt.Println(err)
		if config, err = clientcmd.BuildConfigFromFlags("", kubeconfig); err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func createInformer(clientset *kubernetes.Clientset, listTime time.Duration) (informers.SharedInformerFactory, listerAppsV1.DeploymentLister, listerCoreV1.PodLister) {
	informerFactory := informers.NewSharedInformerFactory(clientset, listTime)
	podInformer := informerFactory.Core().V1().Pods()
	deployInformer := informerFactory.Apps().V1().Deployments()

	pinformer := podInformer.Informer()
	podLister := podInformer.Lister()
	deployLister := deployInformer.Lister()
	pinformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    onAdd,
		UpdateFunc: onUpdate,
		DeleteFunc: onDelete,
	})

	log.Println("informerFactory deployLister podLister initialization completed")
	return informerFactory, deployLister, podLister
}
func offUnHealthyDeployment(deploy []*appsV1.Deployment,
	namespace string,
	clientset *kubernetes.Clientset,
	deployName string) error {

	for _, dp := range deploy {
		//找到对应的deployName
		if dp.Name == deployName {
			//判断deployment是否是属于开启状态 不等于0就是开启状态
			if dp.Status.UpdatedReplicas != 0 && dp.Status.AvailableReplicas == 0 {
				//判断当前运行的pod  等于0就是deployment是开启状态且running的pod为0 确认为0了以后才可以执行
				var replicas int32 = 0
				dp.Spec.Replicas = &replicas
				_, err := clientset.AppsV1().Deployments(dp.Namespace).Update(dp)
				if err != nil {
					log.Fatal("deployment replica set to 0 failed", err)
				}
			}
		}
	}
	return nil
}

func resouceRegistry(podLister listerCoreV1.PodLister, deployLister listerAppsV1.DeploymentLister, namespace string) ([]*coreV1.Pod, []*appsV1.Deployment) {
	pod, err := podLister.Pods(namespace).List(labels.Everything())
	if err != nil {
		log.Fatalf("PodLister Get Pod List Fatal %s\n", err)
	}

	deploy, err := deployLister.Deployments(namespace).List(labels.Everything())
	if err != nil {
		log.Fatalf("DeployLister Get Deployment List Fatal %s\n", err)
	}
	return pod, deploy
}

func offUnHealthyApp(pod []*coreV1.Pod, namespace string, restartNumber int32, deploy []*appsV1.Deployment, clientset *kubernetes.Clientset) {
	var rslice []string
	for _, p := range pod {
		//判断是否为ReplicaSet资源
		//这不是一个呆b行为 二次range是有原因的
		for _, owner := range p.OwnerReferences {
			//完成kind判断以后 就可以避开静态pod不存在这个字段而报错 index out of range [0] with length 0
			if owner.Kind == "ReplicaSet" {
				if len(p.Status.ContainerStatuses) != 0 {
					if p.Status.ContainerStatuses[0].Ready == false && p.Status.ContainerStatuses[0].RestartCount >= restartNumber {
						//将获取的relicaset名字切割末尾的11位字符串(包含-)得到deployment名字 存入切片为了pod去重
						rslice = append(rslice, p.OwnerReferences[0].Name[0:len(p.OwnerReferences[0].Name)-11])
					}
				}
			}
		}
	}
	if len(rslice) != 0 {
		offUnHealthyDeployment(deploy, namespace, clientset, rslice[0])
	}
	return
}

func init() {
	restartNumber = fmt.Sprintf(os.Getenv("POD_RESTART_NUMBER"))

}

func main() {
	log.Println("Project address: https://github.com/huangjc7/sunshine")

	r, _ := strconv.ParseInt(restartNumber, 10, 32) //string转int32

	clientset := getKubeConfig()
	stopCh := make(chan struct{})
	informerFactory, deployLister, podLister := createInformer(clientset, time.Minute*5)
	informerFactory.Start(stopCh)
	log.Println("informerFactory started")
	informerFactory.WaitForCacheSync(stopCh)
	log.Println("Start successfully")
	for {
		time.Sleep(time.Second * 2)
		pod, deploy := resouceRegistry(podLister, deployLister, "")
		offUnHealthyApp(pod, "", int32(r), deploy, clientset)
	}

}

func onAdd(obj interface{}) {}

func onUpdate(old, new interface{}) {}

func onDelete(obj interface{}) {
	pod := obj.(*coreV1.Pod) //断言 是否是deployment类型
	log.Printf("Delete Pod Namespace:%s Pod Name: %s\n", pod.Namespace, pod.Name)
}
