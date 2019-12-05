/*
Copyright 2019 The Tekton Authors
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

package utils

import (
	"context"
	"crypto/tls"
	restful "github.com/emicklei/go-restful"
	logging "github.com/tektoncd/dashboard/pkg/logging"
	"golang.org/x/oauth2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"
	"net/http"
	"strings"
)

// RespondError - logs and writes an error response with a desired status code
func RespondError(response *restful.Response, err error, statusCode int) {
	logging.Log.Error("Error: ", strings.Replace(err.Error(), "/", "", -1))
	response.AddHeader("Content-Type", "text/plain")
	response.WriteError(statusCode, err)
}

// RespondErrorMessage - logs and writes an error message with a desired status code
func RespondErrorMessage(response *restful.Response, message string, statusCode int) {
	logging.Log.Debugf("Error message: %s", message)
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(statusCode, message)
}

// RespondMessageAndLogError - logs and writes an error message with a desired status code and logs the error
func RespondMessageAndLogError(response *restful.Response, err error, message string, statusCode int) {
	logging.Log.Error("Error: ", strings.Replace(err.Error(), "/", "", -1))
	logging.Log.Debugf("Message: %s", message)
	response.AddHeader("Content-Type", "text/plain")
	response.WriteErrorString(statusCode, message)
}

// createOAuth2Client returns an HTTP client with oauth2 authentication using the provided accessToken
func CreateOAuth2Client(ctx context.Context, accessToken string) *http.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	return oauth2.NewClient(ctx, ts)
}

func GetClientAllowsSelfSigned() *http.Client {
	transport := &http.Transport{
		// ToDo - investigate InsecureSkipVerify here
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: transport}
	return client
}

// getWebhookSecretTokens returns the "secretToken" and "accessToken" stored in the Secret
// with the name specified by the parameter, and in the namespace specified by r.Defaults.Namespace.
func GetWebhookSecretTokens(kubeClient k8sclient.Interface, namespace, name string) (accessToken string, secretToken string, err error) {
	// Access token is stored as 'accessToken' and secret as 'secretToken'
	secret, err := kubeClient.CoreV1().Secrets(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return "", "", err
	}
	accessToken = string(secret.Data["accessToken"])
	secretToken = string(secret.Data["secretToken"])
	return accessToken, secretToken, nil
}
