package services

import (
	"AskharovA/go-rest-api/models"
	"AskharovA/go-rest-api/repositories"
)

type EventService struct {
	eventRepo repositories.EventRepository
}

func NewEventService(repo repositories.EventRepository) *EventService {
	return &EventService{
		eventRepo: repo,
	}
}

func (s *EventService) CreateEvent(event *models.Event) error {
	eventId, err := s.eventRepo.Save(event)
	if err != nil {
		return err
	}

	event.ID = eventId
	return nil
}

func (s *EventService) GetEvents(page, per_page int) ([]models.Event, error) {
	return s.eventRepo.GetAllEvents(page, per_page)
}

func (s *EventService) GetEvent(eventId int64) (*models.Event, error) {
	return s.eventRepo.GetEventByID(eventId)
}

func (s *EventService) UpdateEvent(event *models.Event) error {
	return s.eventRepo.Update(event)
}

func (s *EventService) DeleteEvent(event *models.Event) error {
	return s.eventRepo.Delete(event)
}

func (s *EventService) Register(userId int64, event *models.Event) error {
	return s.eventRepo.Register(userId, event)
}

func (s *EventService) CancelRegistration(userId int64, event *models.Event) error {
	return s.eventRepo.CancelRegistration(userId, event)
}
