package handler

import (
	"context"
	"traffic-guarder/internal/infrastructure/app"
	"traffic-guarder/internal/infrastructure/errorsx"
	"traffic-guarder/internal/service"
	"traffic-guarder/internal/viewmodel"
)

type ExclusionHandler struct {
	as service.AnomalyService
}

func NewExclusionHandler(as service.AnomalyService) *ExclusionHandler {
	return &ExclusionHandler{as: as}
}

func (h *ExclusionHandler) GetExclusionList(c *app.Ctx) errorsx.APIError {
	var input viewmodel.ExclusionRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}

	resp, err := h.as.GetAnomalyEvents(context.Background(), &input)
	if err != nil {
		return errorsx.DatabaseError(err)
	}

	return c.SuccessResponse(resp, 1, "Exclusion list retrieved successfully")
}
