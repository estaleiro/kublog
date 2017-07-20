package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

        "k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Started...")
	fmt.Println("Operation - Namespace - Name - Date - Pod IP - Host IP")
	// create watch for pods
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
    	_, controller := cache.NewInformer(
        	watchlist,
        	&v1.Pod{},
        	time.Second * 0,
        	cache.ResourceEventHandlerFuncs{
            		AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				namespace := pod.ObjectMeta.Namespace
				name := pod.ObjectMeta.Name
				creationTime := pod.ObjectMeta.CreationTimestamp

				fmt.Printf("add - %s - %s - %s - %s - %s \n", namespace, name, creationTime, pod.Status.PodIP, pod.Status.HostIP)
            		},
            		DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
				namespace := pod.ObjectMeta.Namespace
				name := pod.ObjectMeta.Name
				creationTime := pod.ObjectMeta.CreationTimestamp

				fmt.Printf("delete - %s - %s - %s - %s - %s \n", namespace, name, creationTime, pod.Status.PodIP, pod.Status.HostIP)
            		},
            		UpdateFunc:func(oldObj, newObj interface{}) {
				oldPod := oldObj.(*v1.Pod)
				namespace := oldPod.ObjectMeta.Namespace
				name := oldPod.ObjectMeta.Name
				creationTime := oldPod.ObjectMeta.CreationTimestamp

				fmt.Printf("change - %s - %s - %s - %s - %s \n", namespace, name, creationTime, oldPod.Status.PodIP, oldPod.Status.HostIP)
            		},
        	},
    	)
    	stop := make(chan struct{})
    	go controller.Run(stop)
	for{
        	time.Sleep(time.Second)
    	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
