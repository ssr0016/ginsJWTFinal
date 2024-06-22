package service

import (
	"gins/data/request"
	"gins/data/response"
	"gins/helper"
	"gins/model"
	"gins/repository"

	"github.com/go-playground/validator/v10"
)

type TagsServiceImpl struct {
	tagsRepository repository.TagsRepository
	validate       *validator.Validate
}

func NewTagsServiceImpl(tagsRepository repository.TagsRepository, validate *validator.Validate) TagsService {
	return &TagsServiceImpl{
		tagsRepository: tagsRepository,
		validate:       validate,
	}
}

func (t *TagsServiceImpl) Create(tags request.CreateTagsRequest) {
	err := t.validate.Struct(tags)
	helper.ErrorPanic(err)
	tagModel := model.Tags{
		Name: tags.Name,
	}

	t.tagsRepository.Save(tagModel)
}

func (t *TagsServiceImpl) Delete(tagsId int) {
	t.tagsRepository.Delete(tagsId)
}

func (t *TagsServiceImpl) FindAll() []response.TagsResponse {
	result := t.tagsRepository.FindAll()

	var tags []response.TagsResponse
	for _, value := range result {
		tag := response.TagsResponse{
			Id:   value.Id,
			Name: value.Name,
		}
		tags = append(tags, tag)
	}

	return tags
}

func (t *TagsServiceImpl) FindById(tagsId int) response.TagsResponse {
	tagData, err := t.tagsRepository.FindById(tagsId)
	helper.ErrorPanic(err)

	tagResponse := response.TagsResponse{
		Id:   tagData.Id,
		Name: tagData.Name,
	}

	return tagResponse
}

func (t *TagsServiceImpl) Update(tags request.UpdateTagsRequest) {
	tagData, err := t.tagsRepository.FindById(tags.Id)
	helper.ErrorPanic(err)
	tagData.Name = tags.Name

	t.tagsRepository.Update(tagData)
}
