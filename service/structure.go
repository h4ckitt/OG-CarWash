package service

import (
	"car_wash/apperror"
	"car_wash/model"
	"context"
	"time"
)

type update struct {
	userChan chan model.Wash
	lastUsed int64
}

type updatesStore struct {
	updatesChan map[string]*update
	locker      chan struct{}
}

type apiStore struct {
	credsJar map[string]creds
	locker   chan struct{}
}

type creds struct {
	apiKey     string
	insertTime int64
}

func newUpdatesStore(ttl int) (u *updatesStore) {
	u = &updatesStore{updatesChan: make(map[string]*update), locker: make(chan struct{}, 1)}

	go func() {
		for now := range time.Tick(time.Hour) {
			u.Lock()
			for key, cred := range u.updatesChan {
				if now.Unix()-cred.lastUsed > int64(ttl) {
					delete(u.updatesChan, key)
				}
			}
			u.Unlock()
		}
	}()

	return
}

func newCredsJar(ttl int) (a *apiStore) {
	a = &apiStore{credsJar: make(map[string]creds), locker: make(chan struct{}, 1)}

	go func() {
		for now := range time.Tick(time.Second * 5) {
			a.Lock()
			for key, cred := range a.credsJar {
				if now.Unix()-cred.insertTime > int64(ttl) {
					delete(a.credsJar, key)
				}
			}
			a.Unlock()
		}
	}()

	return
}

func (a *apiStore) Lock() {
	a.locker <- struct{}{}
}

func (a *apiStore) Unlock() {
	<-a.locker
}

func (a *apiStore) Insert(key, value string) {
	defer a.Unlock()

	a.Lock()

	a.credsJar[key] = creds{value, time.Now().Unix()}
}

func (a *apiStore) Get(key string) (string, error) {
	defer a.Unlock()

	a.Lock()

	if cred, exists := a.credsJar[key]; exists {
		return cred.apiKey, nil
	}

	return "", &apperror.NotFound
}

func (a *apiStore) Delete(key string) {
	defer a.Unlock()

	a.Lock()

	delete(a.credsJar, key)
}

func (c *updatesStore) Lock() {
	c.locker <- struct{}{}
}

func (c *updatesStore) Unlock() {
	<-c.locker
}

func (c *updatesStore) Get(uid string) <-chan model.Wash {
	defer c.Unlock()

	c.Lock()

	ch, exists := c.updatesChan[uid]

	if !exists {
		ch = &update{
			userChan: make(chan model.Wash),
			lastUsed: time.Now().Unix(),
		}
		c.updatesChan[uid] = ch
	}

	return ch.userChan
}

func (c *updatesStore) Check(uid string) (ch *update, exists bool) {
	defer c.Unlock()

	c.Lock()

	ch, exists = c.updatesChan[uid]

	return
}

func (up *update) Send(ctx context.Context, wash model.Wash) {
	select {
	case <-ctx.Done():
		return
	case up.userChan <- wash:
		up.lastUsed = time.Now().Unix()
		return
	}
}
