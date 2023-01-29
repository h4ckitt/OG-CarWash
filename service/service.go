package service

import (
	"car_wash/apperror"
	"car_wash/config"
	"car_wash/model"
	"car_wash/repository"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"
)

// TODO: TEMPORARY FIX, IMPLEMENT SOMETHING BETTER LATER
var (
	credsCache      *apiStore
	updatesCache    *updatesStore
	uzbekPlateRegex = regexp.MustCompile(`^(\d{2}[a-zA-Z]?\d{3}[a-zA-Z](?:[a-zA-Z]{1,2}|\*)?|[a-zA-Z]{3}\d{3})$`)
)

type CarWashSvc struct {
	Repo repository.Repo
}

func NewService(repository repository.Repo) *CarWashSvc {
	updatesCache = newUpdatesStore(18000) // ttl of 5 hours
	credsCache = newCredsJar(60)          // ttl of 1 minute
	return &CarWashSvc{repository}
}

func (c CarWashSvc) GetUpdatesChannel(ctx context.Context) <-chan model.Wash {
	uid := ctx.Value("ID").(string)
	log.Println(uid)
	return updatesCache.Get(uid)
}

func (c CarWashSvc) RegisterNewOwner(ctx context.Context, owner model.Owner) (string, error) {

	if ctx.Value("UID").(string) != owner.UUID {
		log.Println("UID's Don't Match")
		return "", &apperror.BadRequest
	}

	err := c.Repo.RegisterOwner(owner)

	if err != nil {
		log.Println(err)
		return "", err
	}

	key, err := c.CreateAPIKey(ctx)

	if err != nil {
		return "", err
	}

	return key, nil
}

func (c CarWashSvc) FetchDataByDate(ctx context.Context, date string) (model.WebSocketResult, error) {

	t, err := time.Parse("02.01.2006", date)

	if err != nil {
		return model.WebSocketResult{}, errors.New("invalid date specified")
	}

	res, err := c.Repo.FetchAllCarWashDataByDate(t.Format("2006-01-02"), ctx.Value("ID").(string))

	if err != nil {
		log.Println(err)
	}

	return res, err
}

func (c CarWashSvc) RegisterCarWash(ctx context.Context, carWash model.CarWash) (string, error) {
	id := ctx.Value("ID").(string)

	if carWash.Name == "" {
		return "", &apperror.BadRequest
	}

	id, err := c.Repo.RegisterCarWash(id, carWash)

	if err != nil {
		log.Println(err)
		return "", err
	}

	return id, nil
}

func (c CarWashSvc) CreateAPIKey(ctx context.Context) (string, error) {
	key, err := generateAPIKey()

	if err != nil {
		log.Println(err)
		return "", err
	}
	//var appError *apperror.AppError

	if err := c.Repo.StoreAPIKey(ctx.Value("UID").(string), key); err != nil {
		log.Println(err)
		return "", err
	}

	return key, nil
}

func (c CarWashSvc) SaveWashDetails(ctx context.Context, wash model.Wash) error {
	if !uzbekPlateRegex.MatchString(wash.NumberPlate) {
		log.Println("Unsupported plate format")
		return &apperror.UnprocessableEntity
	}

	id := ctx.Value("ID").(string)

	dt, err := time.Parse(time.RFC3339, wash.DateEntered)

	if err != nil {
		log.Println("Bad date format")
		return &apperror.BadRequest
	}

	wash.ImageName, err = c.SaveImage(ctx, wash)

	if err != nil {
		log.Println(err)
		return err
	}

	err = c.Repo.SaveWashDetails(id, wash)

	if err != nil {
		log.Println(err)
		return err
	}

	timeEntered := dt.Format("15:04:05")
	wash.TimeEntered = timeEntered

	wash.DateEntered = ""

	if ch, exists := updatesCache.Check(id); exists {
		ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		ch.Send(ctx, wash)
	}

	return nil
}

func (c CarWashSvc) SaveImage(ctx context.Context, wash model.Wash) (string, error) {
	dt, err := time.Parse(time.RFC3339, wash.DateEntered)

	if err != nil {
		retErr := &apperror.BadRequest
		retErr.Wrap(err)
		return "", retErr

	}
	imgConfig := config.GetConfig().ImageConfig
	filepath := fmt.Sprintf("%s/%s", imgConfig.Location, dt.Format("2006/01/02"))
	timeEntered := dt.Format("15:04:05")

	_ = os.MkdirAll(filepath, 0755)

	destFile, err := os.Create(fmt.Sprintf("%s/%s-%s.%s", filepath, imgConfig.Template, timeEntered, wash.ImageExt))

	if err != nil {
		retErr := &apperror.ServerError
		retErr.Wrap(err)
		return "", retErr
	}

	size, err := io.Copy(destFile, wash.Image)

	if size == 0 {
		retErr := &apperror.BadRequest
		retErr.Wrap(errors.New("empty image received"))
		return "", err
	}

	if err != nil {
		retErr := &apperror.ServerError
		retErr.Wrap(err)
		return "", retErr
	}

	return destFile.Name(), nil
}

func (c CarWashSvc) CacheCreds(ctx context.Context, hash string) error {
	key, err := c.Repo.RetrieveAPIKey(ctx.Value("UID").(string))

	if err != nil {
		log.Println(err)
		return err
	}

	credsCache.Insert(hash, key)

	return nil
}

func (c CarWashSvc) CheckCreds(hash string) (string, error) {
	key, err := credsCache.Get(hash)

	if err == nil {
		credsCache.Delete(hash)
	}

	return key, err
}

func (c CarWashSvc) ChangeAPIKey(ctx context.Context) error {
	key, err := generateAPIKey()

	if err != nil {
		log.Println(err)
		return &apperror.ServerError
	}

	var appError *apperror.AppError

	if err := c.Repo.UpdateAPIKey(ctx.Value("UID").(string), key); err != nil {
		if errors.As(err, &appError) {
			log.Println(appError.Error())
			return err
		}

		log.Println(err)

		return &apperror.ServerError
	}

	return nil
}
