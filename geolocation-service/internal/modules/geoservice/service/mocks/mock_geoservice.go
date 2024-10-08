// Code generated by MockGen. DO NOT EDIT.
// Source: geoservice.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	models "progekt/dating-app/geolocation-service/internal/models"
	reflect "reflect"
	time "time"
)

// MockMessageQueuer is a mock of MessageQueuer interface
type MockMessageQueuer struct {
	ctrl     *gomock.Controller
	recorder *MockMessageQueuerMockRecorder
}

// MockMessageQueuerMockRecorder is the mock recorder for MockMessageQueuer
type MockMessageQueuerMockRecorder struct {
	mock *MockMessageQueuer
}

// NewMockMessageQueuer creates a new mock instance
func NewMockMessageQueuer(ctrl *gomock.Controller) *MockMessageQueuer {
	mock := &MockMessageQueuer{ctrl: ctrl}
	mock.recorder = &MockMessageQueuerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMessageQueuer) EXPECT() *MockMessageQueuerMockRecorder {
	return m.recorder
}

// Publish mocks base method
func (m *MockMessageQueuer) Publish(topic string, message []byte) error {
	ret := m.ctrl.Call(m, "Publish", topic, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish
func (mr *MockMessageQueuerMockRecorder) Publish(topic, message interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockMessageQueuer)(nil).Publish), topic, message)
}

// Close mocks base method
func (m *MockMessageQueuer) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockMessageQueuerMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMessageQueuer)(nil).Close))
}

// MockUserStorager is a mock of UserStorager interface
type MockUserStorager struct {
	ctrl     *gomock.Controller
	recorder *MockUserStoragerMockRecorder
}

// MockUserStoragerMockRecorder is the mock recorder for MockUserStorager
type MockUserStoragerMockRecorder struct {
	mock *MockUserStorager
}

// NewMockUserStorager creates a new mock instance
func NewMockUserStorager(ctrl *gomock.Controller) *MockUserStorager {
	mock := &MockUserStorager{ctrl: ctrl}
	mock.recorder = &MockUserStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockUserStorager) EXPECT() *MockUserStoragerMockRecorder {
	return m.recorder
}

// AddOrUpdateUserLocation mocks base method
func (m *MockUserStorager) AddOrUpdateUserLocation(ctx context.Context, location *models.User) error {
	ret := m.ctrl.Call(m, "AddOrUpdateUserLocation", ctx, location)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrUpdateUserLocation indicates an expected call of AddOrUpdateUserLocation
func (mr *MockUserStoragerMockRecorder) AddOrUpdateUserLocation(ctx, location interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrUpdateUserLocation", reflect.TypeOf((*MockUserStorager)(nil).AddOrUpdateUserLocation), ctx, location)
}

// GetUser mocks base method
func (m *MockUserStorager) GetUser(ctx context.Context, userID int64) (*models.User, error) {
	ret := m.ctrl.Call(m, "GetUser", ctx, userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockUserStoragerMockRecorder) GetUser(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserStorager)(nil).GetUser), ctx, userID)
}

