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
package helmrelease

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

type helmReleasesGetter struct {
	informers externalversions.SharedInformerFactory
}

func New(informers externalversions.SharedInformerFactory) v1alpha3.Interface {
	return &helmReleasesGetter{
		informers: informers,
	}
}

func (r *helmReleasesGetter) Get(_, name string) (runtime.Object, error) {
	app, err := r.informers.Application().V1alpha1().HelmReleases().Lister().Get(name)
	if err != nil {
		klog.Error(err)
		return nil, err
	}
	return app, nil
}

func (r *helmReleasesGetter) List(_ string, query *query2.Query) (*api.ListResult, error) {
	var rls []*v1alpha1.HelmRelease
	var err error

	rls, err = r.informers.Application().V1alpha1().HelmReleases().Lister().List(query.Selector())

	if err != nil {
		klog.Error(err)
		return nil, err
	}
	var result []runtime.Object
	for i := range rls {
		result = append(result, rls[i])
	}

	return v1alpha3.DefaultList(result, query, r.compare, r.filter), nil
}

func (r *helmReleasesGetter) compare(left runtime.Object, right runtime.Object, field query2.Field) bool {

	leftRls, ok := left.(*v1alpha1.HelmRelease)
	if !ok {
		return false
	}

	rightRls, ok := right.(*v1alpha1.HelmRelease)
	if !ok {
		return false
	}
	switch field {
	case query2.FieldName:
		return strings.Compare(leftRls.Spec.Name, rightRls.Spec.Name) > 0
	default:
		return v1alpha3.DefaultObjectMetaCompare(leftRls.ObjectMeta, rightRls.ObjectMeta, field)
	}
}

func (r *helmReleasesGetter) filter(object runtime.Object, filter query2.Filter) bool {
	rls, ok := object.(*v1alpha1.HelmRelease)
	if !ok {
		return false
	}

	switch filter.Field {
	case query2.FieldName:
		return strings.Contains(rls.Spec.Name, string(filter.Value))
	case query2.FieldStatus:
		return strings.Contains(rls.Status.State, string(filter.Value))
	default:
		return v1alpha3.DefaultObjectMetaFilter(rls.ObjectMeta, filter)
	}
}
