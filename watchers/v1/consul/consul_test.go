package consul

import (
	"os"
	"sync"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/testutil"
	"github.com/stretchr/testify/assert"

	"github.com/rapid7/cps/pkg/kv"
)

func TestGetServiceHealth(t *testing.T) {
	os.Stdout, _ = os.Open(os.DevNull)

	srv1, err := testutil.NewTestServer()
	if err != nil {
		t.Fatal(err)
	}
	defer srv1.Stop()

	srv1.AddAddressableService(t, "service-one", api.HealthPassing, "127.0.0.1", 8192, []string{"test"})
	srv1.AddCheck(t, "service:service-one", "service-one", api.HealthPassing)

	client, err := setUpConsulClient(srv1.HTTPAddr)
	assert.Nil(t, err, "Expected no error setting up consul client")

	services, qo, err := getServices(client)
	assert.Nil(t, err, "Expected no error getting services")

	healthyNodes = make(map[string][]string)
	var mutex = &sync.Mutex{}
	for key, _ := range services {
		getServiceHealth(key, client, qo, mutex)
	}

	writeProperties()

	em := map[string][]string{
		"service-one": {"127.0.0.1"},
		"consul":      {"127.0.0.1"},
	}
	c := kv.GetProperty("consul")

	assert.Equal(t, c, em, "Expected consul maps to be equal")
}
