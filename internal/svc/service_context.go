package svc

import (
	"github.com/XXueTu/wise/internal/config"
	"github.com/XXueTu/wise/internal/model"
)

type ServiceContext struct {
	Config        config.Config
	ModelsModel   *model.ModelsModel
	ResourceModel *model.ResourceModel
	TagsModel     *model.TagsModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	db := model.InitDB()
	return &ServiceContext{
		Config:        c,
		ModelsModel:   model.NewModelsModel(db),
		ResourceModel: model.NewResourceModel(db),
		TagsModel:     model.NewTagsModel(db),
	}
}
