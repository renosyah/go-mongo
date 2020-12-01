package model

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	User struct {
		ID    primitive.ObjectID `json:"id" bson:"_id"`
		Name  string             `json:"name" bson:"name"`
		Email string             `json:"email" bson:"email"`
	}

	UserResponse struct {
		ID    primitive.ObjectID `json:"id"`
		Name  string             `json:"name"`
		Email string             `json:"email"`
	}

	AllUser struct {
		SearchBy    string `json:"search_by"`
		SearchValue string `json:"search_value"`
		OrderBy     string `json:"order_by"`
		OrderDir    string `json:"order_dir"`
		Offset      int    `json:"offset"`
		Limit       int    `json:"limit"`
	}
)

func (u *User) Response() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func (u *User) Add(ctx context.Context, db *mongo.Database) (primitive.ObjectID, error) {
	result, err := db.Collection("user").InsertOne(ctx,
		bson.M{
			"name":  u.Name,
			"email": u.Email,
		},
	)
	if err != nil {
		return u.ID, err
	}

	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		u.ID = id
	}

	return u.ID, nil
}

func (u *User) All(ctx context.Context, db *mongo.Database, param AllUser) ([]User, error) {
	all := []User{}

	offset, limit := int64(param.Offset), int64(param.Limit)
	search := bson.D{
		primitive.E{
			Key: param.SearchBy,
			Value: primitive.Regex{
				Pattern: param.SearchValue,
				Options: "",
			},
		},
	}
	opts := options.FindOptions{
		Skip:  &offset,
		Limit: &limit,
		Sort: bson.D{
			primitive.E{
				Key:   param.OrderBy,
				Value: param.OrderToInt(),
			},
		},
	}

	cur, err := db.Collection("user").Find(ctx, search, &opts)
	if err != nil {
		return all, err
	}

	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var one User
		err := cur.Decode(&one)
		if err != nil {
			return all, err
		}

		all = append(all, one)
	}

	if err := cur.Err(); err != nil {
		return all, err
	}

	return all, nil
}

func (u *User) One(ctx context.Context, db *mongo.Database) (User, error) {
	one := User{}
	result := db.Collection("user").FindOne(ctx,
		bson.D{
			primitive.E{
				Key:   "_id",
				Value: u.ID,
			},
		},
	)
	err := result.Decode(&one)
	if err != nil {
		return one, err
	}
	return one, nil
}

func (u *User) Update(ctx context.Context, db *mongo.Database) (primitive.ObjectID, error) {
	var id primitive.ObjectID

	result, err := db.Collection("user").UpdateOne(ctx,
		bson.M{
			"_id": u.ID,
		},
		bson.M{
			"$set": bson.M{
				"name":  u.Name,
				"email": u.Email,
			},
		},
	)
	if err != nil {
		return id, err
	}

	if num := result.ModifiedCount; num > 0 {
		id = u.ID
	}

	return id, nil
}

func (u *User) Delete(ctx context.Context, db *mongo.Database) (primitive.ObjectID, error) {
	var id primitive.ObjectID

	result, err := db.Collection("user").DeleteOne(ctx,
		bson.M{
			"_id": u.ID,
		},
	)
	if err != nil {
		return id, err
	}

	if num := result.DeletedCount; num > 0 {
		id = u.ID
	}

	return id, nil
}
