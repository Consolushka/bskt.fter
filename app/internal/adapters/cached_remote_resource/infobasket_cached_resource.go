package cached_remote_resource

import (
	"IMP/app/internal/infra/infobasket"
	"IMP/app/internal/ports"
	"strconv"
	"time"
)

type InfoBasketCachedResource struct {
	getScheduleGamesFunc func(compId int) ([]infobasket.GameScheduleDto, error)
	compId               int
}

func (i InfoBasketCachedResource) Load() (any, error) {
	data, err := i.getScheduleGamesFunc(i.compId)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (i InfoBasketCachedResource) LocalFileName() string {
	return "infobasket_schedule_games_" + strconv.Itoa(i.compId) + ".json"
}

func (i InfoBasketCachedResource) GetLifeTime() time.Duration {
	return time.Hour * 24 * 4
}

func NewInfobasketCachedResource(getScheduleGamesFunc func(compId int) ([]infobasket.GameScheduleDto, error), compId int) ports.CachedRemoteResource {
	return &InfoBasketCachedResource{
		getScheduleGamesFunc: getScheduleGamesFunc,
		compId:               compId,
	}
}
