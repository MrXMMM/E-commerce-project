package appinfohandlers

import (
	"github.com/MrXMMM/E-commerce-Project/config"
	appinfousecases "github.com/MrXMMM/E-commerce-Project/modules/appinfo/appinfoUsecases"
)

type IAppinfoHandler interface {
}

type appinfoHandler struct {
	cfg             config.IConfig
	appinfoUsecases appinfousecases.IAppinfoUsecase
}

func AppinfoHandler(cfg config.IConfig, appinfoUsecases appinfousecases.IAppinfoUsecase) IAppinfoHandler {
	return &appinfoHandler{
		cfg:             cfg,
		appinfoUsecases: appinfoUsecases,
	}
}
