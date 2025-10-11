package services

import "AskharovA/go-rest-api/repositories"

type EventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(repo repositories.EventRepository) *EventService {
	return &EventService{
		eventRepo: repo,
	}
}