// GetLocationHistory mocks base method
func (m *MockUserStorager) GetLocationHistory(ctx context.Context, userID int64, limit int) (*models.User, error) {
	ret := m.ctrl.Call(m, "GetLocationHistory", ctx, userID, limit)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLocationHistory indicates an expected call of GetLocationHistory
func (mr *MockUserStoragerMockRecorder) GetLocationHistory(ctx, userID, limit interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocationHistory", reflect.TypeOf((*MockUserStorager)(nil).GetLocationHistory), ctx, userID, limit)
}

// ShareLocation mocks base method
func (m *MockUserStorager) ShareLocation(ctx context.Context, userID int64, sharing *models.LocationSharing) error {
	ret := m.ctrl.Call(m, "ShareLocation", ctx, userID, sharing)
	ret0, _ := ret[0].(error)
	return ret0
}

// ShareLocation indicates an expected call of ShareLocation
func (mr *MockUserStoragerMockRecorder) ShareLocation(ctx, userID, sharing interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShareLocation", reflect.TypeOf((*MockUserStorager)(nil).ShareLocation), ctx, userID, sharing)
}

// StopSharingLocation mocks base method
func (m *MockUserStorager) StopSharingLocation(ctx context.Context, sharerID, receiverID int64) error {
	ret := m.ctrl.Call(m, "StopSharingLocation", ctx, sharerID, receiverID)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopSharingLocation indicates an expected call of StopSharingLocation
func (mr *MockUserStoragerMockRecorder) StopSharingLocation(ctx, sharerID, receiverID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopSharingLocation", reflect.TypeOf((*MockUserStorager)(nil).StopSharingLocation), ctx, sharerID, receiverID)
}

// GetActiveSharings mocks base method
func (m *MockUserStorager) GetActiveSharings(ctx context.Context, userID int64) ([]models.LocationSharing, error) {
	ret := m.ctrl.Call(m, "GetActiveSharings", ctx, userID)
	ret0, _ := ret[0].([]models.LocationSharing)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveSharings indicates an expected call of GetActiveSharings
func (mr *MockUserStoragerMockRecorder) GetActiveSharings(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveSharings", reflect.TypeOf((*MockUserStorager)(nil).GetActiveSharings), ctx, userID)
}

// UpdatePrivacy mocks base method
func (m *MockUserStorager) UpdatePrivacy(ctx context.Context, userID int64, privacy string) error {
	ret := m.ctrl.Call(m, "UpdatePrivacy", ctx, userID, privacy)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePrivacy indicates an expected call of UpdatePrivacy
func (mr *MockUserStoragerMockRecorder) UpdatePrivacy(ctx, userID, privacy interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePrivacy", reflect.TypeOf((*MockUserStorager)(nil).UpdatePrivacy), ctx, userID, privacy)
}

// DeleteUser mocks base method
func (m *MockUserStorager) DeleteUser(ctx context.Context, userID int64) error {
	ret := m.ctrl.Call(m, "DeleteUser", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockUserStoragerMockRecorder) DeleteUser(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockUserStorager)(nil).DeleteUser), ctx, userID)
}

// DeleteHistory mocks base method
func (m *MockUserStorager) DeleteHistory(ctx context.Context, userID int64) error {
	ret := m.ctrl.Call(m, "DeleteHistory", ctx, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteHistory indicates an expected call of DeleteHistory
func (mr *MockUserStoragerMockRecorder) DeleteHistory(ctx, userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteHistory", reflect.TypeOf((*MockUserStorager)(nil).DeleteHistory), ctx, userID)
}

// CheckUserIsInActiveSharings mocks base method
func (m *MockUserStorager) CheckUserIsInActiveSharings(ctx context.Context, responder, target int64) (bool, error) {
	ret := m.ctrl.Call(m, "CheckUserIsInActiveSharings", ctx, responder, target)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserIsInActiveSharings indicates an expected call of CheckUserIsInActiveSharings
func (mr *MockUserStoragerMockRecorder) CheckUserIsInActiveSharings(ctx, responder, target interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserIsInActiveSharings", reflect.TypeOf((*MockUserStorager)(nil).CheckUserIsInActiveSharings), ctx, responder, target)
}

// MockGeoCacher is a mock of GeoCacher interface
type MockGeoCacher struct {
	ctrl     *gomock.Controller
	recorder *MockGeoCacherMockRecorder
}

// MockGeoCacherMockRecorder is the mock recorder for MockGeoCacher
type MockGeoCacherMockRecorder struct {
	mock *MockGeoCacher
}

// NewMockGeoCacher creates a new mock instance
func NewMockGeoCacher(ctrl *gomock.Controller) *MockGeoCacher {
	mock := &MockGeoCacher{ctrl: ctrl}
	mock.recorder = &MockGeoCacherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockGeoCacher) EXPECT() *MockGeoCacherMockRecorder {
	return m.recorder
}

// UpdateLocation mocks base method
func (m *MockGeoCacher) UpdateLocation(user models.User, maxAge time.Duration) error {
	ret := m.ctrl.Call(m, "UpdateLocation", user, maxAge)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateLocation indicates an expected call of UpdateLocation
func (mr *MockGeoCacherMockRecorder) UpdateLocation(user, maxAge interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateLocation", reflect.TypeOf((*MockGeoCacher)(nil).UpdateLocation), user, maxAge)
}

// DeleteUser mocks base method
func (m *MockGeoCacher) DeleteUser(userID int64) error {
	ret := m.ctrl.Call(m, "DeleteUser", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser
func (mr *MockGeoCacherMockRecorder) DeleteUser(userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockGeoCacher)(nil).DeleteUser), userID)
}

// GetAllUsersInH3Cell mocks base method
func (m *MockGeoCacher) GetAllUsersInH3Cell(h3Index string) ([]int64, error) {
	ret := m.ctrl.Call(m, "GetAllUsersInH3Cell", h3Index)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllUsersInH3Cell indicates an expected call of GetAllUsersInH3Cell
func (mr *MockGeoCacherMockRecorder) GetAllUsersInH3Cell(h3Index interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllUsersInH3Cell", reflect.TypeOf((*MockGeoCacher)(nil).GetAllUsersInH3Cell), h3Index)
}

// GetUserLocation mocks base method
func (m *MockGeoCacher) GetUserLocation(userID int64) (*models.User, error) {
	ret := m.ctrl.Call(m, "GetUserLocation", userID)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserLocation indicates an expected call of GetUserLocation
func (mr *MockGeoCacherMockRecorder) GetUserLocation(userID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserLocation", reflect.TypeOf((*MockGeoCacher)(nil).GetUserLocation), userID)
}

// SetPrivacy mocks base method
func (m *MockGeoCacher) SetPrivacy(ctx context.Context, userID int64, privacy string, maxAge time.Duration) error {
	ret := m.ctrl.Call(m, "SetPrivacy", ctx, userID, privacy, maxAge)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPrivacy indicates an expected call of SetPrivacy
func (mr *MockGeoCacherMockRecorder) SetPrivacy(ctx, userID, privacy, maxAge interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPrivacy", reflect.TypeOf((*MockGeoCacher)(nil).SetPrivacy), ctx, userID, privacy, maxAge)
}

// GetUsersWithPrivacy mocks base method
func (m *MockGeoCacher) GetUsersWithPrivacy(privacy string) ([]int64, error) {
	ret := m.ctrl.Call(m, "GetUsersWithPrivacy", privacy)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUsersWithPrivacy indicates an expected call of GetUsersWithPrivacy
func (mr *MockGeoCacherMockRecorder) GetUsersWithPrivacy(privacy interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUsersWithPrivacy", reflect.TypeOf((*MockGeoCacher)(nil).GetUsersWithPrivacy), privacy)
}

// MockPaymentServicer is a mock of PaymentServicer interface
type MockPaymentServicer struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentServicerMockRecorder
}

// MockPaymentServicerMockRecorder is the mock recorder for MockPaymentServicer
type MockPaymentServicerMockRecorder struct {
	mock *MockPaymentServicer
}

// NewMockPaymentServicer creates a new mock instance
func NewMockPaymentServicer(ctrl *gomock.Controller) *MockPaymentServicer {
	mock := &MockPaymentServicer{ctrl: ctrl}
	mock.recorder = &MockPaymentServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPaymentServicer) EXPECT() *MockPaymentServicerMockRecorder {
	return m.recorder
}

// CreateCustomer mocks base method
func (m *MockPaymentServicer) CreateCustomer(user *models.User) (*models.User, error) {
	ret := m.ctrl.Call(m, "CreateCustomer", user)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCustomer indicates an expected call of CreateCustomer
func (mr *MockPaymentServicerMockRecorder) CreateCustomer(user interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCustomer", reflect.TypeOf((*MockPaymentServicer)(nil).CreateCustomer), user)
}

// CreateSubscription mocks base method
func (m *MockPaymentServicer) CreateSubscription(user *models.User) error {
	ret := m.ctrl.Call(m, "CreateSubscription", user)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateSubscription indicates an expected call of CreateSubscription
func (mr *MockPaymentServicerMockRecorder) CreateSubscription(user interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSubscription", reflect.TypeOf((*MockPaymentServicer)(nil).CreateSubscription), user)
}

// GetSubscriptionEndDate mocks base method
func (m *MockPaymentServicer) GetSubscriptionEndDate(user *models.User) (time.Time, error) {
	ret := m.ctrl.Call(m, "GetSubscriptionEndDate", user)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionEndDate indicates an expected call of GetSubscriptionEndDate
func (mr *MockPaymentServicerMockRecorder) GetSubscriptionEndDate(user interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionEndDate", reflect.TypeOf((*MockPaymentServicer)(nil).GetSubscriptionEndDate), user)
}

// GetSubscriptionStatus mocks base method
func (m *MockPaymentServicer) GetSubscriptionStatus(user *models.User) (string, error) {
	ret := m.ctrl.Call(m, "GetSubscriptionStatus", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSubscriptionStatus indicates an expected call of GetSubscriptionStatus
func (mr *MockPaymentServicerMockRecorder) GetSubscriptionStatus(user interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSubscriptionStatus", reflect.TypeOf((*MockPaymentServicer)(nil).GetSubscriptionStatus), user)
}
