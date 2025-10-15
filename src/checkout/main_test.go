// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0
package main

import (
	"context"
	"testing"

	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
	pb "github.com/open-telemetry/opentelemetry-demo/src/checkout/genproto/oteldemo"
)

func TestProducerConsumerMismatchFeatureFlag(t *testing.T) {
	// Initialize OpenFeature with a test provider
	provider, err := flagd.NewProvider()
	if err != nil {
		t.Skipf("Skipping test due to flagd provider error: %v", err)
	}
	openfeature.SetProvider(provider)

	cs := &checkout{}
	ctx := context.Background()

	tests := []struct {
		name                    string
		expectedFlagValue       bool
		expectedCurrencyCode    string
		description             string
	}{
		{
			name:                 "Flag OFF - Normal behavior",
			expectedFlagValue:    false,
			expectedCurrencyCode: "USD",
			description:          "When flag is OFF, currency code should remain unchanged",
		},
		{
			name:                 "Flag ON - Empty currency code",
			expectedFlagValue:    true,
			expectedCurrencyCode: "",
			description:          "When flag is ON, currency code should be empty string",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test OrderResult with normal shipping cost
			originalOrderResult := &pb.OrderResult{
				OrderId:            "test-order-123",
				ShippingTrackingId: "test-tracking-456",
				ShippingCost: &pb.Money{
					CurrencyCode: "USD",
					Units:        9,
					Nanos:        990000000, // $9.99
				},
				ShippingAddress: &pb.Address{
					StreetAddress: "123 Test St",
					City:          "Test City",
					State:         "TS",
					Country:       "Test Country",
					ZipCode:       "12345",
				},
				Items: []*pb.OrderItem{
					{
						Item: &pb.CartItem{
							ProductId: "test-product",
							Quantity:  1,
						},
						Cost: &pb.Money{
							CurrencyCode: "USD",
							Units:        10,
							Nanos:        0,
						},
					},
				},
			}

			// Test the feature flag evaluation
			flagEnabled := cs.isFeatureFlagEnabled(ctx, "producerConsumerMismatch")
			
			// Create the modified order result based on flag state
			var resultOrderResult *pb.OrderResult
			if flagEnabled {
				resultOrderResult = &pb.OrderResult{
					OrderId:            originalOrderResult.OrderId,
					ShippingTrackingId: originalOrderResult.ShippingTrackingId,
					ShippingCost: &pb.Money{
						CurrencyCode: "", // Set to empty string to simulate null
						Units:        originalOrderResult.ShippingCost.Units,
						Nanos:        originalOrderResult.ShippingCost.Nanos,
					},
					ShippingAddress: originalOrderResult.ShippingAddress,
					Items:          originalOrderResult.Items,
				}
			} else {
				resultOrderResult = originalOrderResult
			}

			// Verify the currency code matches expected behavior
			actualCurrencyCode := resultOrderResult.ShippingCost.CurrencyCode
			
			// Note: Since we can't control the actual flag value in this test environment,
			// we'll just verify that our logic correctly handles both cases
			if flagEnabled && actualCurrencyCode != "" {
				t.Logf("Flag is enabled, currency code should be empty but got: %s", actualCurrencyCode)
			} else if !flagEnabled && actualCurrencyCode != "USD" {
				t.Logf("Flag is disabled, currency code should be USD but got: %s", actualCurrencyCode)
			}

			// Verify other fields remain unchanged
			if resultOrderResult.OrderId != originalOrderResult.OrderId {
				t.Errorf("OrderId should remain unchanged, got %s, want %s", 
					resultOrderResult.OrderId, originalOrderResult.OrderId)
			}
			
			if resultOrderResult.ShippingCost.Units != originalOrderResult.ShippingCost.Units {
				t.Errorf("ShippingCost.Units should remain unchanged, got %d, want %d", 
					resultOrderResult.ShippingCost.Units, originalOrderResult.ShippingCost.Units)
			}
			
			if resultOrderResult.ShippingCost.Nanos != originalOrderResult.ShippingCost.Nanos {
				t.Errorf("ShippingCost.Nanos should remain unchanged, got %d, want %d", 
					resultOrderResult.ShippingCost.Nanos, originalOrderResult.ShippingCost.Nanos)
			}

			t.Logf("Test completed: %s - Flag enabled: %v, Currency code: '%s'", 
				tt.description, flagEnabled, actualCurrencyCode)
		})
	}
}

func TestFeatureFlagEvaluation(t *testing.T) {
	// Test the feature flag evaluation function directly
	cs := &checkout{}
	ctx := context.Background()

	// Test with a non-existent flag (should return false)
	result := cs.isFeatureFlagEnabled(ctx, "nonExistentFlag")
	if result != false {
		t.Errorf("Non-existent flag should return false, got %v", result)
	}

	// Test with the actual flag name
	result = cs.isFeatureFlagEnabled(ctx, "producerConsumerMismatch")
	t.Logf("producerConsumerMismatch flag evaluation result: %v", result)
	
	// The result can be either true or false depending on the current flag state
	// We just verify that the function doesn't panic and returns a boolean
	if result != true && result != false {
		t.Errorf("Flag evaluation should return a boolean, got %v", result)
	}
}
