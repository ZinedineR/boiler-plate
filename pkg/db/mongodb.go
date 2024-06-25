package db

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"math"
	"os"
)

type MongoDBClientRepository struct {
	Client *mongo.Client
	DB     *mongo.Database
}

func NewMongoDBRepository(host, uname, pass string, port int) (
	*MongoDBClientRepository, error,
) {
	//var dsn string
	//dsn = fmt.Sprintf(os.Getenv("MONGODB_URL"))
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGODB_URL"))
	client, err := mongo.Connect(context.Background(), clientOptions)
	Users := client.Database(os.Getenv("MONGODB_DB"))
	if err != nil {
		logrus.Error(fmt.Sprintf("Cannot connect to MongoDB. %v", err))
		return nil, errors.Wrap(err, "Cannot connect to MongoDB")
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		logrus.Error(fmt.Sprintf("Cannot send/receive data with MongoDB. %v", err))
		return nil, errors.Wrap(err, "Cannot send/receive data with MongoDB")
	}
	return &MongoDBClientRepository{Client: client, DB: Users}, nil
}

type MongoPaginate struct {
	Limit      int64 `json:"limit,omitempty"`
	Page       int64 `json:"page,omitempty"`
	TotalRows  int64 `json:"total_rows,omitempty"`
	TotalPages int64 `json:"total_pages,omitempty"`
}

func (mp *MongoPaginate) GetPaginatedOpts() *options.FindOptions {
	l := mp.Limit
	skip := mp.Page*mp.Limit - mp.Limit
	fOpt := options.FindOptions{Limit: &l, Skip: &skip}

	return &fOpt
}

func (mp *MongoPaginate) InitiateTotal(count int64) {
	if mp.Limit == 0 {
		mp.Limit = count
	}
	if mp.Page == 0 {
		mp.Page = 1
	}
	mp.TotalRows = count
	mp.TotalPages = int64(math.Ceil(float64(count) / float64(mp.Limit)))
}
