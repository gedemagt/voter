package domain

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type Poll struct {
	UUID     uuid.UUID `json:"uuid"`
	Admin    User      `json:"admin"`
	Name     string    `json:"name"`
	SubPolls []SubPoll `json:"sub_polls"`
}

type SubPoll struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Open        bool      `json:"open"`
	Options     []Option  `json:"options"`
	Poll        *Poll     `json:"-"`
}

type Option struct {
	UUID        uuid.UUID `json:"uuid"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Votes       []User    `json:"votes"`
	SubPoll     *SubPoll  `json:"-"`
}

func (poll *Poll) AddSubPool(title string, description string) *SubPoll {
	subPoll := SubPoll{
		UUID:        uuid.New(),
		Title:       title,
		Description: description,
		Open:        false,
		Poll:        poll,
	}

	poll.SubPolls = append(poll.SubPolls, subPoll)
	return &poll.SubPolls[len(poll.SubPolls)-1]
}

func (subPoll *SubPoll) AddOption(title string, description string) *Option {
	option := Option{
		UUID:        uuid.New(),
		Title:       title,
		Description: description,
		SubPoll:     subPoll,
	}

	subPoll.Options = append(subPoll.Options, option)
	return &subPoll.Options[len(subPoll.Options)-1]
}

type PollRepository interface {
	GetPolls() ([]Poll, error)
	GetPollByUUID(uuid uuid.UUID) (Poll, error)
	CreateOrUpdatePoll(poll Poll) error
	DeletePoll(uuid string) error
	DeleteSubPoll(uuid string) error
	DeleteOption(uuid string) error
}

func (poll *Poll) ToJson() (string, error) {
	fmt.Printf("%#v", *poll)
	result, err := json.MarshalIndent(poll, "", "    ")
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func FromJson(data []byte) (Poll, error) {

	result := Poll{}

	if err := json.Unmarshal(data, &result); err != nil {
		return result, err
	}

	return result, nil
}
