// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.frg.tech/cloud/fanplane/pkg/apis/fanplane/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEnvoyBootstraps implements EnvoyBootstrapInterface
type FakeEnvoyBootstraps struct {
	Fake *FakeFanplaneV1alpha1
	ns   string
}

var envoybootstrapsResource = schema.GroupVersionResource{Group: "fanplane.io", Version: "v1alpha1", Resource: "envoybootstraps"}

var envoybootstrapsKind = schema.GroupVersionKind{Group: "fanplane.io", Version: "v1alpha1", Kind: "EnvoyBootstrap"}

// Get takes name of the envoyBootstrap, and returns the corresponding envoyBootstrap object, and an error if there is any.
func (c *FakeEnvoyBootstraps) Get(name string, options v1.GetOptions) (result *v1alpha1.EnvoyBootstrap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(envoybootstrapsResource, c.ns, name), &v1alpha1.EnvoyBootstrap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EnvoyBootstrap), err
}

// List takes label and field selectors, and returns the list of EnvoyBootstraps that match those selectors.
func (c *FakeEnvoyBootstraps) List(opts v1.ListOptions) (result *v1alpha1.EnvoyBootstrapList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(envoybootstrapsResource, envoybootstrapsKind, c.ns, opts), &v1alpha1.EnvoyBootstrapList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EnvoyBootstrapList{}
	for _, item := range obj.(*v1alpha1.EnvoyBootstrapList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested envoyBootstraps.
func (c *FakeEnvoyBootstraps) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(envoybootstrapsResource, c.ns, opts))

}

// Create takes the representation of a envoyBootstrap and creates it.  Returns the server's representation of the envoyBootstrap, and an error, if there is any.
func (c *FakeEnvoyBootstraps) Create(envoyBootstrap *v1alpha1.EnvoyBootstrap) (result *v1alpha1.EnvoyBootstrap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(envoybootstrapsResource, c.ns, envoyBootstrap), &v1alpha1.EnvoyBootstrap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EnvoyBootstrap), err
}

// Update takes the representation of a envoyBootstrap and updates it. Returns the server's representation of the envoyBootstrap, and an error, if there is any.
func (c *FakeEnvoyBootstraps) Update(envoyBootstrap *v1alpha1.EnvoyBootstrap) (result *v1alpha1.EnvoyBootstrap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(envoybootstrapsResource, c.ns, envoyBootstrap), &v1alpha1.EnvoyBootstrap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EnvoyBootstrap), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeEnvoyBootstraps) UpdateStatus(envoyBootstrap *v1alpha1.EnvoyBootstrap) (*v1alpha1.EnvoyBootstrap, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(envoybootstrapsResource, "status", c.ns, envoyBootstrap), &v1alpha1.EnvoyBootstrap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EnvoyBootstrap), err
}

// Delete takes name of the envoyBootstrap and deletes it. Returns an error if one occurs.
func (c *FakeEnvoyBootstraps) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(envoybootstrapsResource, c.ns, name), &v1alpha1.EnvoyBootstrap{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEnvoyBootstraps) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(envoybootstrapsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.EnvoyBootstrapList{})
	return err
}

// Patch applies the patch and returns the patched envoyBootstrap.
func (c *FakeEnvoyBootstraps) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EnvoyBootstrap, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(envoybootstrapsResource, c.ns, name, data, subresources...), &v1alpha1.EnvoyBootstrap{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EnvoyBootstrap), err
}
