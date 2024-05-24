package main

import (
	"log"
	"time"
)

type User struct {
	Email string
}

type UserRepository interface {
	CreateUserAccount(u User) error
}

type NotificationsClient interface {
	SendNotification(u User) error
}

type NewsletterClient interface {
	AddToNewsletter(u User) error
}

type Handler struct {
	repository          UserRepository
	newsletterClient    NewsletterClient
	notificationsClient NotificationsClient
}

func NewHandler(
	repository UserRepository,
	newsletterClient NewsletterClient,
	notificationsClient NotificationsClient,
) Handler {
	return Handler{
		repository:          repository,
		newsletterClient:    newsletterClient,
		notificationsClient: notificationsClient,
	}
}

func (h Handler) SignUp(u User) error {
	if err := h.repository.CreateUserAccount(u); err != nil {
		return err
	}

	go func() {
		for {
			err := h.newsletterClient.AddToNewsletter(u)
			if err != nil {
				log.Printf("error in add to news letter failed to add user to the newsletter: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
	}()

	go func() {
		for {
			err := h.notificationsClient.SendNotification(u)
			if err != nil {
				log.Printf("error in send notification failed to add user to the newsletter: %v", err)
				time.Sleep(1 * time.Second)
				continue
			}
			break
		}
	}()

	return nil
}
