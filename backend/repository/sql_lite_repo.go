package repository

import (
	"errors"
	"time"

	"example.com/hello/domain"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Base struct {
	Uuid      uuid.UUID `json:"uuid" gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PollModel struct {
	Base
	Admin uuid.UUID `json:"admin"`
	Name  string    `json:"name"`
}

type SubPollModel struct {
	Base
	Title       string `json:"title"`
	Description string `json:"description"`
	Open        bool   `json:"open"`
	PollID      uuid.UUID
	Poll        PollModel
}

type OptionModel struct {
	Base
	Title       string `json:"title"`
	Description string `json:"description"`
	SubPollID   uuid.UUID
	SubPoll     SubPollModel
}

type VoteModel struct {
	Base
	User      string `json:"user"`
	OptionID  uuid.UUID
	Option    OptionModel
	SubPollID uuid.UUID
	SubPoll   SubPollModel
}

type PollInviteModel struct {
	Base
	Token  string `json:"token"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	PollID uuid.UUID
	Poll   PollModel
}

func (repo *SQLiteRepository) InitDB(path string) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&PollModel{})
	db.AutoMigrate(&SubPollModel{})
	db.AutoMigrate(&OptionModel{})
	db.AutoMigrate(&VoteModel{})
	db.AutoMigrate(&PollInviteModel{})

	repo.DB = db
}

type SQLiteRepository struct {
	DB *gorm.DB
}

func (repo *SQLiteRepository) GetUserByUUID(_uuid uuid.UUID) (domain.User, error) {
	return domain.User{
		UUID:  uuid.New(),
		Name:  "Admin",
		Email: "admin@email.com",
		Role:  domain.SuperAdmin,
	}, nil
}

func (repo *SQLiteRepository) GetPollByUUID(uuid uuid.UUID) (domain.Poll, error) {
	var pollModel PollModel

	err := repo.DB.First(&pollModel, "uuid = ?", uuid).Error
	if err != nil {
		return domain.Poll{}, err
	}

	user, err := repo.GetUserByUUID(pollModel.Admin)

	if err != nil {
		return domain.Poll{}, err
	}

	poll := domain.Poll{
		UUID:  pollModel.Uuid,
		Admin: user,
		Name:  pollModel.Name,
	}

	return poll, nil
}

func (repo *SQLiteRepository) BuildDomainOption(option *OptionModel, subPoll *domain.SubPoll) (domain.Option, error) {

	result := domain.Option{
		UUID:        option.Uuid,
		Title:       option.Title,
		Description: option.Description,
		SubPoll:     subPoll,
	}

	return result, nil
}

func (repo *SQLiteRepository) BuildDomainSubPoll(subPoll *SubPollModel, poll *domain.Poll) (domain.SubPoll, error) {

	result := domain.SubPoll{
		UUID:        subPoll.Uuid,
		Title:       subPoll.Title,
		Description: subPoll.Description,
		Open:        subPoll.Open,
		Poll:        poll,
	}

	var options []OptionModel

	if err := repo.DB.Find(&options).Error; err != nil {
		return result, err
	}

	for _, optionModel := range options {
		option, _ := repo.BuildDomainOption(&optionModel, &result)
		result.Options = append(result.Options, option)
	}

	return result, nil
}

func (repo *SQLiteRepository) BuildDomainPoll(poll *PollModel) (domain.Poll, error) {
	user, _ := repo.GetUserByUUID(poll.Admin)

	result := domain.Poll{
		UUID:  poll.Uuid,
		Admin: user,
		Name:  poll.Name,
	}

	var subPollModels []SubPollModel

	if err := repo.DB.Find(&subPollModels).Error; err != nil {
		return result, err
	}

	for _, subPollModel := range subPollModels {
		subPoll, _ := repo.BuildDomainSubPoll(&subPollModel, &result)
		result.SubPolls = append(result.SubPolls, subPoll)
	}

	return result, nil
}

func (repo *SQLiteRepository) GetPolls() ([]domain.Poll, error) {
	var pollModels []PollModel
	var result []domain.Poll

	if err := repo.DB.Find(&pollModels).Error; err != nil {
		return result, err
	}

	for _, pollModel := range pollModels {
		poll, _ := repo.BuildDomainPoll(&pollModel)
		result = append(result, poll)
	}

	return result, nil
}

