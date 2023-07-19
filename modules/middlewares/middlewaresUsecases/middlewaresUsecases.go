package middlewaresusecases

import middlewaresrepositories "github.com/MrXMMM/E-commerce-Project/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewaresRepository middlewaresrepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewaresRepository middlewaresrepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewaresRepository: middlewaresRepository,
	}
}
