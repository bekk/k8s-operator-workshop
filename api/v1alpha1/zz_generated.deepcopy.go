//go:build !ignore_autogenerated

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessHoursScaler) DeepCopyInto(out *BusinessHoursScaler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessHoursScaler.
func (in *BusinessHoursScaler) DeepCopy() *BusinessHoursScaler {
	if in == nil {
		return nil
	}
	out := new(BusinessHoursScaler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessHoursScaler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessHoursScalerList) DeepCopyInto(out *BusinessHoursScalerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]BusinessHoursScaler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessHoursScalerList.
func (in *BusinessHoursScalerList) DeepCopy() *BusinessHoursScalerList {
	if in == nil {
		return nil
	}
	out := new(BusinessHoursScalerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *BusinessHoursScalerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessHoursScalerSpec) DeepCopyInto(out *BusinessHoursScalerSpec) {
	*out = *in
	in.DeploymentSelector.DeepCopyInto(&out.DeploymentSelector)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessHoursScalerSpec.
func (in *BusinessHoursScalerSpec) DeepCopy() *BusinessHoursScalerSpec {
	if in == nil {
		return nil
	}
	out := new(BusinessHoursScalerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BusinessHoursScalerStatus) DeepCopyInto(out *BusinessHoursScalerStatus) {
	*out = *in
	in.LastUpdated.DeepCopyInto(&out.LastUpdated)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BusinessHoursScalerStatus.
func (in *BusinessHoursScalerStatus) DeepCopy() *BusinessHoursScalerStatus {
	if in == nil {
		return nil
	}
	out := new(BusinessHoursScalerStatus)
	in.DeepCopyInto(out)
	return out
}
