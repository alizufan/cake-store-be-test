package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zufzuf/cake-store/repository"
	"github.com/zufzuf/cake-store/schema"
	"github.com/zufzuf/cake-store/service"
)

var CakeRepository = &repository.CakeMock{Mock: mock.Mock{}}
var CakeService = &service.Cake{Repo: CakeRepository}

var cake = schema.Cake{
	ID:          1,
	Title:       "Test Title",
	Description: "Test Description",
	Rating:      7,
	Image:       "https://cdn.lorem.space/images/movie/.cache/150x220/totoro-1988.jpeg",
	CreatedAt:   time.Now(),
	UpdatedAt:   time.Now(),
}

var cakes = []schema.Cake{
	cake,
	{
		ID:          2,
		Title:       "Test Title 2",
		Description: "Test Description 2",
		Rating:      8,
		Image:       "https://cdn.lorem.space/images/movie/.cache/150x220/totoro-1988.jpeg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	},
}

func Test_Cake_Service_Find(t *testing.T) {
	tests := []struct {
		Name           string
		Request        int
		ExpectedResult *schema.Cake
		ExpectedError  error
	}{
		{
			Name:           "Found",
			Request:        1,
			ExpectedResult: &cake,
			ExpectedError:  nil,
		},
		{
			Name:           "Not_Found",
			Request:        0,
			ExpectedResult: nil,
			ExpectedError:  repository.ErrRecordNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			CakeRepository.Mock.On("Find", context.Background(), test.Request).Return(test.ExpectedResult, test.ExpectedError)

			res, err := CakeService.Find(context.Background(), test.Request)
			if test.Name == "Not_Found" {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedResult, res)
			}
		})
	}
}

func Test_Cake_Service_Find_All(t *testing.T) {
	tests := []struct {
		Name           string
		Filter         repository.FindAllFilter
		Request        service.FindAllRequest
		ExpectedResult []schema.Cake
		ExpectedError  error
	}{
		{
			Name:           "Not_Found",
			ExpectedResult: nil,
			ExpectedError:  repository.ErrRecordNotFound,
		},
		{
			Name:           "Request_Nil",
			ExpectedResult: nil,
			ExpectedError:  service.ErrRequestNil,
		},
		{
			Name:           "Found_Without_Filter",
			Filter:         repository.FindAllFilter{},
			Request:        service.FindAllRequest{},
			ExpectedResult: cakes,
			ExpectedError:  nil,
		},
		{
			Name: "Found_With_Title_Filter",
			Filter: repository.FindAllFilter{
				Title: "Test Title",
			},
			Request: service.FindAllRequest{
				Title: "Test Title",
			},
			ExpectedResult: cakes,
			ExpectedError:  nil,
		},
		{
			Name: "Found_With_Description_Filter",
			Filter: repository.FindAllFilter{
				Description: "Test Description",
			},
			Request: service.FindAllRequest{
				Description: "Test Description",
			},
			ExpectedResult: cakes,
			ExpectedError:  nil,
		},
		{
			Name: "Found_With_All_Filter",
			Filter: repository.FindAllFilter{
				Title:       "Test Title",
				Description: "Test Description",
			},
			Request: service.FindAllRequest{
				Title:       "Test Title",
				Description: "Test Description",
			},
			ExpectedError: nil,
		},
	}

	for _, v := range tests {
		test := v
		t.Run(test.Name, func(t *testing.T) {
			CakeRepository.Mock.On("FindAll", context.Background(), &test.Filter).Return(test.ExpectedResult, test.ExpectedError).Once()

			res, err := CakeService.FindAll(context.Background(), &test.Request)
			if test.Name == "Not_Found" || test.Name == "Request_Nil" {
				assert.Error(t, err)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedResult, res)
			}
		})
	}
}

func Test_Cake_Service_Insert(t *testing.T) {
	record := cake
	record.ID = 0
	record.CreatedAt = time.Time{}
	record.UpdatedAt = time.Time{}

	tests := []struct {
		Name          string
		Record        *schema.Cake
		Request       *service.CakeRequest
		ExpectedError error
	}{
		{
			Name:   "Insert_Success",
			Record: &record,
			Request: &service.CakeRequest{
				Title:       record.Title,
				Description: record.Description,
				Rating:      record.Rating,
				Image:       record.Image,
			},
			ExpectedError: nil,
		},
		{
			Name:          "Insert_Record_Nil",
			Record:        &record,
			Request:       nil,
			ExpectedError: service.ErrRequestNil,
		},
	}
	for _, v := range tests {
		test := v
		t.Run(test.Name, func(t *testing.T) {
			CakeRepository.Mock.On("Insert", context.Background(), test.Record).Return(test.ExpectedError).Once()
			err := CakeService.Insert(context.Background(), test.Request)
			if test.Name == "Insert_Record_Nil" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.Record.ID, test.Request.ID)
			}
		})
	}
}

func Test_Cake_Service_Update(t *testing.T) {
	record := cake
	record.CreatedAt = time.Time{}
	record.UpdatedAt = time.Time{}

	tests := []struct {
		Name              string
		Record            *schema.Cake
		Request           *service.CakeRequest
		ExpectedResult    *schema.Cake
		ExpectedError     error
		ExpectedFindError error
	}{
		{
			Name:           "Update_Success",
			Record:         &record,
			ExpectedResult: &record,
			Request: &service.CakeRequest{
				ID:          record.ID,
				Title:       record.Title,
				Description: record.Description,
				Rating:      record.Rating,
				Image:       record.Image,
			},
			ExpectedError:     nil,
			ExpectedFindError: nil,
		},
		{
			Name:              "Update_Record_Nil",
			Record:            &record,
			Request:           nil,
			ExpectedError:     service.ErrRequestNil,
			ExpectedFindError: repository.ErrRecordNotFound,
		},
	}
	for _, v := range tests {
		test := v
		t.Run(test.Name, func(t *testing.T) {
			isUpdateRecordNil := test.Name == "Update_Record_Nil"
			id := test.Record.ID
			if isUpdateRecordNil {
				id = 0
			}

			CakeRepository.Mock.On("Find", context.Background(), id).Return(test.ExpectedResult, test.ExpectedFindError).Once()
			CakeRepository.Mock.On("Update", context.Background(), test.Record).Return(test.ExpectedError).Once()

			err2 := CakeService.Update(context.Background(), test.Request)
			if isUpdateRecordNil {
				assert.Error(t, err2)
			} else {
				assert.NoError(t, err2)
			}
		})
	}
}

func Test_Cake_Service_Delete(t *testing.T) {
	record := cake
	record.CreatedAt = time.Time{}
	record.UpdatedAt = time.Time{}

	tests := []struct {
		Name              string
		Request           int
		ExpectedError     error
		ExpectedFindError error
	}{
		{
			Name:              "Delete_Success",
			Request:           1,
			ExpectedError:     nil,
			ExpectedFindError: nil,
		},
		{
			Name:              "Delete_Not_Found",
			Request:           0,
			ExpectedError:     nil,
			ExpectedFindError: repository.ErrRecordNotFound,
		},
	}
	for _, v := range tests {
		test := v
		t.Run(test.Name, func(t *testing.T) {
			notfound := test.Name == "Delete_Not_Found"
			CakeRepository.Mock.On("Find", context.Background(), test.Request).Return(nil, test.ExpectedFindError).Once()
			CakeRepository.Mock.On("Delete", context.Background(), test.Request).Return(test.ExpectedError).Once()

			err := CakeService.Delete(context.Background(), test.Request)
			if notfound {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
