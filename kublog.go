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

	// create watch for pods
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
    	_, controller := cache.NewInformer(
        	watchlist,
        	&v1.Pod{},
        	time.Second * 0,
        	cache.ResourceEventHandlerFuncs{
            		AddFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
                		fmt.Printf("add - %s - %s - %s \n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, pod.ObjectMeta.CreationTimestamp)
            		},
            		DeleteFunc: func(obj interface{}) {
				pod := obj.(*v1.Pod)
                		fmt.Printf("delete: %s - %s - %s \n", pod.ObjectMeta.Namespace, pod.ObjectMeta.Name, pod.ObjectMeta.CreationTimestamp)
            		},
            		UpdateFunc:func(oldObj, newObj interface{}) {
				oldPod := oldObj.(*v1.Pod)
                		fmt.Printf("changed: %s - %s - %s \n", oldPod.ObjectMeta.Namespace, oldPod.ObjectMeta.Name, oldPod.ObjectMeta.CreationTimestamp)
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
