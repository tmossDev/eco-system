package model

import "testing"

func TestValidateOrderStatusTransitionAllowsExpectedLifecycle(t *testing.T) {
	transitions := [][2]string{
		{OrderStatusSubmitted, OrderStatusConfirmed},
		{OrderStatusConfirmed, OrderStatusFulfillment},
		{OrderStatusFulfillment, OrderStatusOutForDelivery},
		{OrderStatusOutForDelivery, OrderStatusDelivered},
		{OrderStatusDelivered, OrderStatusComplete},
		{OrderStatusDelivered, OrderStatusReturned},
		{OrderStatusComplete, OrderStatusReturned},
		{OrderStatusSubmitted, OrderStatusCancelled},
		{OrderStatusConfirmed, OrderStatusCancelled},
		{OrderStatusFulfillment, OrderStatusCancelled},
	}

	for _, transition := range transitions {
		if err := ValidateOrderStatusTransition(transition[0], transition[1]); err != nil {
			t.Fatalf("expected transition %q -> %q to be valid: %v", transition[0], transition[1], err)
		}
	}
}

func TestValidateOrderStatusTransitionRejectsInvalidLifecycleJumps(t *testing.T) {
	transitions := [][2]string{
		{OrderStatusSubmitted, OrderStatusReturned},
		{OrderStatusSubmitted, OrderStatusDelivered},
		{OrderStatusDelivered, OrderStatusCancelled},
		{OrderStatusCancelled, OrderStatusConfirmed},
		{OrderStatusReturned, OrderStatusComplete},
	}

	for _, transition := range transitions {
		if err := ValidateOrderStatusTransition(transition[0], transition[1]); err == nil {
			t.Fatalf("expected transition %q -> %q to be invalid", transition[0], transition[1])
		}
	}
}
