package appinfousecases

import appinforepositories "github.com/MrXMMM/E-commerce-Project/modules/appinfo/appinfoRepositories"

type IAppinfoUsecase interface {
}

type appinfoUsecase struct {
	appinfoRepository appinforepositories.IAppinfoRepository
}

func AppinfoUsecase(appinfoRepository appinforepositories.IAppinfoRepository) IAppinfoUsecase {
	return &appinfoUsecase{
		appinfoRepository: appinfoRepository,
	}
}
