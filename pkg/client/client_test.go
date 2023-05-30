/*
Copyright 2021 The Everoute Authors.

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

package client_test

import (
	"os"
	"testing"

	"github.com/everoute/everoute/plugin/tower/pkg/server/fake"
	. "github.com/onsi/gomega"
)

var (
	server *fake.Server
)

func TestMain(m *testing.M) {
	server = fake.NewServer(nil)
	server.Serve()

	os.Exit(m.Run())
}

func TestClient_Query(t *testing.T) {
	RegisterTestingT(t)
	// todo
}

func TestClient_Subscription(t *testing.T) {
	RegisterTestingT(t)
	// todo
}
