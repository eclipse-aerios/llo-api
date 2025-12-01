package config

import (
	"log"
	"path/filepath"

	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type K8sClient struct {
	ClientSet *kubernetes.Clientset
	Dynamic   *dynamic.DynamicClient
}

var Client K8sClient

func GetKubernetesClient() error {
	log.Println("Loading K8s client from a local kubeconfig file")
	// TODO configure kubeconfig file path (env vars, ...)?
	var kubeconfig string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}
	log.Println("Using kubeconfig file: " + kubeconfig)
	// Use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}
	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	// Create the dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	Client.ClientSet = clientset
	Client.Dynamic = dynamicClient

	return nil
}

func GetKubernetesClientInCluster() error {
	log.Println("Loading K8s client config from in-cluster config")
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}
	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err
	}
	// creates the clientset
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return err
	}

	// link constant client to created clients
	Client.ClientSet = clientset
	Client.Dynamic = dynamicClient

	return err
}

func init() {
	log.Println("Initializing K8s clientset and dynamic client...")
	err := GetKubernetesClientInCluster()

	if err != nil {
		log.Println(err)
		log.Println("The K8s client config cannot be obtained from in-cluster config, so assuming that the application is being run from outside the cluster")

		err := GetKubernetesClient()
		if err != nil {
			log.Println(err)
			Status = UNHEALTHY_STATUS
			log.Println("The K8s client cannot be initialized. Setting status to " + UNHEALTHY_STATUS + "...")
			// panic(err.Error())
		}
	}
}
