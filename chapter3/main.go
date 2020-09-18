package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	home "k8s.io/client-go/util/homedir"
	wq "k8s.io/client-go/util/workqueue"
)

func main() {
	kubeconfig := flag.String("kubeconfig", fmt.Sprintf("%s/%s", home.HomeDir(), filepath.Join(".kube", "config")), "kubeconfig file")
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	// config.AcceptContentTypes = "application/vnd.kubernetes.protobuf,application/json"
	// config.ContentType = "application/vnd.kubernetes.protobuf"
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	q := wq.New()
	informerFactory := informers.NewFilteredSharedInformerFactory(clientset, 30*time.Second, "book", func(opts *metav1.ListOptions) {})
	podInformer := informerFactory.Core().V1().Pods()
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			if pod, ok := newObj.(*v1.Pod); ok {
				q.Add(pod.Name)
			}
		},
	})
	go func() {
		for {
			n, s := q.Get()
			if s {
				break
			}
			fmt.Println(n)
			q.Done(n)
		}
	}()
	informerFactory.Start(wait.NeverStop)
	informerFactory.WaitForCacheSync(wait.NeverStop)
	fmt.Println(clientset.Discovery().ServerVersion())
	pods, err := podInformer.Lister().Pods("book").List(labels.Everything())
	if err != nil {
		panic(err)
	}
	pod := pods[0]
	if err != nil {
		panic(err)
	}
	res, err := restmapper.GetAPIGroupResources(clientset.Discovery())
	if err != nil {
		panic(err)
	}
	m := restmapper.NewDiscoveryRESTMapper(res)
	t, err := m.RESTMapping(schema.ParseGroupKind("Deployment.apps"))
	if err != nil {
		panic(err)
	}
	fmt.Println(t.GroupVersionKind)
	fmt.Println(t.Resource)
	fmt.Println(pod.GetObjectKind().GroupVersionKind().Empty())
	fmt.Println(pod.Name)
	fmt.Println(pod.Status.PodIP)
	select {}
}
