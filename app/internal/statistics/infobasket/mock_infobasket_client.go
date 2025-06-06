// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/statistics/infobasket/client.go

// Package mock is a generated GoMock package.
package infobasket

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockClientInterface is a mock of ClientInterface interface.
type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

// MockClientInterfaceMockRecorder is the mock recorder for MockClientInterface.
type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

// NewMockClientInterface creates a new mock instance.
func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

// BoxScore mocks base method.
func (m *MockClientInterface) BoxScore(gameId string) GameBoxScoreResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BoxScore", gameId)
	ret0, _ := ret[0].(GameBoxScoreResponse)
	return ret0
}

// BoxScore indicates an expected call of BoxScore.
func (mr *MockClientInterfaceMockRecorder) BoxScore(gameId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BoxScore", reflect.TypeOf((*MockClientInterface)(nil).BoxScore), gameId)
}

// ScheduledGames mocks base method.
func (m *MockClientInterface) ScheduledGames(compId int) []GameScheduleDto {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduledGames", compId)
	ret0, _ := ret[0].([]GameScheduleDto)
	return ret0
}

// ScheduledGames indicates an expected call of ScheduledGames.
func (mr *MockClientInterfaceMockRecorder) ScheduledGames(compId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduledGames", reflect.TypeOf((*MockClientInterface)(nil).ScheduledGames), compId)
}

// TeamGames mocks base method.
func (m *MockClientInterface) TeamGames(teamId string) TeamScheduleResponse {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TeamGames", teamId)
	ret0, _ := ret[0].(TeamScheduleResponse)
	return ret0
}

// TeamGames indicates an expected call of TeamGames.
func (mr *MockClientInterfaceMockRecorder) TeamGames(teamId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TeamGames", reflect.TypeOf((*MockClientInterface)(nil).TeamGames), teamId)
}
