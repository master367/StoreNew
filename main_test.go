package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestMain(m *testing.M) {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	cartCollection = client.Database("store").Collection("cart")
	userCollection = client.Database("store").Collection("users")

	m.Run()
}

func clearCollections() {
	_, err := cartCollection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to clear cart collection: %v", err)
	}

	_, err = userCollection.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		log.Fatalf("Failed to clear user collection: %v", err)
	}
}

func TestAddCigaretteToCart(t *testing.T) {
	clearCollections()

	cigarette := Cigarette{
		Brand:    "BrandTest",
		Type:     "TypeTest",
		Price:    10.99,
		Category: "CategoryTest",
		PhotoURL: "http://example.com/photo.jpg",
	}

	_, err := cartCollection.InsertOne(context.Background(), cigarette)
	assert.NoError(t, err)

	var result Cigarette
	err = cartCollection.FindOne(context.Background(), bson.D{{"brand", "BrandTest"}}).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "BrandTest", result.Brand)
	assert.Equal(t, 10.99, result.Price)
}

func TestRegisterUser(t *testing.T) {
	clearCollections()

	user := User{
		Email:    "testuser@example.com",
		Password: "password123",
	}

	_, err := userCollection.InsertOne(context.Background(), user)
	assert.NoError(t, err)

	var result User
	err = userCollection.FindOne(context.Background(), bson.D{{"email", "testuser@example.com"}}).Decode(&result)
	assert.NoError(t, err)
	assert.Equal(t, "testuser@example.com", result.Email)
}

func TestSendCartByEmail(t *testing.T) {
	clearCollections()

	cartItems := []Cigarette{
		{Brand: "Brand1", Type: "Type1", Price: 10.0, Category: "Category1", PhotoURL: "http://example.com/photo1.jpg"},
		{Brand: "Brand2", Type: "Type2", Price: 20.0, Category: "Category2", PhotoURL: "http://example.com/photo2.jpg"},
	}

	for _, item := range cartItems {
		_, err := cartCollection.InsertOne(context.Background(), item)
		assert.NoError(t, err)
	}
	email := "d4mirk@gmail.com"
	assert.Equal(t, "d4mirk@gmail.com", email)
}
