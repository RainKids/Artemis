package k8s

import (
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type Clients struct {
	ClientMap map[string]*Client
	logger    *zap.Logger
}

type Client struct {
	KubeConf   *clientcmdapi.Config
	RestConf   *rest.Config
	KubeClient *kubernetes.Clientset
}
