package integration_test

import (
	"context"
	"testing"
	"time"

	"progekt/dating-app/geolocation-service/internal/infrastructure/db/postgres"
	"progekt/dating-app/geolocation-service/internal/models"
	"progekt/dating-app/geolocation-service/internal/modules/geoservice/storage"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
)

type UserStorageTestSuite struct {
	suite.Suite
	dbPool  *pgxpool.Pool
	storage *storage.UserStorage
}

func (suite *UserStorageTestSuite) SetupSuite() {
	var err error
	suite.dbPool, err = pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable")
	suite.Require().NoError(err)

	suite.storage = storage.NewUserLocationStorage(postgres.NewPostgresDBForTest(suite.dbPool))
}

func (suite *UserStorageTestSuite) TearDownSuite() {
	suite.dbPool.Close()
}

func (suite *UserStorageTestSuite) SetupTest() {
	_, err := suite.dbPool.Exec(context.Background(), "TRUNCATE TABLE users_table, location_history RESTART IDENTITY CASCADE")
	suite.Require().NoError(err)
}

func (suite *UserStorageTestSuite) TestAddOrUpdateUserLocation() {
	ctx := context.Background()
	user := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user)
	suite.Require().NoError(err)

	var count int
	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM users_table WHERE id=$1", user.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(1, count)

	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM location_history WHERE user_id=$1", user.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(1, count)
}

