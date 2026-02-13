// SPDX-FileCopyrightText: Copyright 2025 Carabiner Systems, Inc
// SPDX-License-Identifier: Apache-2.0

package sts

import (
	"github.com/carabiner-dev/signer-extras/sts/providers/spiffe"
	osts "github.com/carabiner-dev/signer/sts"
)

// RegisterExtraProviders registers all the providers from the
// extra package on the signer
func RegisterExtraProviders() {
	osts.RegisterProvider("spire", &spiffe.Spire{})
}
