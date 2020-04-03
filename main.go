package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/sanity-io/litter"
	flag "github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	var labels []string
	var gvks []string

	flag.StringArrayVarP(&labels, "labels", "l", nil, "labels to select objects from a cluster")
	flag.StringArrayVar(&gvks, "gvk", nil, "GroupVersionKinds to list by labels")

	flag.Parse()

	parsedGvks := []schema.GroupVersionKind{}

	for _, gvkStr := range gvks {
		parts := strings.Split(gvkStr, ".")
		if len(parts) != 3 {
			fmt.Printf("failed to parse 3 GVK tokens from flag: %s\n", gvkStr)
			os.Exit(1)
		}

		found := schema.GroupVersionKind{
			Group:   parts[0],
			Version: parts[1],
			Kind:    parts[2],
		}

		parsedGvks = append(parsedGvks, found)
	}

	parsedLabels := map[string]string{}

	for _, labelStr := range labels {
		parts := strings.Split(labelStr, "=")
		if len(parts) != 2 {
			fmt.Printf("failed to parse key=value tokens from flag: %s\n", labelStr)
			os.Exit(1)
		}

		parsedLabels[parts[0]] = parts[1]
	}

	kubeconfig, err := ctrl.GetConfig()
	if err != nil {
		fmt.Printf("failed to get kubeconfig: %v\n", err)
		os.Exit(1)
	}

	kubeclient, err := client.New(kubeconfig, client.Options{})
	if err != nil {
		fmt.Printf("failed to create client: %v\n", err)
		os.Exit(1)
	}

	result := []runtime.Object{}
	for _, gvk := range parsedGvks {
		obj := new(unstructured.UnstructuredList)
		obj.SetGroupVersionKind(gvk)
		if err := kubeclient.List(context.Background(), obj, client.MatchingLabels(parsedLabels)); err != nil {
			fmt.Printf("failed to list %s, error: %v\n", gvk.String(), err)
			os.Exit(1)
		}
		for _, item := range obj.Items {
			result = append(result, &item)
		}
	}

	fmt.Println("printing results")
	litter.Dump(result)
}
