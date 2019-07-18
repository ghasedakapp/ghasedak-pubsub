package pubsub

import "errors"

type ErrorsDef struct {
	SubscriptionNotFound error
	TopicNotFound        error
}

var Errors = ErrorsDef{
	SubscriptionNotFound: errors.New("subscription not found"),
	TopicNotFound:        errors.New("topic not found"),
}
