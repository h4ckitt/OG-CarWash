package postgres

import (
	"car_wash/apperror"
	"car_wash/config"
	"car_wash/model"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Repo struct {
	conn *sql.DB
}

const (
	RETRIEVECARWASHESBYOWNERUID    = "SELECT id FROM car_washes WHERE owner_id = $1"
	RETRIEVECARWASHOWNERUID        = "SELECT owner_id FROM car_washes WHERE id = $1"
	INSERTWASHDETAILS              = "INSERT INTO cars_washed (car_wash_id, license,img_file_path, date) VALUES ($1, $2, $3, $4)"
	RETRIEVEUSERWITHAPIKEY         = "SELECT user_id FROM api_keys WHERE key = $1"
	RETRIEVEAPIKEYWITHUID          = "SELECT key FROM api_keys WHERE user_id = $1"
	CHECKUSEREXISTENCEBEFOREINSERT = "SELECT EXISTS(SELECT 1 FROM owners WHERE id = $1)"
	INSERTCARWASH                  = "INSERT INTO car_washes (owner_id, name, address) VALUES ($1, $2, $3) RETURNING ID"
	INSERTAPIKEY                   = "INSERT INTO api_keys (key, user_id) VALUES ($1, $2)"
	INSERTNEWUSER                  = "INSERT INTO owners (id, first_name, last_name, email, phone) VALUES ($1, $2, $3, $4, $5)"
	FETCHOWNERPHONE                = "SELECT phone FROM owners INNER JOIN car_washes ON owners.id = car_washes.owner_id WHERE owner_id = $1"
	FETCHCARSWASHESBYCARWASHID     = "SELECT cws.Name, cwd.license, cwd.date FROM cars_washed cwd INNER JOIN car_washes cws ON cwd.car_wash_id = cws.id WHERE cws.id = $1 AND DATE(cwd.date) = $2"
)

func (r Repo) RetrieveAPIKey(uid string) (string, error) {
	var key string

	row := r.conn.QueryRow(RETRIEVEAPIKEYWITHUID, uid)

	if row.Err() != nil {
		e := &apperror.ServerError
		e.Wrap(row.Err())
		return "", e
	}

	if err := row.Scan(&key); err != nil {
		e := &apperror.AppError{}
		if err == sql.ErrNoRows {
			e = &apperror.NotFound
		} else {
			e = &apperror.ServerError
			e.Wrap(row.Err())
		}
		return "", e
	}

	return key, nil
}

func (r Repo) StoreAPIKey(uid string, key string) error {
	_, err := r.conn.Exec(INSERTAPIKEY, key, uid)

	if err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return e
	}

	return nil
}

func (r Repo) UpdateAPIKey(uid string, key string) error {
	//TODO implement me
	panic("implement me")
}

func (r Repo) RegisterOwner(owner model.Owner) error {
	row := r.conn.QueryRow(CHECKUSEREXISTENCEBEFOREINSERT, owner.UUID)

	if row.Err() != nil {
		e := &apperror.ServerError
		e.Wrap(row.Err())
		return e
	}

	var exists bool

	if err := row.Scan(&exists); err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return e
	}

	log.Println(exists)

	if exists {
		return &apperror.Conflict
	}
	_, err := r.conn.Exec(INSERTNEWUSER, owner.UUID, owner.FirstName, owner.LastName, owner.Email, owner.PhoneNumber)

	if err != nil {
		retErr := &apperror.ServerError
		retErr.Wrap(err)
		return retErr
	}
	return nil
}

func (r Repo) FetchCarWashDataByDate(date string, carWashID string) (model.WebSocketResult, error) {
	var (
		result        model.WebSocketResult
		clientResult  model.WebSocketClientResult
		carWashResult []model.Wash
		carWashName   string
	)

	result.Date = time.Now().Format("2006.01.02")

	row := r.conn.QueryRow(FETCHOWNERPHONE, carWashID)

	if err := row.Err(); err != nil {
		var e *apperror.AppError
		if err == sql.ErrNoRows {
			e = &apperror.NotFound
			return model.WebSocketResult{}, e
		}

		e = &apperror.ServerError
		e.Wrap(err)
		return model.WebSocketResult{}, e
	}

	if err := row.Scan(&clientResult.ClientNumber); err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return model.WebSocketResult{}, e
	}

	rows, err := r.conn.Query(FETCHCARSWASHESBYCARWASHID, carWashID, date)

	if err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return model.WebSocketResult{}, e
	}

	for rows.Next() {
		var wash model.Wash
		if err := rows.Scan(&carWashName, &wash.NumberPlate, &wash.DateEntered); err != nil {
			if err == sql.ErrNoRows {
				return model.WebSocketResult{}, nil
			}
			e := &apperror.ServerError
			e.Wrap(err)
			return model.WebSocketResult{}, e
		}
		dt, _ := time.Parse(time.RFC3339, wash.DateEntered)
		wash.TimeEntered = dt.Format("15:04:05")
		wash.DateEntered = ""
		carWashResult = append(carWashResult, wash)
	}

	clientResult.Result = []model.CarWashes{
		{
			CarWashID:   carWashID,
			CarWashName: carWashName,
			CarsEntered: len(carWashResult),
			Cars:        carWashResult,
		},
	}
	result.Date = date
	result.Clients = clientResult

	return result, nil
}

