// Copyright (c) David Bond, Tailscale Inc, & Contributors
// SPDX-License-Identifier: MIT

package tailscale

import (
	"context"
	"net/http"
)

// ServicesResource provides access to https://tailscale.com/api#tag/services.
type ServicesResource struct {
	*Client
}

// Service is a Tailscale service with a stable virtual IP address.
type Service struct {
	Name        string            `json:"name,omitempty"`
	Addrs       []string          `json:"addrs,omitempty"`
	Comment     string            `json:"comment,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Ports       []string          `json:"ports,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
}

// ServiceHost contains host details for a device advertising a Service.
type ServiceHost struct {
	StableNodeID  string `json:"stableNodeID,omitempty"`
	ApprovalLevel string `json:"approvalLevel,omitempty"`
	Configured    string `json:"configured,omitempty"`
}

// ServiceApproval contains the approval status for a Service on a device.
type ServiceApproval struct {
	Approved     bool `json:"approved"`
	AutoApproved bool `json:"autoApproved"`
}

type serviceList struct {
	Services []Service `json:"vipServices"`
}

type serviceHostsList struct {
	Hosts []ServiceHost `json:"hosts"`
}

type updateServiceApprovalRequest struct {
	Approved bool `json:"approved"`
}

// List lists every [Service] in the tailnet.
func (sr *ServicesResource) List(ctx context.Context) ([]Service, error) {
	req, err := sr.buildRequest(ctx, http.MethodGet, sr.buildTailnetURL("services"))
	if err != nil {
		return nil, err
	}

	resp, err := body[serviceList](sr, req)
	if err != nil {
		return nil, err
	}
	return resp.Services, nil
}

// Get retrieves a specific [Service] by name.
func (sr *ServicesResource) Get(ctx context.Context, name string) (*Service, error) {
	req, err := sr.buildRequest(ctx, http.MethodGet, sr.buildTailnetURL("services", name))
	if err != nil {
		return nil, err
	}

	return body[Service](sr, req)
}

// CreateOrUpdate creates or updates a [Service].
func (sr *ServicesResource) CreateOrUpdate(ctx context.Context, svc Service) error {
	_, err := sr.Upsert(ctx, svc.Name, svc)
	return err
}

// Upsert creates or updates a [Service] using the current resource name in the request path.
// This allows callers to rename an existing Service by providing the current path name separately
// from the desired service name in the request body.
func (sr *ServicesResource) Upsert(ctx context.Context, serviceName string, svc Service) (*Service, error) {
	req, err := sr.buildRequest(ctx, http.MethodPut, sr.buildTailnetURL("services", serviceName), requestBody(svc))
	if err != nil {
		return nil, err
	}

	return body[Service](sr, req)
}

// Delete deletes a specific [Service].
func (sr *ServicesResource) Delete(ctx context.Context, name string) error {
	req, err := sr.buildRequest(ctx, http.MethodDelete, sr.buildTailnetURL("services", name))
	if err != nil {
		return err
	}

	return sr.do(req, nil)
}

// ListHosts lists all devices hosting the specified [Service].
func (sr *ServicesResource) ListHosts(ctx context.Context, serviceName string) ([]ServiceHost, error) {
	req, err := sr.buildRequest(ctx, http.MethodGet, sr.buildTailnetURL("services", serviceName, "devices"))
	if err != nil {
		return nil, err
	}

	resp, err := body[serviceHostsList](sr, req)
	if err != nil {
		return nil, err
	}
	return resp.Hosts, nil
}

// GetDeviceApproval retrieves the approval status for the specified [Service] on a device.
func (sr *ServicesResource) GetDeviceApproval(ctx context.Context, serviceName, deviceID string) (*ServiceApproval, error) {
	req, err := sr.buildRequest(ctx, http.MethodGet, sr.buildTailnetURL("services", serviceName, "device", deviceID, "approved"))
	if err != nil {
		return nil, err
	}

	return body[ServiceApproval](sr, req)
}

// UpdateDeviceApproval updates the approval status for the specified [Service] on a device.
func (sr *ServicesResource) UpdateDeviceApproval(ctx context.Context, serviceName, deviceID string, approved bool) (*ServiceApproval, error) {
	req, err := sr.buildRequest(ctx, http.MethodPost, sr.buildTailnetURL("services", serviceName, "device", deviceID, "approved"), requestBody(updateServiceApprovalRequest{Approved: approved}))
	if err != nil {
		return nil, err
	}

	return body[ServiceApproval](sr, req)
}

// VIPService is an alias for [Service].
// Deprecated: use [Service] instead.
type VIPService = Service

// VIPServicesResource is an alias for [ServicesResource].
// Deprecated: use [ServicesResource] instead.
type VIPServicesResource = ServicesResource

// VIPServiceApproval is an alias for [ServiceApproval].
// Deprecated: use [ServiceApproval] instead.
type VIPServiceApproval = ServiceApproval
