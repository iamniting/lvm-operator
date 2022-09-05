/*
Copyright 2021.

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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var lvmclusterlog = logf.Log.WithName("lvmcluster-webhook")

func (r *LVMCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-lvm-topolvm-io-v1alpha1-lvmcluster,mutating=true,failurePolicy=fail,sideEffects=None,groups=lvm.topolvm.io,resources=lvmclusters,verbs=create;update,versions=v1alpha1,name=mlvmcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &LVMCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *LVMCluster) Default() {
	lvmclusterlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-lvm-topolvm-io-v1alpha1-lvmcluster,mutating=false,failurePolicy=fail,sideEffects=None,groups=lvm.topolvm.io,resources=lvmclusters,verbs=create;update,versions=v1alpha1,name=vlvmcluster.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &LVMCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *LVMCluster) ValidateCreate() error {
	lvmclusterlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *LVMCluster) ValidateUpdate(old runtime.Object) error {
	lvmclusterlog.Info("validate update", "name", r.Name)

	oldLVMCluster, ok := old.(*LVMCluster)
	if !ok {
		return fmt.Errorf("Failed to parse LVMCluster.")
	}

	for _, deviceClass := range r.Spec.Storage.DeviceClasses {
		var newDevices, oldDevices []string

		if deviceClass.DeviceSelector != nil {
			newDevices = deviceClass.DeviceSelector.Paths
		}

		oldDevices = oldLVMCluster.GetPathsOfDeviceClass(deviceClass.Name)

		// if devices are removed now
		if len(oldDevices) > len(newDevices) {
			return fmt.Errorf("Invalid:devices can not be removed from the LVMCluster once added.")
		}

		// if devices are added now
		if len(oldDevices) == 0 && len(newDevices) > 0 {
			return fmt.Errorf("Invalid:devices can not be added in the LVMCluster once created without devices.")
		}

		deviceMap := make(map[string]bool)

		for _, device := range oldDevices {
			deviceMap[device] = true
		}

		for _, device := range newDevices {
			delete(deviceMap, device)
		}

		// if any old device is removed now
		if len(deviceMap) != 0 {
			return fmt.Errorf("Invalid:some of devices are deleted from the LVMCluster. "+
				"Device can not be removed from the LVMCluster once added. "+
				"oldDevices:%s, newDevices:%s", oldDevices, newDevices)
		}
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *LVMCluster) ValidateDelete() error {
	lvmclusterlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *LVMCluster) GetPathsOfDeviceClass(deviceClassName string) []string {

	for _, deviceClass := range r.Spec.Storage.DeviceClasses {
		if deviceClass.Name == deviceClassName {
			if deviceClass.DeviceSelector != nil {
				return deviceClass.DeviceSelector.Paths
			}
			return []string{}
		}
	}

	return []string{}
}
