package events

import (
	services "github.com/crowdeco/bima/services"
)

type Model struct {
	Data       interface{}
	Id         string
	Repository *services.Repository
}
