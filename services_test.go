// Copyright (c) David Bond, Tailscale Inc, & Contributors
// SPDX-License-Identifier: MIT

package tailscale

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_ListServices(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expected := []Service{
		{
			Name:    "svc:my-service",
			Addrs:   []string{"100.64.0.1", "fd7a:115c:a1e0::1"},
			Comment: "test service",
			Ports:   []string{"tcp:443"},
			Tags:    []string{"tag:web"},
		},
	}
	server.ResponseBody = serviceList{Services: expected}

	actual, err := client.Services().List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services", server.Path)
	assert.Equal(t, expected, actual)
}

func TestClient_GetService(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expected := &Service{
		Name:    "svc:my-service",
		Addrs:   []string{"100.64.0.1", "fd7a:115c:a1e0::1"},
		Comment: "test service",
		Ports:   []string{"tcp:443"},
		Tags:    []string{"tag:web"},
	}
	server.ResponseBody = expected

	actual, err := client.Services().Get(context.Background(), "svc:my-service")
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services/svc:my-service", server.Path)
	assert.Equal(t, expected, actual)
}

func TestClient_CreateOrUpdateService(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK
	server.ResponseBody = &Service{Name: "svc:my-service"}

	svc := Service{
		Name:    "svc:my-service",
		Comment: "new service",
		Ports:   []string{"tcp:443"},
		Tags:    []string{"tag:web"},
	}

	err := client.Services().CreateOrUpdate(context.Background(), svc)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPut, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services/svc:my-service", server.Path)

	var received Service
	err = json.Unmarshal(server.Body.Bytes(), &received)
	assert.NoError(t, err)
	assert.Equal(t, svc, received)
}

func TestClient_UpsertService(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expected := &Service{
		Name:    "svc:renamed-service",
		Comment: "renamed service",
		Ports:   []string{"tcp:443"},
	}
	server.ResponseBody = expected

	actual, err := client.Services().Upsert(context.Background(), "svc:old-service", *expected)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPut, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services/svc:old-service", server.Path)
	assert.Equal(t, expected, actual)

	var received Service
	err = json.Unmarshal(server.Body.Bytes(), &received)
	assert.NoError(t, err)
	assert.Equal(t, *expected, received)
}

func TestClient_DeleteService(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	err := client.Services().Delete(context.Background(), "svc:my-service")
	assert.NoError(t, err)
	assert.Equal(t, http.MethodDelete, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services/svc:my-service", server.Path)
}

func TestClient_ListServiceHosts(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expected := []ServiceHost{{
		StableNodeID:  "n292kg92CNTRL",
		ApprovalLevel: "approved:manual",
		Configured:    "ready",
	}}
	server.ResponseBody = serviceHostsList{Hosts: expected}

	actual, err := client.Services().ListHosts(context.Background(), "svc:my-service")
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services/svc:my-service/devices", server.Path)
	assert.Equal(t, expected, actual)
}

func TestClient_GetServiceDeviceApproval(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expected := &ServiceApproval{Approved: true, AutoApproved: false}
	server.ResponseBody = expected

	actual, err := client.Services().GetDeviceApproval(context.Background(), "svc:my-service", "n123")
	assert.NoError(t, err)
	assert.Equal(t, http.MethodGet, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services/svc:my-service/device/n123/approved", server.Path)
	assert.Equal(t, expected, actual)
}

func TestClient_UpdateServiceDeviceApproval(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK

	expected := &ServiceApproval{Approved: true, AutoApproved: false}
	server.ResponseBody = expected

	actual, err := client.Services().UpdateDeviceApproval(context.Background(), "svc:my-service", "n123", true)
	assert.NoError(t, err)
	assert.Equal(t, http.MethodPost, server.Method)
	assert.Equal(t, "/api/v2/tailnet/example.com/services/svc:my-service/device/n123/approved", server.Path)
	assert.Equal(t, expected, actual)

	var received updateServiceApprovalRequest
	err = json.Unmarshal(server.Body.Bytes(), &received)
	assert.NoError(t, err)
	assert.True(t, received.Approved)
}

func TestClient_GetService_NotFound(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusNotFound
	server.ResponseBody = APIError{Message: "not found"}

	_, err := client.Services().Get(context.Background(), "svc:nonexistent")
	assert.Error(t, err)
	assert.True(t, IsNotFound(err))
}

func TestClient_VIPServicesAlias(t *testing.T) {
	t.Parallel()

	client, server := NewTestHarness(t)
	server.ResponseCode = http.StatusOK
	server.ResponseBody = serviceList{}

	_, err := client.VIPServices().List(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "/api/v2/tailnet/example.com/services", server.Path)
}
