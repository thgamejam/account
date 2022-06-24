package data

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"account/internal/conf"
	"github.com/thgamejam/pkg/cache"
	pkgConf "github.com/thgamejam/pkg/conf"
	"github.com/thgamejam/pkg/database"
)

var (
	Conf = &conf.Data{
		Database: &pkgConf.Database{
			Source:          "root:123456@tcp(127.0.0.1:3306)/touhou_gamejam?charset=utf8mb4&parseTime=True&loc=Local",
			MaxIdleConn:     0,
			MaxOpenConn:     0,
			ConnMaxLifetime: nil,
		},
		Redis: &pkgConf.Redis{
			Network:      "tcp",
			Addr:         "127.0.0.1:6379",
			Password:     "",
			ReadTimeout:  nil,
			WriteTimeout: nil,
		},
	}
)

func TestAccountRepo_GetAccountByID(t *testing.T) {
	Cache, _ := cache.NewCache(Conf.Redis)
	DataBase, _ := database.NewDataBase(Conf.Database)
	data, _, _ := NewData(DataBase, Cache, nil)
	assert.NotNil(t, data)
	ctx, _ := context.WithCancel(context.Background())

	_ = DataBase.AutoMigrate(&Account{})
	repo := NewAccountRepo(data, nil)
	user, err := repo.GetAccountByID(ctx, 1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	u, _ := json.Marshal(user)
	t.Logf("%v\n", string(u))
}
