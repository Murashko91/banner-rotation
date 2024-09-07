package banner

import (
	"context"
	"math"

	"github.com/otus-murashko/banners-rotation/internal/storage"
)

type BannerSelector interface {
	GetBanner(ctx context.Context, slotID, sGroupID int) (storage.Banner, error)
}

type BannerBanditSelector struct {
	db storage.Storage
}

func NewBannerBanditSelector(db storage.Storage) BannerBanditSelector {
	return BannerBanditSelector{db: db}
}

func (bs BannerBanditSelector) GetBanner(ctx context.Context, slotID, sGroupID int) (storage.Banner, error) {

	// get all statistic for the banners and social group

	banners, err := bs.db.GetBannersBySlot(ctx, slotID)

	if err != nil {
		return storage.Banner{}, err
	}

	bannerIDs := make([]int, 0, len(banners))
	bannersMap := make(map[int]storage.Banner)

	for _, banner := range banners {
		bannerIDs = append(bannerIDs, banner.ID)
		bannersMap[banner.ID] = banner
	}

	stats, err := bs.db.GetBannersStat(ctx, slotID, sGroupID, bannerIDs)

	if err != nil {
		return storage.Banner{}, err
	}

	totalShowsCount := 0

	// count total shows
	for _, stat := range stats {
		totalShowsCount += stat.ShowsCount
	}

	lnFormulaPart := 2 * math.Log(float64(totalShowsCount))

	bestStat := storage.Statistic{}
	var bestBannerWeight float64

	for _, stat := range stats {
		totalShowsCount += stat.ShowsCount
		if stat.ShowsCount == 0 {
			bestStat = stat
			break
		}

		weight := float64(stat.ClicksCount)/float64(stat.ShowsCount) +
			math.Sqrt(2*lnFormulaPart/float64(stat.ShowsCount))

		if weight > float64(bestBannerWeight) {
			bestBannerWeight = weight
			bestStat = stat
		}
	}

	bs.db.UpdateShowStat(ctx, bestStat)

	return bannersMap[bestStat.BannerID], nil

}
