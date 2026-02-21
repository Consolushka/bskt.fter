package infobasket

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

type MockClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockClientInterfaceMockRecorder
}

type MockClientInterfaceMockRecorder struct {
	mock *MockClientInterface
}

func NewMockClientInterface(ctrl *gomock.Controller) *MockClientInterface {
	mock := &MockClientInterface{ctrl: ctrl}
	mock.recorder = &MockClientInterfaceMockRecorder{mock}
	return mock
}

func (m *MockClientInterface) EXPECT() *MockClientInterfaceMockRecorder {
	return m.recorder
}

func (m *MockClientInterface) BoxScore(gameId string) (GameBoxScoreResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BoxScore", gameId)
	ret0, _ := ret[0].(GameBoxScoreResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockClientInterfaceMockRecorder) BoxScore(gameId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BoxScore", reflect.TypeOf((*MockClientInterface)(nil).BoxScore), gameId)
}

func (m *MockClientInterface) ScheduledGames(compId int) ([]GameScheduleDto, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ScheduledGames", compId)
	ret0, _ := ret[0].([]GameScheduleDto)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockClientInterfaceMockRecorder) ScheduledGames(compId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ScheduledGames", reflect.TypeOf((*MockClientInterface)(nil).ScheduledGames), compId)
}
