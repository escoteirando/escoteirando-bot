package services

import (
	"fmt"
	"github.com/guionardo/escoteirando-bot/src/domain"
	"sync"
	"time"
)

type cachedAuthsType map[int64]domain.CachedAuthKey

var cachedAuths = struct {
	sync.RWMutex
	Auths cachedAuthsType
}{Auths: make(cachedAuthsType, 10)}

func GetAuthorizationFromCache(chatId int64) (string, error) {
	cachedAuths.RLock()
	defer cachedAuths.RUnlock()
	if auth, ok := cachedAuths.Auths[chatId]; ok {
		if auth.ValidUntil.Before(time.Now()) {
			delete(cachedAuths.Auths, chatId)
			return "", fmt.Errorf("invalidated authorization")
		}
		return auth.AuthKey, nil
	}
	return "", fmt.Errorf("authorization not found in cache")
}

func SetAuthorizationInCache(chat domain.Chat) error {
	cachedAuths.RLock()
	defer cachedAuths.RUnlock()
	if chat.ID == 0 {
		return fmt.Errorf("invalid chat id")
	}
	if len(chat.AuthKey) == 0 {
		return fmt.Errorf("invalid authorization key")
	}
	if chat.AuthValidUntil.Before(time.Now()) {
		return fmt.Errorf("authorization valid until %v", chat.AuthValidUntil)
	}
	cachedAuths.Auths[chat.ID] = domain.CachedAuthKey{
		ValidUntil: chat.AuthValidUntil,
		AuthKey:    chat.AuthKey,
	}
	return nil
}
