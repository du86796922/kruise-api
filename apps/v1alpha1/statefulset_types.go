/*
Copyright 2020 The Kruise Authors.

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

package v1alpha1

import (
	appspub "github.com/openkruise/kruise-api/apps/pub"
	apps "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// MaxMinReadySeconds is the max value of MinReadySeconds
	MaxMinReadySeconds = 300
)

// StatefulSetUpdateStrategy indicates the strategy that the StatefulSet
// controller will use to perform updates. It includes any additional parameters
// necessary to perform the update for the indicated strategy.
type StatefulSetUpdateStrategy struct {
	// Type indicates the type of the StatefulSetUpdateStrategy.
	// Default is RollingUpdate.
	// +optional
	Type apps.StatefulSetUpdateStrategyType `json:"type,omitempty"`
	// RollingUpdate is used to communicate parameters when Type is RollingUpdateStatefulSetStrategyType.
	// +optional
	RollingUpdate *RollingUpdateStatefulSetStrategy `json:"rollingUpdate,omitempty"`
}

// RollingUpdateStatefulSetStrategy is used to communicate parameter for RollingUpdateStatefulSetStrategyType.
type RollingUpdateStatefulSetStrategy struct {
	// Partition indicates the ordinal at which the StatefulSet should be partitioned by default.
	// But if unorderedUpdate has been set:
	//   - Partition indicates the number of pods with non-updated revisions when rolling update.
	//   - It means controller will update $(replicas - partition) number of pod.
	// Default value is 0.
	// +optional
	Partition *int32 `json:"partition,omitempty"`
	// The maximum number of pods that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired pods (ex: 10%).
	// Absolute number is calculated from percentage by rounding down.
	// Also, maxUnavailable can just be allowed to work with Parallel podManagementPolicy.
	// Defaults to 1.
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`
	// PodUpdatePolicy indicates how pods should be updated
	// Default value is "ReCreate"
	// +optional
	PodUpdatePolicy PodUpdateStrategyType `json:"podUpdatePolicy,omitempty"`
	// Paused indicates that the StatefulSet is paused.
	// Default value is false
	// +optional
	Paused bool `json:"paused,omitempty"`
	// UnorderedUpdate contains strategies for non-ordered update.
	// If it is not nil, pods will be updated with non-ordered sequence.
	// Noted that UnorderedUpdate can only be allowed to work with Parallel podManagementPolicy
	// +optional
	UnorderedUpdate *UnorderedUpdateStrategy `json:"unorderedUpdate,omitempty"`
	// InPlaceUpdateStrategy contains strategies for in-place update.
	// +optional
	InPlaceUpdateStrategy *appspub.InPlaceUpdateStrategy `json:"inPlaceUpdateStrategy,omitempty"`
	// MinReadySeconds indicates how long will the pod be considered ready after it's updated.
	// MinReadySeconds works with both OrderedReady and Parallel podManagementPolicy.
	// It affects the pod scale up speed when the podManagementPolicy is set to be OrderedReady.
	// Combined with MaxUnavailable, it affects the pod update speed regardless of podManagementPolicy.
	// Default value is 0, max is 300.
	// +optional
	MinReadySeconds *int32 `json:"minReadySeconds,omitempty"`
}

// UnorderedUpdateStrategy defines strategies for non-ordered update.
type UnorderedUpdateStrategy struct {
	// Priorities are the rules for calculating the priority of updating pods.
	// Each pod to be updated, will pass through these terms and get a sum of weights.
	// +optional
	PriorityStrategy *appspub.UpdatePriorityStrategy `json:"priorityStrategy,omitempty"`
}

// PodUpdateStrategyType is a string enumeration type that enumerates
// all possible ways we can update a Pod when updating application
type PodUpdateStrategyType string

const (
	// RecreatePodUpdateStrategyType indicates that we always delete Pod and create new Pod
	// during Pod update, which is the default behavior
	RecreatePodUpdateStrategyType PodUpdateStrategyType = "ReCreate"
	// InPlaceIfPossiblePodUpdateStrategyType indicates that we try to in-place update Pod instead of
	// recreating Pod when possible. Currently, only image update of pod spec is allowed. Any other changes to the pod
	// spec will fall back to ReCreate PodUpdateStrategyType where pod will be recreated.
	InPlaceIfPossiblePodUpdateStrategyType PodUpdateStrategyType = "InPlaceIfPossible"
	// InPlaceOnlyPodUpdateStrategyType indicates that we will in-place update Pod instead of
	// recreating pod. Currently we only allow image update for pod spec. Any other changes to the pod spec will be
	// rejected by kube-apiserver
	InPlaceOnlyPodUpdateStrategyType PodUpdateStrategyType = "InPlaceOnly"
)

// StatefulSetSpec defines the desired state of StatefulSet
type StatefulSetSpec struct {
	// replicas is the desired number of replicas of the given Template.
	// These are replicas in the sense that they are instantiations of the
	// same Template, but individual replicas also have a consistent identity.
	// If unspecified, defaults to 1.
	// TODO: Consider a rename of this field.
	// +optional
	Replicas *int32 `json:"replicas,omitempty"`

	// selector is a label query over pods that should match the replica count.
	// It must match the pod template's labels.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	Selector *metav1.LabelSelector `json:"selector"`

	// template is the object that describes the pod that will be created if
	// insufficient replicas are detected. Each pod stamped out by the StatefulSet
	// will fulfill this Template, but have a unique identity from the rest
	// of the StatefulSet.
	Template v1.PodTemplateSpec `json:"template"`

	// volumeClaimTemplates is a list of claims that pods are allowed to reference.
	// The StatefulSet controller is responsible for mapping network identities to
	// claims in a way that maintains the identity of a pod. Every claim in
	// this list must have at least one matching (by name) volumeMount in one
	// container in the template. A claim in this list takes precedence over
	// any volumes in the template, with the same name.
	// TODO: Define the behavior if a claim already exists with the same name.
	// +optional
	VolumeClaimTemplates []v1.PersistentVolumeClaim `json:"volumeClaimTemplates,omitempty"`

	// serviceName is the name of the service that governs this StatefulSet.
	// This service must exist before the StatefulSet, and is responsible for
	// the network identity of the set. Pods get DNS/hostnames that follow the
	// pattern: pod-specific-string.serviceName.default.svc.cluster.local
	// where "pod-specific-string" is managed by the StatefulSet controller.
	ServiceName string `json:"serviceName,omitempty"`

	// podManagementPolicy controls how pods are created during initial scale up,
	// when replacing pods on nodes, or when scaling down. The default policy is
	// `OrderedReady`, where pods are created in increasing order (pod-0, then
	// pod-1, etc) and the controller will wait until each pod is ready before
	// continuing. When scaling down, the pods are removed in the opposite order.
	// The alternative policy is `Parallel` which will create pods in parallel
	// to match the desired scale without waiting, and on scale down will delete
	// all pods at once.
	// +optional
	PodManagementPolicy apps.PodManagementPolicyType `json:"podManagementPolicy,omitempty"`

	// updateStrategy indicates the StatefulSetUpdateStrategy that will be
	// employed to update Pods in the StatefulSet when a revision is made to
	// Template.
	UpdateStrategy StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`

	// revisionHistoryLimit is the maximum number of revisions that will
	// be maintained in the StatefulSet's revision history. The revision history
	// consists of all revisions not represented by a currently applied
	// StatefulSetSpec version. The default value is 10.
	RevisionHistoryLimit *int32 `json:"revisionHistoryLimit,omitempty"`
}

// StatefulSetStatus defines the observed state of StatefulSet
type StatefulSetStatus struct {
	// observedGeneration is the most recent generation observed for this StatefulSet. It corresponds to the
	// StatefulSet's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`

	// replicas is the number of Pods created by the StatefulSet controller.
	Replicas int32 `json:"replicas"`

	// readyReplicas is the number of Pods created by the StatefulSet controller that have a Ready Condition.
	ReadyReplicas int32 `json:"readyReplicas"`

	// AvailableReplicas is the number of Pods created by the StatefulSet controller that have been ready for
	//minReadySeconds.
	AvailableReplicas int32 `json:"availableReplicas"`

	// currentReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version
	// indicated by currentRevision.
	CurrentReplicas int32 `json:"currentReplicas"`

	// updatedReplicas is the number of Pods created by the StatefulSet controller from the StatefulSet version
	// indicated by updateRevision.
	UpdatedReplicas int32 `json:"updatedReplicas"`

	// currentRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the
	// sequence [0,currentReplicas).
	CurrentRevision string `json:"currentRevision,omitempty"`

	// updateRevision, if not empty, indicates the version of the StatefulSet used to generate Pods in the sequence
	// [replicas-updatedReplicas,replicas)
	UpdateRevision string `json:"updateRevision,omitempty"`

	// collisionCount is the count of hash collisions for the StatefulSet. The StatefulSet controller
	// uses this field as a collision avoidance mechanism when it needs to create the name for the
	// newest ControllerRevision.
	// +optional
	CollisionCount *int32 `json:"collisionCount,omitempty"`

	// Represents the latest available observations of a statefulset's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []apps.StatefulSetCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type"`

	// LabelSelector is label selectors for query over pods that should match the replica count used by HPA.
	LabelSelector string `json:"labelSelector,omitempty"`
}

// These are valid conditions of a statefulset.
const (
	FailedCreatePod apps.StatefulSetConditionType = "FailedCreatePod"
	FailedUpdatePod apps.StatefulSetConditionType = "FailedUpdatePod"
)

// +genclient
// +genclient:method=GetScale,verb=get,subresource=scale,result=k8s.io/api/autoscaling/v1.Scale
// +genclient:method=UpdateScale,verb=update,subresource=scale,input=k8s.io/api/autoscaling/v1.Scale,result=k8s.io/api/autoscaling/v1.Scale
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.labelSelector
// +kubebuilder:resource:shortName=sts;asts
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.replicas,selectorpath=.status.labelSelector
// +kubebuilder:printcolumn:name="DESIRED",type="integer",JSONPath=".spec.replicas",description="The desired number of pods."
// +kubebuilder:printcolumn:name="CURRENT",type="integer",JSONPath=".status.replicas",description="The number of currently all pods."
// +kubebuilder:printcolumn:name="UPDATED",type="integer",JSONPath=".status.updatedReplicas",description="The number of pods updated."
// +kubebuilder:printcolumn:name="READY",type="integer",JSONPath=".status.readyReplicas",description="The number of pods ready."
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp",description="CreationTimestamp is a timestamp representing the server time when this object was created. It is not guaranteed to be set in happens-before order across separate operations. Clients may not set this value. It is represented in RFC3339 form and is in UTC."

// StatefulSet is the Schema for the statefulsets API
type StatefulSet struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StatefulSetSpec   `json:"spec,omitempty"`
	Status StatefulSetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StatefulSetList contains a list of StatefulSet
type StatefulSetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StatefulSet `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StatefulSet{}, &StatefulSetList{})
}
