// Copyright 2020 MongoDB Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// +build unit

package store

import (
	"testing"

	"github.com/mongodb/mongocli/internal/config"
	"go.mongodb.org/atlas/mongodbatlas"
)

type auth struct {
	username string
	password string
}

func (a auth) PublicAPIKey() string {
	return a.username
}

func (a auth) PrivateAPIKey() string {
	return a.password
}

var _ CredentialsGetter = &auth{}

func TestService(t *testing.T) {
	c, err := New(Service(config.CloudService))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if c.service != config.CloudService {
		t.Errorf("New() service = %s; expected %s", c.service, "cloud")
	}
}

func TestWithBaseURL(t *testing.T) {
	c, err := New(Service(config.CloudService), WithBaseURL("http://test"))
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if c.baseURL != "http://test" {
		t.Errorf("New() baseURL = %s; expected %s", c.baseURL, "http://test")
	}
}

func TestSkipVerify(t *testing.T) {
	c, err := New(Service(config.CloudService), SkipVerify())
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if !c.skipVerify {
		t.Error("New() skipVerify not set")
	}
}

func TestWithPublicPathBaseURL(t *testing.T) {
	c, err := New(Service(config.CloudService), WithBaseURL("http://test"), WithPublicPathBaseURL())
	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	expected := "http://test" + mongodbatlas.APIPublicV1Path
	if c.baseURL != expected {
		t.Errorf("New() baseURL = %s; expected %s", c.baseURL, expected)
	}
}

type testConfig struct {
	url string
	auth
}

func (c testConfig) OpsManagerCACertificate() string {
	return ""
}

func (c testConfig) OpsManagerSkipVerify() string {
	return "false"
}

func (c testConfig) Service() string {
	return config.CloudService
}

func (c testConfig) OpsManagerURL() string {
	return c.url
}

var _ AuthenticatedConfig = &testConfig{}

func TestPrivateAuthenticatedPreset(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c, err := New(PrivateAuthenticatedPreset(testConfig{}))
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		if c.baseURL != mongodbatlas.CloudURL {
			t.Errorf("New() baseURL = %s; expected %s", c.baseURL, mongodbatlas.CloudURL)
		}
	})
	t.Run("with a base url", func(t *testing.T) {
		const url = "http://test"
		c, err := New(PrivateAuthenticatedPreset(testConfig{url: url}))
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		if c.baseURL != url {
			t.Errorf("New() baseURL = %s; expected %s", c.baseURL, url)
		}
	})
}

func TestPrivateUnauthenticatedPreset(t *testing.T) {
	t.Run("default", func(t *testing.T) {
		c, err := New(PrivateUnauthenticatedPreset(testConfig{}))
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		if c.baseURL != mongodbatlas.CloudURL {
			t.Errorf("New() baseURL = %s; expected %s", c.baseURL, mongodbatlas.CloudURL)
		}
	})
	t.Run("with a base url", func(t *testing.T) {
		const url = "http://test"
		c, err := New(PrivateUnauthenticatedPreset(testConfig{url: url}))
		if err != nil {
			t.Fatalf("New() unexpected error: %v", err)
		}

		if c.baseURL != url {
			t.Errorf("New() baseURL = %s; expected %s", c.baseURL, url)
		}
	})
}

func TestWithAuthentication(t *testing.T) {
	a := auth{
		username: "username",
		password: "password",
	}
	c, err := New(Service("cloud"), WithAuthentication(a))

	if err != nil {
		t.Fatalf("New() unexpected error: %v", err)
	}

	if c.username != a.username {
		t.Errorf("New() username = %s; expected %s", c.username, a.username)
	}
	if c.password != a.password {
		t.Errorf("New() password = %s; expected %s", c.password, a.password)
	}
}
