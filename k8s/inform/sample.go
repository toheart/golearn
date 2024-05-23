package main

/**
@file:
@author: levi.Tang
@time: 2024/5/16 11:22
@description:
**/

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"time"
)

func main() {
	config, err := clientcmd.BuildConfigFromFlags("", "/home/ubuntu/.kube/config")
	if err != nil {
		panic(err)
	}
	config.Timeout = time.Second * 5
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	stopCh := make(chan struct{})
	defer close(stopCh)
	log.Printf("create new inforemer \n")
	sharedInformers := informers.NewSharedInformerFactory(clientset, time.Minute)

	informer := sharedInformers.Core().V1().Pods().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("New Pod Added to Store: %s", mObj.GetName())
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			oObj := oldObj.(v1.Object)
			nObj := newObj.(v1.Object)

			log.Printf("%s Pod Update to %s", oObj.GetName(), nObj.GetName())
		},
		DeleteFunc: func(obj interface{}) {
			mObj := obj.(v1.Object)
			log.Printf("Pod Deleted from Store:%s", mObj.GetName())
		},
	})
	informer.Run(stopCh)
}