func (r Repo) FetchAllCarWashDataByDate(date string, ownerID string) (model.WebSocketResult, error) {
	var (
		result       model.WebSocketResult
		clientResult model.WebSocketClientResult
		i            = make(map[string][]model.Wash)
	)
	row := r.conn.QueryRow("SELECT phone FROM owners WHERE id = $1", ownerID)

	if err := row.Err(); err != nil {
		var e *apperror.AppError
		if err == sql.ErrNoRows {
			e = &apperror.NotFound
			return model.WebSocketResult{}, e
		}

		e = &apperror.ServerError
		e.Wrap(err)
		return model.WebSocketResult{}, e
	}

	if err := row.Scan(&clientResult.ClientNumber); err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return model.WebSocketResult{}, e
	}

	carsWashIDs, err := r.conn.Query(RETRIEVECARWASHESBYOWNERUID, ownerID)

	if err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return model.WebSocketResult{}, e
	}

	for carsWashIDs.Next() {
		var (
			carWashID string
		)

		if err := carsWashIDs.Scan(&carWashID); err != nil {
			if err == sql.ErrNoRows {
				return model.WebSocketResult{}, nil
			}

			e := &apperror.ServerError
			e.Wrap(err)

			return model.WebSocketResult{}, e
		}

		washes, err := r.conn.Query(FETCHCARSWASHESBYCARWASHID, carWashID, date)

		if err != nil {
			e := &apperror.ServerError
			e.Wrap(err)
			return model.WebSocketResult{}, e
		}

		for washes.Next() {
			var wash model.Wash
			if err := washes.Scan(&wash.CarWashName, &wash.NumberPlate, &wash.DateEntered); err != nil {
				if err == sql.ErrNoRows {
					continue
				} else {
					e := &apperror.ServerError
					e.Wrap(err)
					return model.WebSocketResult{}, e
				}
			}

			dt, _ := time.Parse(time.RFC3339, wash.DateEntered)

			wash.TimeEntered = dt.Format("15:04:05")

			wash.DateEntered = ""

			i[carWashID] = append(i[carWashID], wash)
		}
	}

	for id, washes := range i {
		clientResult.Result = append(clientResult.Result, model.CarWashes{
			CarWashID:   id,
			CarWashName: washes[0].CarWashName,
			CarsEntered: len(washes),
			Cars:        washes,
		})
	}

	result.Date = date
	result.Clients = clientResult
	return result, nil
}

func (r Repo) VerifyAPIKey(key string) (string, error) {
	row := r.conn.QueryRow(RETRIEVEUSERWITHAPIKEY, key)

	if row.Err() != nil {
		e := &apperror.ServerError
		e.Wrap(row.Err())
		return "", e
	}

	var uid string
	if err := row.Scan(&uid); err != nil {
		var e *apperror.AppError
		if err == sql.ErrNoRows {
			e = &apperror.NotFound
			return "", e
		}

		return "", &apperror.ServerError
	}

	return uid, nil
}

func (r Repo) RegisterCarWash(userID string, wash model.CarWash) (string, error) {
	res, err := r.conn.Query(INSERTCARWASH, userID, wash.Name, wash.Address)

	if err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return "", e
	}

	id := ""
	res.Next()
	if err := res.Scan(&id); err != nil {
		e := &apperror.ServerError
		e.Wrap(err)
		return "", e
	}

	return id, nil
}

func (r Repo) SaveWashDetails(uid string, wash model.Wash) error {
	row := r.conn.QueryRow(RETRIEVECARWASHOWNERUID, wash.CarWashID)

	if row.Err() != nil {
		e := &apperror.ServerError
		e.Wrap(row.Err())
		return e
	}

	var carWashOwnerID string

	if err := row.Scan(&carWashOwnerID); err != nil {
		e := &apperror.AppError{}

		if err == sql.ErrNoRows {
			e = &apperror.NotFound
		} else {
			e = &apperror.ServerError
			e.Wrap(err)
		}
		return e
	}

	if carWashOwnerID != uid {
		return &apperror.Unauthorized
	}

	_, err := r.conn.Exec(INSERTWASHDETAILS, wash.CarWashID, wash.NumberPlate, wash.ImageName, wash.DateEntered)

	if err != nil {
		var e *apperror.AppError
		if strings.Contains(err.Error(), "img_uniq") {
			e = &apperror.UnprocessableEntity
			e.Wrap(errors.New("duplicate record with timestamp passed"))
		} else {
			e = &apperror.ServerError
			e.Wrap(err)
		}
		return e
	}

	return nil
}

func NewPostgresRepo() (*Repo, error) {
	pgConfig := config.GetConfig().PGConfig

	dsn := fmt.Sprintf("host=%s port=%s user='%s' password='%s' dbname=%s sslmode=disable", pgConfig.DatabaseHost, pgConfig.DatabasePort, pgConfig.DatabaseUser, pgConfig.DatabasePassword, pgConfig.DatabaseName)

	db, err := sql.Open("postgres", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Repo{
		conn: db,
	}, nil
}
