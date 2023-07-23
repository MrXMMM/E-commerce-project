package fileshandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/MrXMMM/E-commerce-Project/config"
	"github.com/MrXMMM/E-commerce-Project/modules/entities"
	"github.com/MrXMMM/E-commerce-Project/modules/files"
	filesusecases "github.com/MrXMMM/E-commerce-Project/modules/files/filesUsecases"
	"github.com/MrXMMM/E-commerce-Project/pkg/util"
	"github.com/gofiber/fiber/v2"
)

type filesHandlerErrCode string

const (
	uploadErr filesHandlerErrCode = "file-001"
	deleteErr filesHandlerErrCode = "file-002"
)

type IFilesHandler interface {
	UploadFiles(c *fiber.Ctx) error
	DeleteFile(c *fiber.Ctx) error
}

type filesHandler struct {
	cfg         config.IConfig
	fileUsecase filesusecases.IFilesUsecase
}

func FilesHandler(cfg config.IConfig, fileUsecase filesusecases.IFilesUsecase) IFilesHandler {
	return &filesHandler{
		cfg:         cfg,
		fileUsecase: fileUsecase,
	}
}

func (h *filesHandler) UploadFiles(c *fiber.Ctx) error {
	req := make([]*files.FileReq, 0)

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}
	filesReq := form.File["files"]
	destination := c.FormValue("destination")

	// Files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}
	for _, file := range filesReq {
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				"extension is not acceptable",
			).Res()
		}

		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(uploadErr),
				fmt.Sprintf("file size must less than %d MiB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := util.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: destination + "/" + filename,
			FileName:    filename,
			Extension:   ext,
		})
	}

	res, err := h.fileUsecase.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(uploadErr),
			err.Error(),
		).Res()
	}

	// If you want to upload files to your computer please use this function below instead

	// res, err := h.filesUsecase.UploadToStorage(req)
	// if err != nil {
	// 	return entities.NewResponse(c).Error(
	// 		fiber.ErrInternalServerError.Code,
	// 		string(uploadErr),
	// 		err.Error(),
	// 	).Res()
	// }
	return entities.NewResponse(c).Success(fiber.StatusCreated, res).Res()
}

func (h *filesHandler) DeleteFile(c *fiber.Ctx) error {
	req := make([]*files.DeleteFileReq, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	if err := h.fileUsecase.DeleteFileOnGCP(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(deleteErr),
			err.Error(),
		).Res()
	}

	// If you want to upload files to your computer please use this function below instead

	// if err := h.filesUsecase.DeleteFileOnStorage(req); err != nil {
	// 	return entities.NewResponse(c).Error(
	// 		fiber.ErrInternalServerError.Code,
	// 		string(deleteErr),
	// 		err.Error(),
	// 	).Res()
	// }
	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}
