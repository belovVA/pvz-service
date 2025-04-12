package service

type Repository interface {
	UserRepository
	PvzRepository
	ReceptionRepository
	ProductRepository
}

type Service struct {
	*AuthService
	*PvzService
	*ReceptionService
	*ProductService
	*InfoService
}

func NewService(repo Repository, jwtSecret string) *Service {
	return &Service{
		AuthService:      NewAuthService(repo, jwtSecret),
		PvzService:       NewPvzService(repo),
		ReceptionService: NewReceptionService(repo),
		ProductService:   NewProductService(repo, repo),
		InfoService:      NewInfoService(repo, repo, repo),
	}
}
