package handler

import (
	"context"
	"traffic-guarder/internal/infrastructure/app"
	"traffic-guarder/internal/infrastructure/errorsx"
	"traffic-guarder/internal/service"
	"traffic-guarder/internal/viewmodel"
)

type TrafficHandler struct {
	ts service.TrafficService
}

func NewTrafficHandler(ts service.TrafficService) *TrafficHandler {
	return &TrafficHandler{ts: ts}
}

func (h *TrafficHandler) CreateTrafficLog(c *app.Ctx) errorsx.APIError {
	var input viewmodel.CreateTrafficLogRequest
	if errs := c.BodyParseValidate(&input); len(errs) > 0 {
		return errorsx.ValidationError(errs)
	}
	err := h.ts.CreateLogAndGoBucket(context.Background(), &input)
	if err != nil {
		return errorsx.DatabaseError(err)
	}
	return c.SuccessResponse("", 0, "Traffic log created successfully!")
}
