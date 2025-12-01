package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ServiceComponent struct {
	ApiVersionAux     string `yaml:"apiVersion,omitempty" json:"apiVersion,omitempty"`
	metav1.TypeMeta   `json:",inline" yaml:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec              ServiceComponentSpec `json:"spec" yaml:"spec"`
}

type ServiceComponentSpec struct {
	SelectedIE     SelectedIE     `yaml:"selectedIE,omitempty" json:"selectedIE,omitempty"`
	Image          string         `yaml:"image,omitempty" json:"image,omitempty"`
	ImageRegistry  ImageRegistry  `yaml:"imageRegistry,omitempty" json:"imageRegistry,omitempty"`
	IsJob          bool           `yaml:"isJob,omitempty" json:"isJob,omitempty"`
	Ports          []NetworkPort  `yaml:"ports,omitempty" json:"ports,omitempty"`
	ExposePorts    bool           `yaml:"exposePorts,omitempty" json:"exposePorts,omitempty"`
	EnvVars        []KeyValue     `yaml:"envVars,omitempty" json:"envVars,omitempty"`
	CliArgs        []KeyValue     `yaml:"cliArgs,omitempty" json:"cliArgs,omitempty"`
	Privileged     bool           `yaml:"privileged,omitempty" json:"privileged,omitempty"`
	NetworkOverlay NetworkOverlay `yaml:"networkOverlay,omitempty" json:"networkOverlay,omitempty"`
}

type ImageRegistry struct {
	Url      string `yaml:"url,omitempty" json:"url,omitempty"`
	Username string `yaml:"username,omitempty" json:"username,omitempty"`
	Password string `yaml:"password,omitempty" json:"password,omitempty"`
}

type KeyValue struct {
	Key   string `yaml:"key,omitempty" json:"key,omitempty"`
	Value string `yaml:"value,omitempty" json:"value,omitempty"`
}

type NetworkOverlay struct {
	InternalIp          string `yaml:"internalIp,omitempty" json:"internalIp,omitempty"`
	Dns                 string `yaml:"dns,omitempty" json:"dns,omitempty"`
	PublicKey           string `yaml:"publicKey,omitempty" json:"publicKey,omitempty"`
	PrivateKey          string `yaml:"privateKey,omitempty" json:"privateKey,omitempty"`
	Endpoint            string `yaml:"endpoint,omitempty" json:"endpoint,omitempty"`
	Mtu                 int    `yaml:"mtu,omitempty" json:"mtu,omitempty"`
	AllowedIps          string `yaml:"allowedIps,omitempty" json:"allowedIps,omitempty"`
	PersistentKeepalive int    `yaml:"persistentKeepalive,omitempty" json:"persistentKeepalive,omitempty"`
}

type NetworkPort struct {
	Number   int32  `yaml:"number,omitempty" json:"number,omitempty"`
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

type SelectedIE struct {
	Id       string `yaml:"id,omitempty" json:"id,omitempty"`
	Hostname string `yaml:"hostname,omitempty" json:"hostname,omitempty"`
}
