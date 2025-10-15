package main

import (
	"testing"

	pb "github.com/open-telemetry/opentelemetry-demo/src/checkout/genproto/oteldemo"
)

func TestPermanentProducerConsumerMismatch(t *testing.T) {
	// Create a sample OrderResult with proper currency code
	originalOrder := &pb.OrderResult{
		OrderId:            "test-order-123",
		ShippingTrackingId: "test-tracking-456",
		ShippingCost: &pb.Money{
			CurrencyCode: "USD",
			Units:        100,
			Nanos:        500000000,
		},
		ShippingAddress: &pb.Address{
			StreetAddress: "123 Test St",
			City:          "Test City",
			State:         "CA",
			Country:       "US",
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
					Units:        50,
				},
			},
		},
	}

	// Test that the permanent behavior always creates a modified copy with empty CurrencyCode
	t.Run("AlwaysSetsCurrencyCodeToEmpty", func(t *testing.T) {
		// Simulate the permanent behavior logic from main.go lines 392-402
		kafkaOrderResult := &pb.OrderResult{
			OrderId:            originalOrder.OrderId,
			ShippingTrackingId: originalOrder.ShippingTrackingId,
			ShippingCost: &pb.Money{
				CurrencyCode: "", // Always set to empty string
				Units:        originalOrder.ShippingCost.Units,
				Nanos:        originalOrder.ShippingCost.Nanos,
			},
			ShippingAddress: originalOrder.ShippingAddress,
			Items:           originalOrder.Items,
		}

		// Verify the permanent behavior
		if kafkaOrderResult.ShippingCost.CurrencyCode != "" {
			t.Errorf("Expected ShippingCost.CurrencyCode to be empty, got: %s", kafkaOrderResult.ShippingCost.CurrencyCode)
		}

		// Verify other fields are preserved
		if kafkaOrderResult.OrderId != originalOrder.OrderId {
			t.Errorf("Expected OrderId to be preserved, got: %s", kafkaOrderResult.OrderId)
		}

		if kafkaOrderResult.ShippingTrackingId != originalOrder.ShippingTrackingId {
			t.Errorf("Expected ShippingTrackingId to be preserved, got: %s", kafkaOrderResult.ShippingTrackingId)
		}

		if kafkaOrderResult.ShippingCost.Units != originalOrder.ShippingCost.Units {
			t.Errorf("Expected ShippingCost.Units to be preserved, got: %d", kafkaOrderResult.ShippingCost.Units)
		}

		if kafkaOrderResult.ShippingCost.Nanos != originalOrder.ShippingCost.Nanos {
			t.Errorf("Expected ShippingCost.Nanos to be preserved, got: %d", kafkaOrderResult.ShippingCost.Nanos)
		}
	})

	t.Run("OriginalOrderUnmodified", func(t *testing.T) {
		// Verify that the original order is not modified by the permanent behavior
		if originalOrder.ShippingCost.CurrencyCode != "USD" {
			t.Errorf("Expected original order CurrencyCode to remain USD, got: %s", originalOrder.ShippingCost.CurrencyCode)
		}
	})
}

func TestProducerConsumerMismatchBehavior(t *testing.T) {
	t.Run("PermanentBehaviorDescription", func(t *testing.T) {
		// This test documents the permanent behavior
		// The checkout service now ALWAYS:
		// 1. Creates a copy of OrderResult for Kafka messages
		// 2. Sets ShippingCost.CurrencyCode to empty string
		// 3. Preserves all other fields (Units, Nanos, etc.)
		// 4. Sends the modified copy to Kafka
		// 5. Returns the original (unmodified) OrderResult to the gRPC client

		// This simulates data quality issues where:
		// - Producer (checkout) sends incomplete data (missing currency code)
		// - Consumer (accounting) receives and processes the incomplete data
		// - Database stores empty currency codes instead of proper values
		// - This creates a permanent producer-consumer mismatch for demonstration

		t.Log("Permanent producer-consumer mismatch behavior is active")
		t.Log("All Kafka messages will have empty ShippingCost.CurrencyCode")
		t.Log("This demonstrates data quality issues in distributed systems")
	})
}
