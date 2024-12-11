package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// Important: Run "make" to regenerate code after modifying this file
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BusinessHoursScalerSpec defines the desired state of BusinessHoursScaler
type BusinessHoursScalerSpec struct {
	// DeploymentSelector specifies the label selector to match Deployments which should be scaled
	DeploymentSelector metav1.LabelSelector `json:"deploymentSelector"`

	// MinReplicas defines the number of replicas during off-hours
	MinReplicas int32 `json:"minReplicas"`

	// MaxReplicas defines the number of replicas during business hours
	MaxReplicas int32 `json:"maxReplicas"`

	// StartTime is the time when business hours begin, in HH:mm:ss format
	StartTime string `json:"startTime"`

	// EndTime is the time when business hours end, in HH:mm:ss format
	EndTime string `json:"endTime"`

	// The time zone name for the given schedule, see https://en.wikipedia.org/wiki/List_of_tz_database_time_zones.
	// +default:value="Etc/UTC"
	TimeZone string `json:"timeZone,omitempty"`
}

// BusinessHoursScalerStatus defines the observed state of BusinessHoursScaler
type BusinessHoursScalerStatus struct {
	// CurrentReplicas shows the current number of replicas
	CurrentReplicas int32 `json:"currentReplicas"`

	// LastUpdated is the last time the operator reconciled
	LastUpdated metav1.Time `json:"lastUpdated"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=bhs

// BusinessHoursScaler is the Schema for the businesshoursscalers API
type BusinessHoursScaler struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BusinessHoursScalerSpec   `json:"spec,omitempty"`
	Status BusinessHoursScalerStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BusinessHoursScalerList contains a list of BusinessHoursScaler
type BusinessHoursScalerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BusinessHoursScaler `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BusinessHoursScaler{}, &BusinessHoursScalerList{})
}
