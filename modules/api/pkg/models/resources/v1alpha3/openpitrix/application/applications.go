/*
Copyright 2019 The KubeSphere Authors.

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
package application

import (
	"kubesphere.io/openpitrix/pkg/api"
	"kubesphere.io/openpitrix/pkg/api/application/v1alpha1"
	query2 "kubesphere.io/openpitrix/pkg/apiserver/query"
	"kubesphere.io/openpitrix/pkg/client/informers/externalversions"
	"kubesphere.io/openpitrix/pkg/models/resources/v1alpha3"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
)

type helmApplicationsGetter struct {
	informers externalversions.SharedInformerFactory
}

func New(informers externalversions.SharedInformerFactory) v1alpha3.Interface {
	return &helmApplicationsGetter{
		informers: informers,
	}
}

func (r *helmApplicationsGetter) Get(_, name string) (runtime.Object, error) {
	app, err := r.informers.Application().V1alpha1().HelmApplications().Lister().Get(name)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return app, nil
}

func (r *helmApplicationsGetter) List(_ string, query *query2.Query) (*api.ListResult, error) {
	var apps []*v1alpha1.HelmApplication
	var err error

	apps, err = r.informers.Application().V1alpha1().HelmApplications().Lister().List(query.Selector())

	if err != nil {
		klog.Error(err)
		return nil, err
	}
	var result []runtime.Object
	for i := range apps {
		result = append(result, apps[i])
	}

	return v1alpha3.DefaultList(result, query, r.compare, r.filter), nil
}

func (r *helmApplicationsGetter) compare(left runtime.Object, right runtime.Object, field query2.Field) bool {

	leftApp, ok := left.(*v1alpha1.HelmApplication)
	if !ok {
		return false
	}

	rightApp, ok := right.(*v1alpha1.HelmApplication)
	if !ok {
		return false
	}
	switch field {
	case query2.FieldName:
		return strings.Compare(leftApp.Spec.Name, rightApp.Spec.Name) > 0
	default:
		return v1alpha3.DefaultObjectMetaCompare(leftApp.ObjectMeta, rightApp.ObjectMeta, field)
	}
}

func (r *helmApplicationsGetter) filter(object runtime.Object, filter query2.Filter) bool {
	application, ok := object.(*v1alpha1.HelmApplication)
	if !ok {
		return false
	}

	switch filter.Field {
	case query2.FieldName:
		return strings.Contains(application.Spec.Name, string(filter.Value))
	case query2.FieldStatus:
		return strings.Contains(application.Status.State, string(filter.Value))
	default:
		return v1alpha3.DefaultObjectMetaFilter(application.ObjectMeta, filter)
	}
}
