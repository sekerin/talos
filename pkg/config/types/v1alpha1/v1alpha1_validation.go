// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package v1alpha1 provides user-facing v1alpha1 machine configs
// nolint: dupl
package v1alpha1

import (
	"errors"
	"fmt"
	"net"
	"strconv"

	"github.com/hashicorp/go-multierror"

	"github.com/talos-systems/talos/pkg/config/machine"
)

var (
	// General

	// ErrRequiredSection denotes a section is required
	ErrRequiredSection = errors.New("required config section")
	// ErrInvalidVersion denotes that the config file version is invalid
	ErrInvalidVersion = errors.New("invalid config version")

	// Security

	// ErrInvalidCert denotes that the certificate specified is invalid
	ErrInvalidCert = errors.New("certificate is invalid")
	// ErrInvalidCertType denotes that the certificate type is invalid
	ErrInvalidCertType = errors.New("certificate type is invalid")

	// Services

	// ErrUnsupportedCNI denotes that the specified CNI is invalid
	ErrUnsupportedCNI = errors.New("unsupported CNI driver")
	// ErrInvalidTrustdToken denotes that a trustd token has not been specified
	ErrInvalidTrustdToken = errors.New("trustd token is invalid")

	// Networking

	// ErrBadAddressing denotes that an incorrect combination of network
	// address methods have been specified
	ErrBadAddressing = errors.New("invalid network device addressing method")
	// ErrInvalidAddress denotes that a bad address was provided
	ErrInvalidAddress = errors.New("invalid network address")
)

// NetworkDeviceCheck defines the function type for checks.
// nolint: dupl
type NetworkDeviceCheck func(*machine.Device) error

// Validate triggers the specified validation checks to run.
// nolint: dupl
func Validate(d *machine.Device, checks ...NetworkDeviceCheck) error {
	var result *multierror.Error

	if d.Ignore {
		return result.ErrorOrNil()
	}

	for _, check := range checks {
		result = multierror.Append(result, check(d))
	}

	return result.ErrorOrNil()
}

// CheckDeviceInterface ensures that the interface has been specified.
// nolint: dupl
func CheckDeviceInterface() NetworkDeviceCheck {
	return func(d *machine.Device) error {
		var result *multierror.Error

		if d.Interface == "" {
			result = multierror.Append(result, fmt.Errorf("[%s] %q: %w", "networking.os.device.interface", "", ErrRequiredSection))
		}

		return result.ErrorOrNil()
	}
}

// CheckDeviceAddressing ensures that an appropriate addressing method.
// has been specified
// nolint: dupl
func CheckDeviceAddressing() NetworkDeviceCheck {
	return func(d *machine.Device) error {
		var result *multierror.Error

		// Test for both dhcp and cidr specified
		if d.DHCP && d.CIDR != "" {
			result = multierror.Append(result, fmt.Errorf("[%s] %q: %w", "networking.os.device", "", ErrBadAddressing))
		}

		// test for neither dhcp nor cidr specified
		if !d.DHCP && d.CIDR == "" {
			result = multierror.Append(result, fmt.Errorf("[%s] %q: %w", "networking.os.device", "", ErrBadAddressing))
		}

		// ensure cidr is a valid address
		if d.CIDR != "" {
			if _, _, err := net.ParseCIDR(d.CIDR); err != nil {
				result = multierror.Append(result, fmt.Errorf("[%s] %q: %w", "networking.os.device.CIDR", "", err))
			}
		}

		return result.ErrorOrNil()
	}
}

// CheckDeviceRoutes ensures that the specified routes are valid.
// nolint: dupl
func CheckDeviceRoutes() NetworkDeviceCheck {
	return func(d *machine.Device) error {
		var result *multierror.Error

		if len(d.Routes) == 0 {
			return result.ErrorOrNil()
		}

		for idx, route := range d.Routes {
			if _, _, err := net.ParseCIDR(route.Network); err != nil {
				result = multierror.Append(result, fmt.Errorf("[%s] %q: %w", "networking.os.device.route["+strconv.Itoa(idx)+"].Network", route.Network, ErrInvalidAddress))
			}

			if ip := net.ParseIP(route.Gateway); ip == nil {
				result = multierror.Append(result, fmt.Errorf("[%s] %q: %w", "networking.os.device.route["+strconv.Itoa(idx)+"].Gateway", route.Gateway, ErrInvalidAddress))
			}
		}
		return result.ErrorOrNil()
	}
}
