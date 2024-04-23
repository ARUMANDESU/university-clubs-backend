package notificationSrvc

import (
	"context"
	"fmt"
	notifv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/notification"
	"github.com/ARUMANDESU/university-clubs-backend/internal/clients/notification"
	"github.com/ARUMANDESU/university-clubs-backend/internal/domain"
	"github.com/ARUMANDESU/university-clubs-backend/pkg/logger"
	"io"
	"log/slog"
	"sync"
)

type Service struct {
	connections        *sync.Map
	notificationClient *notification.Client
	log                *slog.Logger
}

func NewService(logger *slog.Logger, client *notification.Client) *Service {
	return &Service{connections: new(sync.Map), log: logger, notificationClient: client}
}

func (s *Service) AddNewConnection(userId int64) chan domain.Notification {
	notificationChan := make(chan domain.Notification)
	s.connections.Store(userId, notificationChan)
	s.log.Debug("new connection added", slog.Int64("userID", userId))
	return notificationChan
}

func (s *Service) RemoveConnection(userId int64) {
	if notificationChan, ok := s.connections.Load(userId); ok {
		ch, ok := notificationChan.(chan domain.Notification)
		if ok {
			close(ch)
		}
	}

	s.connections.Delete(userId)
	s.log.Debug("connection removed", slog.Int64("userID", userId))
}

func (s *Service) GetConnection(userId int64) (<-chan domain.Notification, error) {
	const op = "Services.Notification.GetConnection"
	log := s.log.With(slog.String("op", op))

	if notificationChan, ok := s.connections.Load(userId); ok {
		ch, ok := notificationChan.(chan domain.Notification)
		if ok {
			return ch, nil
		} else {
			log.Error("notificationChan is not a channel of domain.Notification")
			return nil, fmt.Errorf("%s: notificationChan is not a channel of domain.Notification", ok)
		}
	}

	return nil, fmt.Errorf("%s connection not found", op)
}

func (s *Service) HandleNotification() {
	const op = "Services.Notification.HandleNotification"
	log := s.log.With(slog.String("op", op))
	stream, err := s.notificationClient.GetNotifications(context.Background(), &notifv1.NotificationRequest{Message: "lol"})
	if err != nil {
		log.Error("error getting notification stream", logger.Err(err))
	}
	log.Info("got notification stream")

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(fmt.Sprintf("%s: %v.GetNotifications(_) = _, %v", op, s.notificationClient, err))
		}
		log.Info(res.Message)

		// if user is online then send notification
		if notificationChan, ok := s.connections.Load(res.UserId); ok {
			ch, ok := notificationChan.(chan domain.Notification)
			if ok {
				// Now you can send data to notificationChan
				ch <- domain.NotificationProtoToDomain(res)
			} else {
				// Handle the case where notificationChan is not a channel of domain.Notification
				log.Error("notificationChan is not a channel of domain.Notification")
			}
		}
		// TODO: handle if user is offline then store notification in the storage
	}

}
