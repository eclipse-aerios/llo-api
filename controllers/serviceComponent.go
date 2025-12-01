package controllers

import (
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/eclipse-aerios/llo-api/config"
	"github.com/eclipse-aerios/llo-api/models"
	"github.com/eclipse-aerios/llo-api/services"
	"github.com/gin-gonic/gin"
	"k8s.io/apimachinery/pkg/api/errors"
)

type ServiceComponentController struct {
	// Deploy(c *gin.Context)
	// GetByName(c *gin.Context, name string)
	// List(c *gin.Context)
	// Delete(c *gin.Context)
	svc services.ServiceComponentSvc
}

func (s *ServiceComponentController) Deploy(c *gin.Context) {
	if c.ContentType() != "application/yaml" && c.ContentType() != "text/yaml" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The content of the body must be in YAML format"})
		return
	}
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body cannot be empty"})
		return
	}

	// var serviceComponent *models.ServiceComponent
	// serviceComponent := new(models.ServiceComponent)
	serviceComponent := &models.ServiceComponent{}
	if err := c.ShouldBindYAML(&serviceComponent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the Kind is a valid K8s Custom Resource
	if !slices.Contains(config.GetSupportedCRs(), serviceComponent.TypeMeta.Kind) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This K8s CR is not currently supported by aeriOS LLOs"})
		return
	}

	log.Println("Orchestration type: " + strings.Split(serviceComponent.Kind, "ServiceComponent")[1])
	log.Println("Selected IE: " + serviceComponent.Spec.SelectedIE.Id + " Hostname: " + serviceComponent.Spec.SelectedIE.Hostname)

	// Deploy to Kubernetes
	log.Println("Deploying CR into the K8s cluster...")
	if err := s.svc.DeployToKubernetes(serviceComponent); err != nil {
		if errors.IsAlreadyExists(err) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Successful ServiceComponent deployment"})
	// c.Status(http.StatusCreated)
}

func (s *ServiceComponentController) List(c *gin.Context) {
	lloType := config.GetCR(c.DefaultQuery("type", config.DEFAULT_LLO_TYPE))
	onlyIds, err := strconv.ParseBool(c.DefaultQuery("onlyNames", "false"))
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "The onlyNames parameter must be true or false"})
		return
	}
	if lloType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This LLO type is not valid"})
		return
	}
	serviceComponents, err := s.svc.GetDeployedServiceComponents(lloType)
	if errors.IsNotFound(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "CRD not found in the K8s cluster"})
		return
	}
	c.Header("Results-Count", strconv.Itoa(len(serviceComponents)))
	c.Header("LLO-Type", lloType)
	if onlyIds {
		serviceComponentsIds, _ := s.svc.GetOnlyServiceComponentsIds(serviceComponents)
		c.JSON(http.StatusOK, serviceComponentsIds)
	} else {
		c.JSON(http.StatusOK, serviceComponents)
	}
}

func (s *ServiceComponentController) GetByName(c *gin.Context) {
	lloType := config.GetCR(c.DefaultQuery("type", config.DEFAULT_LLO_TYPE))
	if lloType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This LLO type is not valid"})
		return
	}
	serviceComponents, err := s.svc.GetDeployedServiceComponent(lloType, c.Param("scName"))
	if errors.IsNotFound(err) {
		if errors.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": c.Param("scName") + " of type " + lloType + " ServiceComponent not found in the K8s cluster"})
			return
		}
		return
	}
	c.JSON(http.StatusOK, serviceComponents)
}

func (s *ServiceComponentController) Delete(c *gin.Context) {
	lloType := config.GetCR(c.DefaultQuery("type", config.DEFAULT_LLO_TYPE))
	if lloType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This LLO type is not valid"})
		return
	}
	err := s.svc.DeleteServiceComponent(lloType, c.Param("scName"))
	if errors.IsNotFound(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": c.Param("scName") + " of type " + lloType + " ServiceComponent not found in the K8s cluster"})
		return
	}
	c.Status(http.StatusOK)
}

func (s *ServiceComponentController) Patch(c *gin.Context) {
	if c.ContentType() != "application/yaml" && c.ContentType() != "text/yaml" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The content of the body must be in YAML format"})
		return
	}
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body cannot be empty"})
		return
	}

	lloType := config.GetCR(c.DefaultQuery("type", config.DEFAULT_LLO_TYPE))
	if lloType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This LLO type is not valid"})
		return
	}

	serviceComponent := &models.ServiceComponentSpec{}
	if err := c.ShouldBindYAML(&serviceComponent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Println("Orchestration type: " + lloType)
	log.Println("Selected IE: " + serviceComponent.SelectedIE.Id + " Hostname: " + serviceComponent.SelectedIE.Hostname)

	// Update in Kubernetes
	log.Println("Deploying CR into the K8s cluster...")
	if err := s.svc.PatchServiceComponent(lloType, c.Param("scName"), serviceComponent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Successful ServiceComponent update"})
}

func (s *ServiceComponentController) Update(c *gin.Context) {
	if c.ContentType() != "application/yaml" && c.ContentType() != "text/yaml" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "The content of the body must be in YAML format"})
		return
	}
	if c.Request.Body == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Body cannot be empty"})
		return
	}

	serviceComponent := &models.ServiceComponent{}
	if err := c.ShouldBindYAML(&serviceComponent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the Kind is a valid K8s Custom Resource
	if !slices.Contains(config.GetSupportedCRs(), serviceComponent.TypeMeta.Kind) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This K8s CR is not currently supported by aeriOS LLOs"})
		return
	}

	log.Println("Orchestration type: " + strings.Split(serviceComponent.Kind, "ServiceComponent")[1])
	log.Println("Selected IE: " + serviceComponent.Spec.SelectedIE.Id + " Hostname: " + serviceComponent.Spec.SelectedIE.Hostname)

	// Update in Kubernetes
	log.Println("Deploying CR into the K8s cluster...")
	if err := s.svc.UpdateServiceComponent(serviceComponent); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, gin.H{"message": "Successful ServiceComponent update"})
}
