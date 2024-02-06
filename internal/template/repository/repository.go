package repository

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ms-batch/internal/template/domain"
	baseModel "ms-batch/pkg/db"
	"ms-batch/pkg/errs"
)

type Repo struct {
	db   *mongo.Database
	base *baseModel.MongoDBClientRepository
}

func NewRepository(db *mongo.Database, base *baseModel.MongoDBClientRepository) Repository {
	return &Repo{db: db, base: base}
}

func (r Repo) CreateMailTemplate(ctx context.Context, model *domain.MailTemplate) error {
	_, err := r.db.Collection(model.CollectionName()).InsertOne(ctx, model)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (r Repo) FindMailTemplate(ctx context.Context, limit, page int64, column, name string) (
	*[]domain.MailTemplate, *baseModel.MongoPaginate, error,
) {
	var results []domain.MailTemplate
	var model domain.MailTemplate
	paginate := &baseModel.MongoPaginate{
		Limit: limit,
		Page:  page,
	}

	// Define an empty filter to match all documents.
	filter := bson.D{}
	if column != "" {
		regexPattern := primitive.Regex{Pattern: name, Options: "i"}
		filter = append(filter, bson.E{column, regexPattern})
	} else if column == "" && name != "" {
		// If column is empty, use regex for all columns.
		regexPattern := primitive.Regex{Pattern: name, Options: "i"}
		filter = bson.D{
			{
				"$or", bson.A{
					bson.D{{"name", regexPattern}},
					bson.D{{"subject", regexPattern}},
					bson.D{{"mail_layout.title", regexPattern}},
				},
			},
		}
	}

	count, err := r.db.Collection(model.CollectionName()).CountDocuments(ctx, filter)
	if err != nil {
		return nil, nil, err
	}
	paginate.InitiateTotal(count)
	opts := options.Find().SetProjection(bson.D{{"mail_layout.content", 0}, {"body", 0}})
	cursor, err := r.db.Collection(model.CollectionName()).Find(ctx, filter, opts, paginate.GetPaginatedOpts())
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var model domain.MailTemplate
		if err := cursor.Decode(&model); err != nil {
			return nil, nil, err
		}
		results = append(results, model)
	}

	if err := cursor.Err(); err != nil {
		return nil, nil, err
	}

	return &results, paginate, nil
}

func (r Repo) FindOneMailTemplate(ctx context.Context, id string) (*domain.MailTemplate, error) {
	var model domain.MailTemplate
	filter := bson.D{{"_id", id}}
	result := r.db.Collection(model.CollectionName()).FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If there are no documents found, return nil without an error.
			return nil, nil
		}
		return nil, err // Return other errors as is.
	}
	if err := result.Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (r Repo) FindOneMailTemplateByName(ctx context.Context, name string) (*domain.MailTemplate, error) {
	var model domain.MailTemplate
	filter := bson.D{{"name", name}}
	result := r.db.Collection(model.CollectionName()).FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// If there are no documents found, return nil without an error.
			return nil, nil
		}
		return nil, err // Return other errors as is.
	}
	if err := result.Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (r Repo) UpdateMailTemplate(ctx context.Context, id string, model *domain.MailTemplate) error {
	// Define a filter to match the document with the specified ID.
	filter := bson.M{"_id": id}

	// Convert the updates to a map (bson.M).
	update := bson.M{
		"$set": bson.M{
			"name":    model.Name,
			"subject": model.Subject,
			"body":    model.Body,
		},
	}

	_, err := r.db.Collection(model.CollectionName()).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r Repo) DeleteMailTemplate(ctx context.Context, id string) error {
	var model domain.MailTemplate
	// Define a filter to match the document with the specified ID.
	filter := bson.M{"_id": id}

	_, err := r.db.Collection(model.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
