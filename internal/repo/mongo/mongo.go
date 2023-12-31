package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nightlord189/ca-url-shortener/internal/config"
	"github.com/nightlord189/ca-url-shortener/internal/entity"
	"github.com/nightlord189/ca-url-shortener/internal/usecase"
	"github.com/nightlord189/ca-url-shortener/pkg/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const usersCollection = "users"

type Repo struct {
	Config config.MongoConfig
	client *mongo.Client
}

func New(cfg config.MongoConfig) (*Repo, error) {
	//nolint: nosprintfhostport //it's ok for database url
	connectURI := fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.User, cfg.Password, cfg.Host, cfg.Port)

	cmdLogger := &event.CommandMonitor{
		Started: func(ctx context.Context, evt *event.CommandStartedEvent) {
			log.Ctx(ctx).Debugf("mongo request: %v", evt.Command)
		},
	}

	opts := options.Client().ApplyURI(connectURI).SetMonitor(cmdLogger).SetTimeout(10 * time.Second)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, fmt.Errorf("connect to mongodb error: %w", err)
	}

	return &Repo{
		Config: cfg,
		client: client,
	}, nil
}

func (r *Repo) CreateUser(ctx context.Context, user *entity.User) error {
	coll := r.client.Database(r.Config.Database).Collection(usersCollection)

	item := UserFromEntity(user)

	_, err := coll.InsertOne(ctx, item)
	return err
}

func (r *Repo) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	coll := r.client.Database(r.Config.Database).Collection(usersCollection)

	filter := bson.D{{"username", username}}
	var item User
	err := coll.FindOne(ctx, filter).Decode(&item)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return nil, usecase.ErrNotFound
	case err == nil:
		itemEntity := item.ToEntity()
		return &itemEntity, nil
	default:
		return nil, fmt.Errorf("find item error: %w", err)
	}
}

func (r *Repo) GetLink(ctx context.Context, shortURL string) (string, error) {
	coll := r.client.Database(r.Config.Database).Collection(usersCollection)

	filter := bson.D{{Key: fmt.Sprintf("links.%s", shortURL), Value: bson.D{{Key: "$exists", Value: "true"}}}}
	var item User
	err := coll.FindOne(ctx, filter).Decode(&item)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		return "", usecase.ErrNotFound
	case err == nil:
		return item.Links[shortURL], nil
	default:
		return "", err
	}
}

func (r *Repo) UpdateUserLinks(ctx context.Context, user *entity.User) error {
	coll := r.client.Database(r.Config.Database).Collection(usersCollection)

	links := bson.M{}
	for key, value := range user.Links {
		links[key] = value
	}

	update := bson.D{{"$set", bson.D{{"links", links}}}}

	_, err := coll.UpdateOne(ctx, bson.D{{"username", user.Username}}, update)
	return err
}
