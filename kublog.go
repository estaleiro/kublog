package main

import (
	"log"
	"os"
	"time"

	"kublog/config"
	"kublog/outputs"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"github.com/namsral/flag"
)

func main() {
	var configFile string
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "KUBLOG", 0)
	fs.StringVar(&configFile, "config", "kublog.toml", "config file (mandatory)")
	fs.Parse(os.Args[1:])

	config, err := config.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	// set log
	f, err := os.OpenFile(config.LogFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	// use the current context in kubeconfig
	cfg, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Started...")

	// create watch for pods
	watchlist := cache.NewListWatchFromClient(clientset.Core().RESTClient(), "pods", v1.NamespaceDefault, fields.Everything())
	_, controller := cache.NewInformer(
		watchlist,
		&v1.Pod{},
		time.Second*0,
		cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				outputs.IndexKublog(obj.(*v1.Pod), "add", *config)
			},
			DeleteFunc: func(obj interface{}) {
				outputs.IndexKublog(obj.(*v1.Pod), "delete", *config)
			},
			UpdateFunc: func(oldObj, newObj interface{}) {
				outputs.IndexKublog(oldObj.(*v1.Pod), "update", *config)
			},
		},
	)
	stop := make(chan struct{})
	go controller.Run(stop)
	for {
		time.Sleep(time.Second * time.Duration(config.Period))
	}
}
