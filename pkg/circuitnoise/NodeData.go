package circuitnoise

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"k8s.io/klog/v2"
)

// Assigning map literal
var nodeDNS = map[string]string{
    "minikube-m02":  "10.105.22.156",
    "minikube-m03": "10.100.241.135",
}

func GetNodeNoise(nodeName string) (float64, error) {
    // Get the port number from the nodePort map
    serviceUrl, ok := nodeDNS[nodeName]
    if !ok {
        klog.Infof("[CircuitNoise] Error in getting node service")
        return 0, fmt.Errorf("nodeName %s not found in nodeDNS map", nodeName)
    }

    // Make a network call to retrieve JSON data
    url := fmt.Sprintf("http://%s:8000/get-noise/", serviceUrl)
    response, err := http.Get(url)
    if err != nil {
        klog.Infof("[CircuitNoise] Failed to make HTTP request")
        return 0, fmt.Errorf("failed to make HTTP request: %v", err)
    }
    defer response.Body.Close()

    // Read the response body
    body, err := io.ReadAll(response.Body)
    if err != nil {
        klog.Infof("[CircuitNoise] Failed to read response body")
        return 0, fmt.Errorf("failed to read response body: %v", err)
    }

    // Define a struct to hold the JSON response
    var responseData struct {
        Result float64 `json:"result"`
    }

    // Unmarshal JSON data
    if err := json.Unmarshal(body, &responseData); err != nil {
        klog.Infof("[CircuitNoise] Failed to unmarshall json")
        return 0, fmt.Errorf("failed to unmarshal JSON: %v", err)
    }

    return responseData.Result, nil
}