package service

type AnalyzeService interface {
}

type analyzeService struct {
}

func NewAnalyzeService() AnalyzeService {
	return &analyzeService{}
}
