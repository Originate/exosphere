package clientRegistry

import (
	"github.com/Originate/exocom/go/exocom/src/message_translator"
)

// Subscriber is the combination of a client name and internal namesapce
type Subscriber struct {
	ClientName        string
	InternalNamespace string
}

// SubscriptionManager manages what clients should be sent messages
type SubscriptionManager struct {
	MessageNameToSubscribers map[string][]Subscriber
	Routing                  Routes
}

// NewSubscriptionManager returns a SubscriptionManager for the given routes
func NewSubscriptionManager(routes Routes) *SubscriptionManager {
	result := new(SubscriptionManager)
	result.MessageNameToSubscribers = map[string][]Subscriber{}
	result.Routing = routes
	return result
}

// GetSubscribersFor returns all subscribers for the given message name
func (s *SubscriptionManager) GetSubscribersFor(messageName string) []Subscriber {
	subscribers, hasKey := s.MessageNameToSubscribers[messageName]
	if hasKey {
		return subscribers
	}
	return []Subscriber{}
}

// AddAll adds subscriptions for all messages the given client receives
// (as defined in the routes)
func (s *SubscriptionManager) AddAll(clientName string) {
	for _, messageName := range s.Routing[clientName].Receives {
		s.add(messageName, clientName)
	}
}

// RemoveAll removes subscriptions for all messages the given client receives
// (as defined in the routes)
func (s *SubscriptionManager) RemoveAll(clientName string) {
	for _, messageName := range s.Routing[clientName].Receives {
		s.remove(messageName, clientName)
	}
}

// Helpers

func (s *SubscriptionManager) add(internalMessageName, clientName string) {
	clientInternalNamespace := s.Routing[clientName].InternalNamespace
	publicMessageName := messageTranslator.GetPublicMessageName(&messageTranslator.GetPublicMessageNameOptions{
		InternalMessageName: internalMessageName,
		ClientName:          clientName,
		Namespace:           clientInternalNamespace,
	})
	if s.MessageNameToSubscribers[publicMessageName] == nil {
		s.MessageNameToSubscribers[publicMessageName] = []Subscriber{}
	}
	s.MessageNameToSubscribers[publicMessageName] = append(
		s.MessageNameToSubscribers[publicMessageName],
		Subscriber{ClientName: clientName, InternalNamespace: clientInternalNamespace})
}

func (s *SubscriptionManager) remove(messageName, clientName string) {
	clientInternalNamespace := s.Routing[clientName].InternalNamespace
	publicMessageName := messageTranslator.GetPublicMessageName(&messageTranslator.GetPublicMessageNameOptions{
		InternalMessageName: messageName,
		ClientName:          clientName,
		Namespace:           clientInternalNamespace,
	})
	if s.MessageNameToSubscribers[publicMessageName] == nil {
		return
	}
	for index, subscription := range s.MessageNameToSubscribers[publicMessageName] {
		if subscription.ClientName == clientName {
			s.MessageNameToSubscribers[publicMessageName] = append(s.MessageNameToSubscribers[publicMessageName][:index], s.MessageNameToSubscribers[publicMessageName][index+1:]...)
			return
		}
	}
}