func (repo *SQLiteRepository) CreateOrUpdateOption(option domain.Option, parent *SubPollModel) error {

	optionModel := OptionModel{

		Title:       option.Title,
		Description: option.Description,
		SubPoll:     *parent,
	}

	optionModel.Uuid = option.UUID

	exists := repo.DB.Where("uuid = ?", option.UUID).First(&optionModel)

	if errors.Is(exists.Error, gorm.ErrRecordNotFound) {
		if err := repo.DB.Create(&optionModel).Error; err != nil {
			return err
		}
	} else {
		if err := exists.Updates(optionModel).Error; err != nil {
			return err
		}
	}

	return nil
}

func (repo *SQLiteRepository) CreateOrUpdateSubPoll(subPoll *domain.SubPoll, parent *PollModel) error {

	subPollModel := SubPollModel{
		Title:       subPoll.Title,
		Description: subPoll.Description,
		Open:        subPoll.Open,
		Poll:        *parent,
	}
	subPollModel.Uuid = subPoll.UUID

	exists := repo.DB.Where("uuid = ?", subPoll.UUID).First(&subPollModel)

	if errors.Is(exists.Error, gorm.ErrRecordNotFound) {
		if err := repo.DB.Create(&subPollModel).Error; err != nil {
			return err
		}
	} else {
		if err := exists.Updates(subPollModel).Error; err != nil {
			return err
		}
	}

	var relevantUUID []uuid.UUID
	for _, option := range subPoll.Options {
		relevantUUID = append(relevantUUID, option.UUID)
	}

	// First we delete existing ones that are not relevant any more
	if err := repo.DB.Where("sub_poll_id = ? and uuid NOT in ?", subPoll.UUID, relevantUUID).Delete(OptionModel{}).Error; err != nil {
		return err
	}

	// Create or update existing / new ones
	for _, option := range subPoll.Options {
		repo.CreateOrUpdateOption(option, &subPollModel)
	}

	return nil
}

func (repo *SQLiteRepository) CreateOrUpdatePoll(poll *domain.Poll) error {

	pollModel := PollModel{
		Admin: poll.Admin.UUID,
		Name:  poll.Name,
	}

	pollModel.Uuid = poll.UUID

	// We update or create it if it doesn't exists

	exists := repo.DB.Where("uuid = ?", poll.UUID).First(&pollModel)
	if errors.Is(exists.Error, gorm.ErrRecordNotFound) {
		if err := repo.DB.Create(&pollModel).Error; err != nil {
			return err
		}
	} else {
		if err := exists.Updates(pollModel).Error; err != nil {
			return err
		}
	}

	var relevantUUID []uuid.UUID
	for _, option := range poll.SubPolls {
		relevantUUID = append(relevantUUID, option.UUID)
	}

	// First we delete existing ones that are not relevant any more
	var subPollsToDelete []SubPollModel
	if err := repo.DB.Where("poll_id = ? and uuid NOT in ?", poll.UUID, relevantUUID).Find(&subPollsToDelete).Error; err != nil {
		return err
	}

	for _, subPoll := range subPollsToDelete {
		repo.DeleteSubPoll(subPoll.Uuid)
	}

	for _, subPoll := range poll.SubPolls {
		repo.CreateOrUpdateSubPoll(&subPoll, &pollModel)
	}

	return nil
}

func (repo *SQLiteRepository) DeletePoll(uuid uuid.UUID) error {

	var subPolls []SubPollModel

	if err := repo.DB.Where("poll_id = ?", uuid).Find(&subPolls).Error; err != nil {
		return err
	}

	for _, subPoll := range subPolls {
		repo.DeleteSubPoll(subPoll.Uuid)
	}

	repo.DB.Where("uuid = ?", uuid).Delete(&PollModel{})

	return nil
}

func (repo *SQLiteRepository) DeleteSubPoll(uuid uuid.UUID) error {

	if err := repo.DB.Where("sub_poll_id = ?", uuid).Delete(&OptionModel{}).Error; err != nil {
		return err
	}

	if err := repo.DB.Where("uuid = ?", uuid).Delete(&SubPollModel{}).Error; err != nil {
		return err
	}

	return nil

}
