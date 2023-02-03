package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"k8s.io/klog"
	"k8s.io/kubernetes/pkg/apis/core/v1"
	"k8s.io/kubernetes/pkg/controller/podautoscaler/metrics"
)

const (
	doToken = "YOUR_DIGITAL_OCEAN_API_TOKEN"
	busy    = "busy"
)

type stats struct {
	Status string `json:"status"`
}

func deleteDroplet(client *godo.Client, dropletID string) error {
	if _, err := client.Droplets.Delete(context.TODO(), dropletID); err != nil {
		return err
	}
	return nil
}

func main() {
	client := godo.NewClient(nil)
	client.Authenticator = &godo.TokenAuthenticator{
		AccessToken: doToken,
	}

	pods, err := metrics.GetPodsForMetrics(context.TODO(), v1.NamespaceAll)
	if err != nil {
		klog.Errorf("Failed to get pods: %v", err)
		return
	}

	var busyPods []*v1.Pod
	for i, pod := range pods.Items {
		response, err := http.Get(fmt.Sprintf("http://%s:8001/v1/stats", pod.Status.PodIP))
		if err != nil {
			klog.Errorf("Failed to get pod stats: %v", err)
			continue
		}

		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			klog.Errorf("Failed to read response body: %v", err)
			continue
		}

		var s stats
		if err := json.Unmarshal(body, &s); err != nil {
			klog.Errorf("Failed to parse response body: %v", err)
			continue
		}

		if s.Status == busy {
			busyPods = append(busyPods, &pod)
			continue
		}

		time.Sleep(5 * time.Minute)

		if err := deleteDroplet(client, strconv.Itoa(i)); err != nil {
			klog.Errorf("Failed to delete droplet: %v", err)
			continue
		}
	}

	if len(busyPods) == len(pods.Items) {
		klog.Infof("All pods are busy, skipping node scale down")
		return
	}

	klog.Infof("Sleeping for %v before removing node", wait)
		time.Sleep(wait)

	if err := deleteDroplet(client, strconv.Itoa(i)); err != nil {
			klog.Errorf("Failed to delete droplet: %v", err)
			continue
		}
	}
