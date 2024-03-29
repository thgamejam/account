package data

import (
	"account/internal/biz"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"math/rand"
	"strconv"

	"github.com/thgamejam/pkg/crypto/ecc"
)

const (
	lockOpenerListMaxLen = 5
)

var lockOpenerCacheKey = func(hash string) string {
	return "lock_opener_" + hash
}

var lockOpenerIDCacheKey = func(id int) string {
	return "lock_opener_id_to_key_" + strconv.Itoa(id)
}

// GetPublicKey 使用Hash值获取公钥
func (r *accountRepo) GetPublicKey(ctx context.Context, hash string) (*biz.PublicKey, error) {
	var lock LockOpener
	ok, err := r.data.Cache.Get(ctx, lockOpenerCacheKey(hash), &lock)
	if err != nil {
		r.log.Error("") // TODO
	}

	if ok {
		return &biz.PublicKey{
			Hash: hash,
			Key:  lock.Public,
		}, nil
	}

	return nil, errors.New("") // TODO
}

// GetRandomlyPublicKey 获取任意的公钥
func (r *accountRepo) GetRandomlyPublicKey(ctx context.Context) (key *biz.PublicKey, err error) {
	// 随机获取id
	id := rand.Intn(lockOpenerListMaxLen)

	// 查找缓存中是否已经存在id对于的密钥对
	var lock LockOpener
	hash, ok, _ := r.data.Cache.GetString(ctx, lockOpenerIDCacheKey(id))
	if ok {
		return r.GetPublicKey(ctx, hash)
	}

	// 缓存不存在密钥对时进行创建
	hash, err = r.CacheCreateLockOpener(ctx, &lock, id)
	if err != nil {
		return nil, err
	}

	return &biz.PublicKey{
		Hash: hash,
		Key:  lock.Public,
	}, nil
}

// GetPrivateKey 使用Hash值获取密钥
func (r *accountRepo) GetPrivateKey(ctx context.Context, hash string) (*biz.PrivateKey, error) {
	var lock LockOpener
	ok, err := r.data.Cache.Get(ctx, lockOpenerCacheKey(hash), &lock)
	if err != nil {
		r.log.Error("") // TODO
	}

	if ok {
		return &biz.PrivateKey{
			Hash: hash,
			Key:  lock.Private,
		}, nil
	}

	return nil, errors.New("") // TODO
}

// hashMd5To16 获取密钥md5 hash值，返回16个字符
var hashMd5To16 = func(privateKey string) string {
	bytes := md5.Sum([]byte(privateKey))
	return hex.EncodeToString(bytes[4:12])
}

// CacheCreateLockOpener 创建钥匙对到缓存中
func (r *accountRepo) CacheCreateLockOpener(ctx context.Context, lock *LockOpener, id int) (hash string, err error) {
	// 生成钥匙对
	privateKey, publicKey, err := ecc.GenerateKey()
	if err != nil {
		return
	}
	// 取密钥hash
	hash = hashMd5To16(privateKey)

	*lock = LockOpener{
		ID:      id,
		Public:  publicKey,
		Private: privateKey,
	}

	// 存缓存
	err = r.data.Cache.Set(ctx, lockOpenerCacheKey(hash), lock, 0)
	if err != nil {
		return
	}
	err = r.data.Cache.SetString(ctx, lockOpenerIDCacheKey(id), hash, 0)
	if err != nil {
		return
	}
	return
}
