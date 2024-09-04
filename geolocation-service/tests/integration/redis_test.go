package integration_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"progekt/dating-app/geolocation-service/internal/models"
	"progekt/dating-app/geolocation-service/internal/modules/geoservice/storage"

	"github.com/go-redis/redis"
	"github.com/stretchr/testify/suite"
)

type GeoCacheTestSuite struct {
	suite.Suite
	rdb   *redis.Client
	cache *storage.GeoCache
}

func (suite *GeoCacheTestSuite) SetupSuite() {
	suite.rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6378",
	})
	suite.cache = storage.NewGeoCache(suite.rdb)
}

func (suite *GeoCacheTestSuite) TearDownSuite() {
	suite.rdb.Close()
}

func (suite *GeoCacheTestSuite) SetupTest() {
	suite.rdb.FlushDB()
}

func (suite *GeoCacheTestSuite) TestUpdateLocation() {
	user := models.User{
		ID:        1,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user, 10*time.Minute)
	suite.Require().NoError(err)

	// Проверяем, что данные пользователя были сохранены
	data, err := suite.rdb.Get("user_location:1").Bytes()
	suite.Require().NoError(err)

	var storedUser models.User
	err = json.Unmarshal(data, &storedUser)
	suite.Require().NoError(err)
	suite.Equal(user.ID, storedUser.ID)
	suite.Equal(user.Address, storedUser.Address)
	suite.Equal(user.Privacy, storedUser.Privacy)
	suite.Equal(user.H3Index, storedUser.H3Index)
}

func (suite *GeoCacheTestSuite) TestDeleteUser() {
	user := models.User{
		ID:        1,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user, 10*time.Minute)
	suite.Require().NoError(err)

	err = suite.cache.DeleteUser(user.ID)
	suite.Require().NoError(err)

	// Проверяем, что данные пользователя были удалены
	data, err := suite.rdb.Get("user_location:1").Bytes()
	suite.Require().Equal(redis.Nil, err)
	suite.Require().Empty(data)
}

func (suite *GeoCacheTestSuite) TestGetAllUsersInH3Cell() {
	user1 := models.User{
		ID:        1,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := models.User{
		ID:        2,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user1, 10*time.Minute)
	suite.Require().NoError(err)
	err = suite.cache.UpdateLocation(user2, 10*time.Minute)
	suite.Require().NoError(err)

	userIDs, err := suite.cache.GetAllUsersInH3Cell("85283473fffffff")
	suite.Require().NoError(err)
	suite.Len(userIDs, 2)
	suite.Contains(userIDs, user1.ID)
	suite.Contains(userIDs, user2.ID)
}

func (suite *GeoCacheTestSuite) TestGetUserLocation() {
	user := models.User{
		ID:        1,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user, 10*time.Minute)
	suite.Require().NoError(err)

	storedUser, err := suite.cache.GetUserLocation(user.ID)
	suite.Require().NoError(err)
	suite.Equal(user.ID, storedUser.ID)
	suite.Equal(user.Address, storedUser.Address)
	suite.Equal(user.Privacy, storedUser.Privacy)
	suite.Equal(user.H3Index, storedUser.H3Index)
}

func (suite *GeoCacheTestSuite) TestSetPrivacy() {
	user := models.User{
		ID:        1,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user, 10*time.Minute)
	suite.Require().NoError(err)

	err = suite.cache.SetPrivacy(context.Background(), user.ID, "private", 10*time.Minute)
	suite.Require().NoError(err)

	storedUser, err := suite.cache.GetUserLocation(user.ID)
	suite.Require().NoError(err)
	suite.Equal("private", storedUser.Privacy)
}

func (suite *GeoCacheTestSuite) TestGetUsersWithPrivacy() {
	user1 := models.User{
		ID:        1,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := models.User{
		ID:        2,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "private",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user1, 10*time.Minute)
	suite.Require().NoError(err)
	err = suite.cache.UpdateLocation(user2, 10*time.Minute)
	suite.Require().NoError(err)

	userIDs, err := suite.cache.GetUsersWithPrivacy("public")
	suite.Require().NoError(err)
	suite.Len(userIDs, 1)
	suite.Contains(userIDs, user1.ID)

	userIDs, err = suite.cache.GetUsersWithPrivacy("private")
	suite.Require().NoError(err)
	suite.Len(userIDs, 1)
	suite.Contains(userIDs, user2.ID)
}

func (suite *GeoCacheTestSuite) TestUpdateLocationWithInvalidData() {

	user := models.User{
		ID:        0,
		Address:   models.Address{Lat: 0, Lng: 0},
		Privacy:   "public",
		H3Index:   "invalid_h3_index",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user, 10*time.Minute)
	suite.Require().Error(err)
}

func (suite *GeoCacheTestSuite) TestDatabaseUnavailable() {

	suite.rdb.Close()

	user := models.User{
		ID:        1,
		Address:   models.Address{Lat: 37.7749, Lng: -122.4194},
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.cache.UpdateLocation(user, 10*time.Minute)
	suite.Require().Error(err)

	suite.rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6378",
	})
	suite.cache = storage.NewGeoCache(suite.rdb)
}

func TestGeoCacheTestSuite(t *testing.T) {
	suite.Run(t, new(GeoCacheTestSuite))
}
