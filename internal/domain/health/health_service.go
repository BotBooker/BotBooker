package health

import "context"

type HealthService struct{}

func NewHealthService() *HealthService {
	return &HealthService{}
}

func (s *HealthService) Check(ctx context.Context) map[string]interface{} {
	return map[string]interface{}{
		"status":  "healthy",
		"service": "botbooker-api",
	}
}
