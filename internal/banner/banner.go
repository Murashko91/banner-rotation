package banner

import (
	"context"
	"math"

	"github.com/otus-murashko/banners-rotation/internal/storage"
)

type Selector interface {
	GetBanner(ctx context.Context, slotID, sGroupID int) (storage.Banner, error)
}

type BanditSelector struct {
	db storage.Storage
}

func NewBannerBanditSelector(db storage.Storage) BanditSelector {
	return BanditSelector{db: db}
}

func (bs BanditSelector) GetBanner(ctx context.Context, slotID, sGroupID int) (storage.Banner, error) {
	// get all statistic for the banners and social group

	bannerIDs, err := bs.db.GetBannersBySlot(ctx, slotID)
	if err != nil {
		return storage.Banner{}, err
	}

	stats, err := bs.db.GetBannersStat(ctx, slotID, sGroupID, bannerIDs)
	if err != nil {
		return storage.Banner{}, err
	}

	bestStat := getBestBannerStat(stats)

	bs.db.UpdateShowStat(ctx, bestStat) // increase shows count in db

	return storage.Banner{ID: bestStat.BannerID}, nil
}

func getBestBannerStat(stats []storage.Statistic) storage.Statistic {
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

	return bestStat
}
