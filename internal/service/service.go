package service

import (
	"account/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"

	pb "account/proto/api/account/v1"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewAccountService)

type AccountService struct {
	pb.UnimplementedAccountServer

	uc  *biz.AccountUseCase
	log *log.Helper
}

func NewAccountService(
	uc *biz.AccountUseCase,
	logger log.Logger,
) *AccountService {
	return &AccountService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}
