// Copyright 2018 OVO Technology
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"log"
	"os"

	"k8s.io/helm/pkg/helm"
	"k8s.io/helm/pkg/tlsutil"
)

//Client creates a Helm client and checks the connection works
func Client(tlsKey, tlsCert, caCert, tlsServerName string, disableTLS bool) (client *helm.Client) {
	options := []helm.Option{
		helm.Host(os.Getenv("TILLER_HOST")),
	}
	if !disableTLS {
		if tlsServerName == "" {
			panic("If using TLS, serverName must be set. This is the value of '/O='" +
				" that you used in the subject field when creating the CSR.")
		}
		opts := tlsutil.Options{
			ServerName:         tlsServerName,
			CaCertFile:         caCert,
			CertFile:           tlsCert,
			KeyFile:            tlsKey,
			InsecureSkipVerify: false,
		}
		tlsCfg, err := tlsutil.ClientConfig(opts)
		PanicCheck(err)
		options = append(options, helm.WithTLS(tlsCfg))
	}

	client = helm.NewClient(options...)
	log.Println("Checking Helm client connection")
	_, errb := client.GetVersion()
	PanicCheck(errb)
	log.Println("Connection: OK")
	return
}