func (suite *UserStorageTestSuite) TestGetUser() {
	ctx := context.Background()
	user := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := suite.dbPool.Exec(ctx, `
		INSERT INTO users_table (id, privacy, h3_index, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Privacy, user.H3Index, user.CreatedAt, user.UpdatedAt)
	suite.Require().NoError(err)

	retrievedUser, err := suite.storage.GetUser(ctx, user.ID)
	suite.Require().NoError(err)
	suite.Equal(user.ID, retrievedUser.ID)
	suite.Equal(user.Privacy, retrievedUser.Privacy)
	suite.Equal(user.H3Index, retrievedUser.H3Index)
}

func (suite *UserStorageTestSuite) TestGetLocationHistory() {
	ctx := context.Background()
	user := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user)
	suite.Require().NoError(err)

	historyUser, err := suite.storage.GetLocationHistory(ctx, user.ID, 10)
	suite.Require().NoError(err)
	suite.Equal(user.ID, historyUser.ID)
	suite.Len(historyUser.LHistory, 1)
	suite.Equal(user.H3Index, historyUser.LHistory[0].H3Index)
}

func (suite *UserStorageTestSuite) TestShareLocation() {
	ctx := context.Background()
	user1 := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := &models.User{
		ID:        2,
		Privacy:   "public",
		H3Index:   "85283473ffffffe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user1)
	suite.Require().NoError(err)
	err = suite.storage.AddOrUpdateUserLocation(ctx, user2)
	suite.Require().NoError(err)

	sharing := &models.LocationSharing{
		ReceiverID: user2.ID,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(1 * time.Hour),
	}

	err = suite.storage.ShareLocation(ctx, user1.ID, sharing)
	suite.Require().NoError(err)

	var count int
	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM location_sharing WHERE sharer_id=$1 AND receiver_id=$2", user1.ID, user2.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(1, count)
}

func (suite *UserStorageTestSuite) TestStopSharingLocation() {
	ctx := context.Background()
	user1 := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := &models.User{
		ID:        2,
		Privacy:   "public",
		H3Index:   "85283473ffffffe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user1)
	suite.Require().NoError(err)
	err = suite.storage.AddOrUpdateUserLocation(ctx, user2)
	suite.Require().NoError(err)

	sharing := &models.LocationSharing{
		ReceiverID: user2.ID,
		StartTime:  time.Now(),
		EndTime:    time.Now().Add(1 * time.Hour),
	}

	err = suite.storage.ShareLocation(ctx, user1.ID, sharing)
	suite.Require().NoError(err)

	err = suite.storage.StopSharingLocation(ctx, user1.ID, user2.ID)
	suite.Require().NoError(err)

	var count int
	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM location_sharing WHERE sharer_id=$1 AND receiver_id=$2 AND end_time IS NOT NULL", user1.ID, user2.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(1, count)
}

func (suite *UserStorageTestSuite) TestGetActiveSharings() {
	ctx := context.Background()
	user1 := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := &models.User{
		ID:        2,
		Privacy:   "public",
		H3Index:   "85283473ffffffe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user1)
	suite.Require().NoError(err)
	err = suite.storage.AddOrUpdateUserLocation(ctx, user2)
	suite.Require().NoError(err)

	sharing := &models.LocationSharing{
		SharerID:   user1.ID,
		ReceiverID: user2.ID,
		StartTime:  time.Now(),
		EndTime:    time.Time{},
	}

	err = suite.storage.ShareLocation(ctx, user1.ID, sharing)
	suite.Require().NoError(err)

	sharings, err := suite.storage.GetActiveSharings(ctx, user1.ID)
	suite.Require().NoError(err)
	suite.Len(sharings, 1)
	suite.Equal(user1.ID, sharings[0].SharerID)
	suite.Equal(user2.ID, sharings[0].ReceiverID)
	suite.True(sharings[0].EndTime.IsZero())

	sharings, err = suite.storage.GetActiveSharings(ctx, user2.ID)
	suite.Require().NoError(err)
	suite.Len(sharings, 1)
	suite.Equal(user1.ID, sharings[0].SharerID)
	suite.Equal(user2.ID, sharings[0].ReceiverID)
	suite.True(sharings[0].EndTime.IsZero())
}

func (suite *UserStorageTestSuite) TestCheckUserIsInActiveSharings() {
	ctx := context.Background()
	user1 := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user2 := &models.User{
		ID:        2,
		Privacy:   "public",
		H3Index:   "85283473ffffffe",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user1)
	suite.Require().NoError(err)
	err = suite.storage.AddOrUpdateUserLocation(ctx, user2)
	suite.Require().NoError(err)

	sharing := &models.LocationSharing{
		ReceiverID: user2.ID,
		StartTime:  time.Now(),
		EndTime:    time.Time{}, // EndTime is NULL for active sharing
	}

	err = suite.storage.ShareLocation(ctx, user1.ID, sharing)
	suite.Require().NoError(err)

	isInActiveSharing, err := suite.storage.CheckUserIsInActiveSharings(ctx, user2.ID, user1.ID)
	suite.Require().NoError(err)
	suite.True(isInActiveSharing)

	// Прекращаем шеринг местоположения
	err = suite.storage.StopSharingLocation(ctx, user1.ID, user2.ID)
	suite.Require().NoError(err)

	isInActiveSharing, err = suite.storage.CheckUserIsInActiveSharings(ctx, user2.ID, user1.ID)
	suite.Require().NoError(err)
	suite.False(isInActiveSharing)
}

func (suite *UserStorageTestSuite) TestUpdatePrivacy() {
	ctx := context.Background()
	user := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user)
	suite.Require().NoError(err)

	// Обновляем настройки приватности
	newPrivacy := "private"
	err = suite.storage.UpdatePrivacy(ctx, user.ID, newPrivacy)
	suite.Require().NoError(err)

	// Проверяем, что настройки приватности обновлены
	retrievedUser, err := suite.storage.GetUser(ctx, user.ID)
	suite.Require().NoError(err)
	suite.Equal(newPrivacy, retrievedUser.Privacy)
}

func (suite *UserStorageTestSuite) TestDeleteUser() {
	ctx := context.Background()
	user := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user)
	suite.Require().NoError(err)

	// Удаляем пользователя
	err = suite.storage.DeleteUser(ctx, user.ID)
	suite.Require().NoError(err)

	// Проверяем, что пользователь удален
	var count int
	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM users_table WHERE id=$1", user.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(0, count)

	// Проверяем, что записи в location_history удалены
	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM location_history WHERE user_id=$1", user.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(0, count)

	// Проверяем, что записи в location_sharing удалены
	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM location_sharing WHERE sharer_id=$1 OR receiver_id=$1", user.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(0, count)
}

func (suite *UserStorageTestSuite) TestDeleteHistory() {
	ctx := context.Background()
	user := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user)
	suite.Require().NoError(err)

	// Удаляем историю местоположений
	err = suite.storage.DeleteHistory(ctx, user.ID)
	suite.Require().NoError(err)

	// Проверяем, что записи в location_history удалены
	var count int
	err = suite.dbPool.QueryRow(ctx, "SELECT COUNT(*) FROM location_history WHERE user_id=$1", user.ID).Scan(&count)
	suite.Require().NoError(err)
	suite.Equal(0, count)
}

func (suite *UserStorageTestSuite) TestAddOrUpdateUserLocationWithInvalidData() {
	ctx := context.Background()
	invalidUser := &models.User{
		ID:        0,
		Privacy:   "public",
		H3Index:   "invalid_h3_index",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, invalidUser)
	suite.Require().Error(err)
}

func (suite *UserStorageTestSuite) TestDatabaseUnavailable() {

	suite.dbPool.Close()
	suite.storage.Close()

	ctx := context.Background()
	user := &models.User{
		ID:        1,
		Privacy:   "public",
		H3Index:   "85283473fffffff",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.storage.AddOrUpdateUserLocation(ctx, user)
	suite.Require().Error(err)

	var errReconnect error
	suite.dbPool, errReconnect = pgxpool.Connect(context.Background(), "postgres://postgres:postgres@localhost:5433/postgres?sslmode=disable")
	suite.Require().NoError(errReconnect)

	suite.storage = storage.NewUserLocationStorage(postgres.NewPostgresDBForTest(suite.dbPool))
}

func TestUserStorageTestSuite(t *testing.T) {
	suite.Run(t, new(UserStorageTestSuite))
}
