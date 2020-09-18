package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	home "k8s.io/client-go/util/homedir"
)

func main() {
	kubeconfig := flag.String("kubeconfig", fmt.Sprintf("%s/%s", home.HomeDir(), filepath.Join(".kube", "config")), "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	pod, err := clientset.CoreV1().Pods("book").Get(context.Background(), "example", metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println(pod.Status.PodIP)
}
