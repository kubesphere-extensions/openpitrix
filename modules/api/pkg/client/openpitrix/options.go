/*
Copyright 2020 KubeSphere Authors

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

package openpitrix

import (
	"kubesphere.io/openpitrix/pkg/client/s3"
	"kubesphere.io/openpitrix/pkg/utils/reflectutils"
	"time"

	"github.com/spf13/pflag"
)

type Options struct {
	S3Options                *s3.Options               `json:"s3,omitempty" yaml:"s3,omitempty" mapstructure:"s3"`
	ReleaseControllerOptions *ReleaseControllerOptions `json:"releaseControllerOptions,omitempty" yaml:"releaseControllerOptions,omitempty" mapstructure:"releaseControllerOptions"`
}

type ReleaseControllerOptions struct {
	MaxConcurrent int           `json:"maxConcurrent,omitempty" yaml:"maxConcurrent,omitempty" mapstructure:"maxConcurrent"`
	WaitTime      time.Duration `json:"waitTime,omitempty" yaml:"waitTime,omitempty" mapstructure:"waitTime"`
}

func NewOptions() *Options {
	return &Options{
		S3Options: &s3.Options{},
		ReleaseControllerOptions: &ReleaseControllerOptions{
			MaxConcurrent: 10,
			WaitTime:      30 * time.Second,
		},
	}
}

// Validate check options values
func (s *Options) Validate() []error {
	var errors []error

	return errors
}

func (s *Options) AppStoreConfIsEmpty() bool {
	return s.S3Options == nil || s.S3Options.Endpoint == ""
}

// ApplyTo overrides options if it's valid, which endpoint is not empty
func (s *Options) ApplyTo(options *Options) {
	if s.S3Options != nil {
		reflectutils.Override(options, s)
	}

	if s.ReleaseControllerOptions != nil {
		reflectutils.Override(options, s)
	}
}

// AddFlags add options flags to command line flags,
func (s *Options) AddFlags(fs *pflag.FlagSet, c *Options) {
	// if s3-endpoint if left empty, following options will be ignored
	fs.StringVar(&s.S3Options.Endpoint, "openpitrix-s3-endpoint", c.S3Options.Endpoint, ""+
		"Endpoint to access to s3 object storage service for openpitrix, if left blank, the following options "+
		"will be ignored.")

	fs.StringVar(&s.S3Options.Region, "openpitrix-s3-region", c.S3Options.Region, ""+
		"Region of s3 that openpitrix will access to, like us-east-1.")

	fs.StringVar(&s.S3Options.AccessKeyID, "openpitrix-s3-access-key-id", c.S3Options.AccessKeyID, "access key of openpitrix s3")

	fs.StringVar(&s.S3Options.SecretAccessKey, "openpitrix-s3-secret-access-key", c.S3Options.SecretAccessKey, "secret access key of openpitrix s3")

	fs.StringVar(&s.S3Options.SessionToken, "openpitrix-s3-session-token", c.S3Options.SessionToken, "session token of openpitrix s3")

	fs.StringVar(&s.S3Options.Bucket, "openpitrix-s3-bucket", c.S3Options.Bucket, "bucket name of openpitrix s3")

	fs.BoolVar(&s.S3Options.DisableSSL, "openpitrix-s3-disable-SSL", c.S3Options.DisableSSL, "disable ssl")

	fs.BoolVar(&s.S3Options.ForcePathStyle, "openpitrix-s3-force-path-style", c.S3Options.ForcePathStyle, "force path style")
	fs.StringVar(&s.S3Options.FilePath, "s3-file-path", c.S3Options.FilePath, "file save path")
	fs.DurationVar(&s.ReleaseControllerOptions.WaitTime, "openpitrix-release-controller-options-wait-time", c.ReleaseControllerOptions.WaitTime, "wait time when check release is ready or not")
	fs.IntVar(&s.ReleaseControllerOptions.MaxConcurrent, "openpitrix-release-controller-options-max-concurrent", c.ReleaseControllerOptions.MaxConcurrent, "the maximum number of concurrent Reconciles which can be run for release controller")
}
