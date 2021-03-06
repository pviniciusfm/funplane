// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	scheme "github.frg.tech/cloud/fanplane/pkg/apis/client/clientset/versioned/scheme"
	v1alpha1 "github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// GatewaysGetter has a method to return a GatewayInterface.
// A group's client should implement this interface.
type GatewaysGetter interface {
	Gateways(namespace string) GatewayInterface
}

// GatewayInterface has methods to work with Gateway resources.
type GatewayInterface interface {
	Create(*v1alpha1.Gateway) (*v1alpha1.Gateway, error)
	Update(*v1alpha1.Gateway) (*v1alpha1.Gateway, error)
	UpdateStatus(*v1alpha1.Gateway) (*v1alpha1.Gateway, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Gateway, error)
	List(opts v1.ListOptions) (*v1alpha1.GatewayList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Gateway, err error)
	GatewayExpansion
}

// gateways implements GatewayInterface
type gateways struct {
	client rest.Interface
	ns     string
}

// newGateways returns a Gateways
func newGateways(c *FanplaneV1alpha1Client, namespace string) *gateways {
	return &gateways{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the gateway, and returns the corresponding gateway object, and an error if there is any.
func (c *gateways) Get(name string, options v1.GetOptions) (result *v1alpha1.Gateway, err error) {
	result = &v1alpha1.Gateway{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("gateways").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Gateways that match those selectors.
func (c *gateways) List(opts v1.ListOptions) (result *v1alpha1.GatewayList, err error) {
	result = &v1alpha1.GatewayList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("gateways").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested gateways.
func (c *gateways) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("gateways").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a gateway and creates it.  Returns the server's representation of the gateway, and an error, if there is any.
func (c *gateways) Create(gateway *v1alpha1.Gateway) (result *v1alpha1.Gateway, err error) {
	result = &v1alpha1.Gateway{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("gateways").
		Body(gateway).
		Do().
		Into(result)
	return
}

// Update takes the representation of a gateway and updates it. Returns the server's representation of the gateway, and an error, if there is any.
func (c *gateways) Update(gateway *v1alpha1.Gateway) (result *v1alpha1.Gateway, err error) {
	result = &v1alpha1.Gateway{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("gateways").
		Name(gateway.Name).
		Body(gateway).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *gateways) UpdateStatus(gateway *v1alpha1.Gateway) (result *v1alpha1.Gateway, err error) {
	result = &v1alpha1.Gateway{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("gateways").
		Name(gateway.Name).
		SubResource("status").
		Body(gateway).
		Do().
		Into(result)
	return
}

// Delete takes name of the gateway and deletes it. Returns an error if one occurs.
func (c *gateways) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("gateways").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *gateways) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("gateways").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched gateway.
func (c *gateways) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Gateway, err error) {
	result = &v1alpha1.Gateway{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("gateways").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
