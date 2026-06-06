package model

import (
	"fmt"

	sharedOrderModel "tmossDev.github.com/eco-system/shared-components/backend/package/order/model"
)

const (
	OrderStatusSubmitted      = sharedOrderModel.OrderStatusSubmitted
	OrderStatusConfirmed      = sharedOrderModel.OrderStatusConfirmed
	OrderStatusFulfillment    = sharedOrderModel.OrderStatusFulfillment
	OrderStatusOutForDelivery = sharedOrderModel.OrderStatusOutForDelivery
	OrderStatusDelivered      = sharedOrderModel.OrderStatusDelivered
	OrderStatusComplete       = sharedOrderModel.OrderStatusComplete
	OrderStatusReturned       = sharedOrderModel.OrderStatusReturned
	OrderStatusCancelled      = sharedOrderModel.OrderStatusCancelled
)

var orderStatusTransitions = map[string]map[string]struct{}{
	OrderStatusSubmitted: {
		OrderStatusConfirmed: {},
		OrderStatusCancelled: {},
	},
	OrderStatusConfirmed: {
		OrderStatusFulfillment: {},
		OrderStatusCancelled:   {},
	},
	OrderStatusFulfillment: {
		OrderStatusOutForDelivery: {},
		OrderStatusCancelled:      {},
	},
	OrderStatusOutForDelivery: {
		OrderStatusDelivered: {},
	},
	OrderStatusDelivered: {
		OrderStatusComplete: {},
		OrderStatusReturned: {},
	},
	OrderStatusComplete: {
		OrderStatusReturned: {},
	},
	OrderStatusReturned:  {},
	OrderStatusCancelled: {},
}

func IsKnownOrderStatus(status string) bool {
	_, ok := orderStatusTransitions[status]
	return ok
}

func CanTransitionOrderStatus(from string, to string) bool {
	allowedTransitions, ok := orderStatusTransitions[from]
	if !ok {
		return false
	}
	_, ok = allowedTransitions[to]
	return ok
}

func ValidateOrderStatusTransition(from string, to string) error {
	if !IsKnownOrderStatus(to) {
		return fmt.Errorf("unknown order status %q", to)
	}
	if from == to {
		return nil
	}
	if !CanTransitionOrderStatus(from, to) {
		return fmt.Errorf("cannot transition order from %q to %q", from, to)
	}
	return nil
}
