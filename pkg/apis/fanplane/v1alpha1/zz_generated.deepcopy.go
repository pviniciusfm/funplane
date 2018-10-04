// +build !ignore_autogenerated

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConnectionPoolSettings) DeepCopyInto(out *ConnectionPoolSettings) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		if *in == nil {
			*out = nil
		} else {
			*out = new(ReadableDuration)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConnectionPoolSettings.
func (in *ConnectionPoolSettings) DeepCopy() *ConnectionPoolSettings {
	if in == nil {
		return nil
	}
	out := new(ConnectionPoolSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConsulServiceEntry) DeepCopyInto(out *ConsulServiceEntry) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConsulServiceEntry.
func (in *ConsulServiceEntry) DeepCopy() *ConsulServiceEntry {
	if in == nil {
		return nil
	}
	out := new(ConsulServiceEntry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DNSRoute) DeepCopyInto(out *DNSRoute) {
	*out = *in
	if in.Type != nil {
		in, out := &in.Type, &out.Type
		if *in == nil {
			*out = nil
		} else {
			*out = new(DNSType)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DNSRoute.
func (in *DNSRoute) DeepCopy() *DNSRoute {
	if in == nil {
		return nil
	}
	out := new(DNSRoute)
	in.DeepCopyInto(out)
	return out
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvoyBootstrap.
func (in *EnvoyBootstrap) DeepCopy() *EnvoyBootstrap {
	if in == nil {
		return nil
	}
	out := new(EnvoyBootstrap)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EnvoyBootstrap) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvoyBootstrapList) DeepCopyInto(out *EnvoyBootstrapList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]EnvoyBootstrap, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvoyBootstrapList.
func (in *EnvoyBootstrapList) DeepCopy() *EnvoyBootstrapList {
	if in == nil {
		return nil
	}
	out := new(EnvoyBootstrapList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *EnvoyBootstrapList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FaultInjection) DeepCopyInto(out *FaultInjection) {
	*out = *in
	if in.AbortPercent != nil {
		in, out := &in.AbortPercent, &out.AbortPercent
		if *in == nil {
			*out = nil
		} else {
			*out = new(float32)
			**out = **in
		}
	}
	if in.Delay != nil {
		in, out := &in.Delay, &out.Delay
		if *in == nil {
			*out = nil
		} else {
			*out = new(ReadableDuration)
			**out = **in
		}
	}
	if in.DelayPercent != nil {
		in, out := &in.DelayPercent, &out.DelayPercent
		if *in == nil {
			*out = nil
		} else {
			*out = new(float32)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FaultInjection.
func (in *FaultInjection) DeepCopy() *FaultInjection {
	if in == nil {
		return nil
	}
	out := new(FaultInjection)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Gateway) DeepCopyInto(out *Gateway) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Status = in.Status
	if in.Spec != nil {
		in, out := &in.Spec, &out.Spec
		if *in == nil {
			*out = nil
		} else {
			*out = new(GatewaySpec)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Gateway.
func (in *Gateway) DeepCopy() *Gateway {
	if in == nil {
		return nil
	}
	out := new(Gateway)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Gateway) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GatewayList) DeepCopyInto(out *GatewayList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Gateway, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GatewayList.
func (in *GatewayList) DeepCopy() *GatewayList {
	if in == nil {
		return nil
	}
	out := new(GatewayList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GatewayList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GatewaySpec) DeepCopyInto(out *GatewaySpec) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Hosts != nil {
		in, out := &in.Hosts, &out.Hosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Listener != nil {
		in, out := &in.Listener, &out.Listener
		if *in == nil {
			*out = nil
		} else {
			*out = new(Listener)
			**out = **in
		}
	}
	if in.Routes != nil {
		in, out := &in.Routes, &out.Routes
		*out = make([]*Route, len(*in))
		for i := range *in {
			if (*in)[i] == nil {
				(*out)[i] = nil
			} else {
				(*out)[i] = new(Route)
				(*in)[i].DeepCopyInto((*out)[i])
			}
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GatewaySpec.
func (in *GatewaySpec) DeepCopy() *GatewaySpec {
	if in == nil {
		return nil
	}
	out := new(GatewaySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Listener) DeepCopyInto(out *Listener) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Listener.
func (in *Listener) DeepCopy() *Listener {
	if in == nil {
		return nil
	}
	out := new(Listener)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancerSettings) DeepCopyInto(out *LoadBalancerSettings) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancerSettings.
func (in *LoadBalancerSettings) DeepCopy() *LoadBalancerSettings {
	if in == nil {
		return nil
	}
	out := new(LoadBalancerSettings)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PrefixRewrite) DeepCopyInto(out *PrefixRewrite) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PrefixRewrite.
func (in *PrefixRewrite) DeepCopy() *PrefixRewrite {
	if in == nil {
		return nil
	}
	out := new(PrefixRewrite)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RetryPolicy) DeepCopyInto(out *RetryPolicy) {
	*out = *in
	if in.RetryOn != nil {
		in, out := &in.RetryOn, &out.RetryOn
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PerTryTimeout != nil {
		in, out := &in.PerTryTimeout, &out.PerTryTimeout
		if *in == nil {
			*out = nil
		} else {
			*out = new(ReadableDuration)
			**out = **in
		}
	}
	if in.MaxTimeout != nil {
		in, out := &in.MaxTimeout, &out.MaxTimeout
		if *in == nil {
			*out = nil
		} else {
			*out = new(ReadableDuration)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RetryPolicy.
func (in *RetryPolicy) DeepCopy() *RetryPolicy {
	if in == nil {
		return nil
	}
	out := new(RetryPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Route) DeepCopyInto(out *Route) {
	*out = *in
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		if *in == nil {
			*out = nil
		} else {
			*out = new(ConsulServiceEntry)
			**out = **in
		}
	}
	if in.DNS != nil {
		in, out := &in.DNS, &out.DNS
		if *in == nil {
			*out = nil
		} else {
			*out = new(DNSRoute)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.Rewrite != nil {
		in, out := &in.Rewrite, &out.Rewrite
		if *in == nil {
			*out = nil
		} else {
			*out = new(PrefixRewrite)
			**out = **in
		}
	}
	if in.TLSContext != nil {
		in, out := &in.TLSContext, &out.TLSContext
		if *in == nil {
			*out = nil
		} else {
			*out = new(TLSContext)
			**out = **in
		}
	}
	if in.LoadBalancerSettings != nil {
		in, out := &in.LoadBalancerSettings, &out.LoadBalancerSettings
		if *in == nil {
			*out = nil
		} else {
			*out = new(LoadBalancerSettings)
			**out = **in
		}
	}
	if in.ConnectionPoolSettings != nil {
		in, out := &in.ConnectionPoolSettings, &out.ConnectionPoolSettings
		if *in == nil {
			*out = nil
		} else {
			*out = new(ConnectionPoolSettings)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.RetryPolicy != nil {
		in, out := &in.RetryPolicy, &out.RetryPolicy
		if *in == nil {
			*out = nil
		} else {
			*out = new(RetryPolicy)
			(*in).DeepCopyInto(*out)
		}
	}
	if in.FaultInjection != nil {
		in, out := &in.FaultInjection, &out.FaultInjection
		if *in == nil {
			*out = nil
		} else {
			*out = new(FaultInjection)
			(*in).DeepCopyInto(*out)
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Route.
func (in *Route) DeepCopy() *Route {
	if in == nil {
		return nil
	}
	out := new(Route)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Status.
func (in *Status) DeepCopy() *Status {
	if in == nil {
		return nil
	}
	out := new(Status)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TLSContext) DeepCopyInto(out *TLSContext) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TLSContext.
func (in *TLSContext) DeepCopy() *TLSContext {
	if in == nil {
		return nil
	}
	out := new(TLSContext)
	in.DeepCopyInto(out)
	return out
}
