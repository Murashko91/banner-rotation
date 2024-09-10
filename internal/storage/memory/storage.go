package memory

import (
	"context"
	"fmt"
	"slices"
	"sync"

	"github.com/otus-murashko/banners-rotation/internal/storage"
)

type Storage struct {
	slotMap    map[int]storage.Slot
	bannerMap  map[int]storage.Banner
	groupMap   map[int]storage.SosialGroup
	rotations  []storage.Rotation
	statistics []storage.Statistic
	mutex      *sync.RWMutex
}

func (s Storage) Connect() error {
	return nil
}

func (s Storage) Close() error {
	return nil
}

func (s Storage) GetBannersBySlot(_ context.Context, slotID int) ([]int, error) {
	s.mutex.RLock()
	result := make([]int, 0)

	for _, rotation := range s.rotations {
		if rotation.SlotID == slotID {
			result = append(result, rotation.BannerID)
		}
	}

	s.mutex.RUnlock()
	return result, nil
}

func (s Storage) GetBannersStat(ctx context.Context, slotID int, groupID int, bannerIDs []int) ([]storage.Statistic, error) {
	s.mutex.RLock()
	result := make([]storage.Statistic, 0)

	fmt.Println(".statistics")

	fmt.Println(s.statistics)

	for _, stat := range s.statistics {

		if stat.SlotID == slotID &&
			slices.Contains(bannerIDs, stat.BannerID) &&
			stat.SosialGroupID == groupID {
			result = append(result, stat)
		}
	}

	s.mutex.RUnlock()
	return result, nil
}

func (s *Storage) AddBannerToSlot(ctx context.Context, bannerID int, slotID int) error {

	s.mutex.Lock()
	hasValue := false
	for _, rotation := range s.rotations {
		if rotation.BannerID == bannerID &&
			rotation.SlotID == slotID {
			hasValue = true
		}
	}

	if !hasValue {

		s.rotations = append(s.rotations, storage.Rotation{BannerID: bannerID, SlotID: slotID})
	}

	existsStatGroupIDs := make([]int, 0)

	for _, stat := range s.statistics {
		if stat.SlotID == slotID &&
			stat.BannerID == bannerID {
			existsStatGroupIDs = append(existsStatGroupIDs, stat.SosialGroupID)
		}
	}

	// Add stats for sosial groups if no exists
	for _, groupID := range getAllGroupIDs(s.groupMap) {
		if !slices.Contains(existsStatGroupIDs, groupID) {
			s.statistics = append(s.statistics,
				storage.Statistic{BannerID: bannerID, SlotID: slotID, SosialGroupID: groupID})
		}
	}

	s.mutex.Unlock()
	return nil
}

func getAllGroupIDs(gMap map[int]storage.SosialGroup) []int {

	result := make([]int, 0, len(gMap))
	for id := range gMap {
		result = append(result, id)

	}
	return result
}

func (s *Storage) DeleteBannerFromSlot(ctx context.Context, bannerID int, slotID int) error {
	s.mutex.Lock()
	hasValue := false
	position := 0

	for i, rotation := range s.rotations {
		if rotation.BannerID == bannerID &&
			rotation.SlotID == slotID {
			hasValue = true
			position = i
		}
	}

	if hasValue {
		s.rotations = slices.Delete(s.rotations, position, position+1)
	}
	s.mutex.Unlock()
	return nil
}

func (s *Storage) CreateBanner(ctx context.Context, desc string) (int, error) {
	s.mutex.Lock()
	newId := getNewId(s.bannerMap)
	s.bannerMap[newId] = storage.Banner{
		ID:    newId,
		Descr: desc,
	}
	s.mutex.Unlock()
	return newId, nil
}

func (s *Storage) CreateSlot(ctx context.Context, desc string) (int, error) {
	s.mutex.Lock()
	newId := getNewId(s.slotMap)
	s.slotMap[newId] = storage.Slot{
		ID:    newId,
		Descr: desc,
	}
	s.mutex.Unlock()
	return newId, nil
}

func (s *Storage) CreateGroup(ctx context.Context, desc string) (int, error) {

	s.mutex.Lock()
	newId := getNewId(s.groupMap)
	s.groupMap[newId] = storage.SosialGroup{
		ID:    newId,
		Descr: desc,
	}
	s.mutex.Unlock()
	return newId, nil
}

func (s *Storage) UpdateShowStat(ctx context.Context, showStat storage.Statistic) error {
	s.mutex.Lock()
	var statInDb storage.Statistic
	hasValue := false
	pos := 0
	for i, stat := range s.statistics {

		if stat.SlotID == showStat.SlotID &&
			stat.BannerID == showStat.BannerID &&
			stat.SosialGroupID == showStat.SosialGroupID {
			pos = i
			hasValue = true
			statInDb = stat
		}
	}

	if hasValue {
		statInDb.ShowsCount++
		s.statistics[pos] = statInDb
	}

	s.mutex.Unlock()
	return nil
}

func (s *Storage) UpdateClickStat(ctx context.Context, showStat storage.Statistic) error {

	s.mutex.Lock()
	var statInDb storage.Statistic
	hasValue := false
	pos := 0
	for i, stat := range s.statistics {

		if stat.SlotID == showStat.SlotID &&
			stat.BannerID == showStat.BannerID &&
			stat.SosialGroupID == showStat.SosialGroupID {
			pos = i
			hasValue = true
			statInDb = stat
		}
	}

	if hasValue {
		statInDb.ClicksCount++
		s.statistics[pos] = statInDb
	}
	s.mutex.Unlock()

	return nil
}

func NewMemoryStorage() *Storage {

	return &Storage{
		slotMap:    make(map[int]storage.Slot),
		bannerMap:  make(map[int]storage.Banner),
		groupMap:   make(map[int]storage.SosialGroup),
		rotations:  make([]storage.Rotation, 0),
		statistics: make([]storage.Statistic, 0),
		mutex:      &sync.RWMutex{},
	}
}

func getNewId[V any](in map[int]V) int {

	maxId := 0
	for key := range in {
		if key > maxId {
			maxId = key
		}
	}

	return maxId + 1
}
