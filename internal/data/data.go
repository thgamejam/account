package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/gorm"

	"account/internal/conf"
	"github.com/thgamejam/pkg/cache"
	"github.com/thgamejam/pkg/database"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(
	NewData,
	NewAccountRepo,
	NewDataBase,
	NewCache,
)

func NewDataBase(c *conf.Data) (*gorm.DB, error) {
	return database.NewDataBase(c.Database)
}

func NewCache(c *conf.Data) (*cache.Cache, error) {
	return cache.NewCache(c.Redis)
}

// Data .
type Data struct {
	Cache    *cache.Cache
	DataBase *gorm.DB
}

// NewData .
func NewData(db *gorm.DB, cache *cache.Cache, logger log.Logger) (*Data, func(), error) {
	data := &Data{
		DataBase: db,
		Cache:    cache,
	}

	cleanup := func() {
		_ = cache.Close()
		log.NewHelper(logger).Info("closing the data resources")
	}
	return data, cleanup, nil
}
