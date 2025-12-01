package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/eclipse-aerios/llo-api/config"
	"github.com/eclipse-aerios/llo-api/models"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ServiceComponentSvc struct{}

func (s *ServiceComponentSvc) GetDeployedServiceComponents(kind string) ([]unstructured.Unstructured, error) {
	dynamicClient := config.Client.Dynamic
	scGVR := createScGroupVersionResource(kind)
	serviceComponents, err := dynamicClient.Resource(scGVR).Namespace(config.Namespace).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("There are %d Service Components in the cluster\n", len(serviceComponents.Items))
	return serviceComponents.Items, nil
}

func (s *ServiceComponentSvc) GetOnlyServiceComponentsIds(serviceComponents []unstructured.Unstructured) (ids []string, err error) {
	ids = []string{}
	for i := 0; i < len(serviceComponents); i++ {
		ids = append(ids, serviceComponents[i].GetName())
	}
	return
}

func (s *ServiceComponentSvc) GetDeployedServiceComponent(kind string, name string) (*unstructured.Unstructured, error) {
	dynamicClient := config.Client.Dynamic
	scGVR := createScGroupVersionResource(kind)
	serviceComponent, err := dynamicClient.Resource(scGVR).Namespace(config.Namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Found Service Component " + name + " of type " + kind + " in the cluster")
	return serviceComponent, nil
}

func (s *ServiceComponentSvc) DeleteServiceComponent(kind string, name string) error {
	dynamicClient := config.Client.Dynamic
	scGVR := createScGroupVersionResource(kind)
	err := dynamicClient.Resource(scGVR).Namespace(config.Namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Deleted CR " + name)
	return nil
}

func (s *ServiceComponentSvc) DeployToKubernetes(sc *models.ServiceComponent) error {
	dynamicClient := config.Client.Dynamic
	scGVR := createScGroupVersionResource(config.GetCR(sc.Kind))
	// Serialize ServiceComponent to JSON
	data, err := json.Marshal(sc)
	if err != nil {
		log.Println("Error serializing Service Component CR to JSON")
		return err
	}
	// Convert JSON to unstructured.Unstructured
	var unstructuredObj unstructured.Unstructured
	if err := unstructuredObj.UnmarshalJSON(data); err != nil {
		log.Println("Error unmarshalling JSON to Unstructured")
		return err
	}
	log.Printf("Creating ServiceComponent with APIVersion: %s and Kind: %s\n", unstructuredObj.GetAPIVersion(), unstructuredObj.GetKind())
	// Create the Custom Resource using the dynamic client
	_, err = dynamicClient.Resource(scGVR).Namespace(config.Namespace).Create(context.Background(), &unstructuredObj, metav1.CreateOptions{})
	if err != nil {
		log.Printf("Error creating custom resource: %v", err)
		return err
	}
	log.Println("Created custom resource: " + sc.GetObjectMeta().GetName())
	return nil
}

func (s *ServiceComponentSvc) PatchServiceComponent(kind string, name string, sc *models.ServiceComponentSpec) error {
	dynamicClient := config.Client.Dynamic
	scGVR := createScGroupVersionResource(kind)
	// Get SC from K8s
	clusterSc, err := dynamicClient.Resource(scGVR).Namespace(config.Namespace).Get(context.Background(), name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Found Service Component " + clusterSc.GetName() + " of type " + clusterSc.GetKind() + " in the cluster")
	// Create a new ServiceComponent from the function input params and from the former CR obtained from K8s
	updatedSc := &models.ServiceComponent{
		ObjectMeta: metav1.ObjectMeta{
			Name:            name,
			Labels:          clusterSc.GetLabels(),
			ResourceVersion: clusterSc.GetResourceVersion(),
		},
		TypeMeta: metav1.TypeMeta{
			APIVersion: clusterSc.GetAPIVersion(),
			Kind:       clusterSc.GetKind(),
		},
		ApiVersionAux: clusterSc.GetAPIVersion(),
		Spec:          *sc,
	}
	// Serialize ServiceComponent to JSON
	data, err := json.Marshal(updatedSc)
	if err != nil {
		log.Println("Error serializing Service Component CR to JSON")
		return err
	}
	// Convert JSON to unstructured.Unstructured
	var unstructuredObj unstructured.Unstructured
	if err := unstructuredObj.UnmarshalJSON(data); err != nil {
		log.Println("Error unmarshalling JSON to Unstructured")
		return err
	}
	log.Printf("Updating ServiceComponent with APIVersion: %s and Kind: %s\n", unstructuredObj.GetAPIVersion(), unstructuredObj.GetKind())
	// Update the Custom Resource using the dynamic client
	_, err = dynamicClient.Resource(scGVR).Namespace(config.Namespace).Update(context.Background(), &unstructuredObj, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Error updating custom resource: %v", err)
		return err
	}
	log.Println("Updated custom resource: " + name)
	return nil
}

func (s *ServiceComponentSvc) UpdateServiceComponent(sc *models.ServiceComponent) error {
	dynamicClient := config.Client.Dynamic
	scGVR := createScGroupVersionResource(config.GetCR(sc.Kind))
	// Get SC from K8s
	clusterSc, err := dynamicClient.Resource(scGVR).Namespace(config.Namespace).Get(context.Background(), sc.Name, metav1.GetOptions{})
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Found Service Component " + clusterSc.GetName() + " of type " + clusterSc.GetKind() + " in the cluster")
	resourceVersion := clusterSc.GetResourceVersion()
	sc.ObjectMeta.ResourceVersion = resourceVersion
	// Serialize ServiceComponent to JSON
	data, err := json.Marshal(sc)
	if err != nil {
		log.Println("Error serializing Service Component CR to JSON")
		return err
	}
	// Convert JSON to unstructured.Unstructured
	var unstructuredObj unstructured.Unstructured
	if err := unstructuredObj.UnmarshalJSON(data); err != nil {
		log.Println("Error unmarshalling JSON to Unstructured")
		return err
	}
	log.Printf("Updating ServiceComponent with APIVersion: %s and Kind: %s\n", unstructuredObj.GetAPIVersion(), unstructuredObj.GetKind())
	// Update the Custom Resource using the dynamic client
	_, err = dynamicClient.Resource(scGVR).Namespace(config.Namespace).Update(context.Background(), &unstructuredObj, metav1.UpdateOptions{})
	if err != nil {
		log.Printf("Error updating custom resource: %v", err)
		return err
	}
	log.Println("Updated custom resource: " + sc.GetObjectMeta().GetName())
	return nil
}

func createScGroupVersionResource(kind string) (scGVR schema.GroupVersionResource) {
	scGVR = schema.GroupVersionResource{
		Group:    config.GVR_GROUP,
		Version:  config.GVR_VERSION,
		Resource: kind,
	}
	return
}
