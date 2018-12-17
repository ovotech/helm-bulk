/*
Copyright The Helm Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

/*
configForContext and getKubeClient funcs copied, at time of writing, from
https://github.com/helm/helm/blob/1b34a511d4ae38e43518e99a8250330515e3a93c/cmd/helm/helm.go
 without any modifications.
*/

package utils

import (
	"fmt"
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	//https://github.com/helm/helm/issues/3806#issuecomment-378072159

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/helm/portforwarder"
	kube "k8s.io/helm/pkg/kube"
	"k8s.io/helm/pkg/tlsutil"
)

//Client creates a Helm client and checks the connection works
func Client() (client *helm.Client) {
	config, kclient, err := getKubeClient("", "")
	PanicCheck(err)
	log.Println("Portforwarding from Kubernetes Tiller pod")
	pf, err := portforwarder.New("kube-system", kclient, config)
	PanicCheck(err)
	log.Println("Starting Helm client")
	tillerHost := fmt.Sprintf("127.0.0.1:%d", pf.Local)
	options := []helm.Option{
		helm.Host(tillerHost),
	}
	helmHome := os.Getenv("HELM_HOME")
	opts := tlsutil.Options{
		ServerName:         "",
		CaCertFile:         helmHome + "/ca.pem",
		CertFile:           helmHome + "/ca.cert.pem",
		KeyFile:            helmHome + "/ca.key.pem",
		InsecureSkipVerify: false,
	}
	tlsCfg, err := tlsutil.ClientConfig(opts)
	PanicCheck(err)
	options = append(options, helm.WithTLS(tlsCfg))
	client = helm.NewClient(options...)
	log.Println("Checking Helm client connection")
	_, errb := client.GetVersion()
	PanicCheck(errb)
	log.Println("Connection: OK")
	return
}

// configForContext creates a Kubernetes REST client configuration for a
// given kubeconfig context.
func configForContext(context string, kubeconfig string) (*rest.Config, error) {
	config, err := kube.GetConfig(context, kubeconfig).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("could not get Kubernetes config for context %q: %s",
			context, err)
	}
	return config, nil
}

// getKubeClient creates a Kubernetes config and client for a given kubeconfig
// context.
func getKubeClient(context string, kubeconfig string) (*rest.Config,
	kubernetes.Interface, error) {
	config, err := configForContext(context, kubeconfig)
	if err != nil {
		return nil, nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("could not get Kubernetes client: %s", err)
	}
	return config, client, nil
}
