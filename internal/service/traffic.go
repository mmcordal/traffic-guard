package service

import (
	"context"
	"errors"
	"fmt"
	"traffic-guarder/internal/repository"
	"traffic-guarder/internal/viewmodel"
)

type TrafficService interface {
	CreateLogAndGoBucket(ctx context.Context, vm *viewmodel.CreateTrafficLogRequest) error
}

type trafficService struct {
	tr repository.TrafficRepository
	bs BucketService
}

func NewTrafficService(tr repository.TrafficRepository, bs BucketService) TrafficService {
	return &trafficService{tr: tr, bs: bs}
}

func (s *trafficService) CreateLogAndGoBucket(ctx context.Context, vm *viewmodel.CreateTrafficLogRequest) error {
	if vm == nil {
		return errors.New("input is required")
	}

	log, err := vm.ToModel()
	if err != nil {
		return fmt.Errorf("vm to traffic log model error: %v", err)
	}

	// Logu db ye kaydet ve tabloda oluşturulanı dön
	createdLog, err := s.tr.CreateLog(ctx, log)
	if err != nil {
		return fmt.Errorf("create log error: %v", err)
	}

	// Gelen Logu ilgili bucket kolonunun değerlerini güncelle ve rediste de ilgili domain bucketini güncelle
	err = s.bs.UpsertBucket(ctx, createdLog)
	if err != nil {
		return fmt.Errorf("upsert bucket error: %v", err)
	}

	return nil
}
