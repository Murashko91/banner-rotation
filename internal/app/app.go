package app

import (
	"context"

	"github.com/otus-murashko/banners-rotation/internal/banner"
	"github.com/otus-murashko/banners-rotation/internal/storage"
)

type Application interface {
	GetBannersBySlot(ctx context.Context, slotID int) ([]int, error)
	GetBannersStat(ctx context.Context, slotID int, groupID int, bannerIDs []int) ([]storage.Statistic, error)
	AddBannerToSlot(ctx context.Context, bannerID int, slotID int) error
	DeleteBannerFromSlot(ctx context.Context, bannerID int, slotID int) error
	CreateBanner(ctx context.Context, desc string) (int, error)
	CreateSlot(ctx context.Context, desc string) (int, error)
	CreateGroup(ctx context.Context, desc string) (int, error)
	GetBannerRotation(ctx context.Context, slotID, sGroupID int) (storage.Banner, error)
	UpdateShowStat(ctx context.Context, stat storage.Statistic) error
	UpdateClickStat(ctx context.Context, stat storage.Statistic) error
}

type BannerSelector interface {
	GetBanner(ctx context.Context, slotID, sGroupID int) (storage.Banner, error)
}

type App struct {
	storage storage.Storage
	bs      BannerSelector
}

func New(storage storage.Storage) *App {
	return &App{
		storage: storage,
		bs:      banner.NewBannerBanditSelector(storage),
	}
}

func (a App) GetBannersBySlot(ctx context.Context, slotID int) ([]int, error) {
	return a.storage.GetBannersBySlot(ctx, slotID)
}

func (a App) GetBannersStat(
	ctx context.Context, slotID int, groupID int, bannerIDs []int,
) ([]storage.Statistic, error) {
	return a.storage.GetBannersStat(ctx, slotID, groupID, bannerIDs)
}

func (a App) AddBannerToSlot(ctx context.Context, bannerID int, slotID int) error {
	return a.storage.AddBannerToSlot(ctx, bannerID, slotID)
}

func (a App) DeleteBannerFromSlot(ctx context.Context, bannerID int, slotID int) error {
	return a.storage.DeleteBannerFromSlot(ctx, bannerID, slotID)
}

func (a App) CreateBanner(ctx context.Context, desc string) (int, error) {
	return a.storage.CreateBanner(ctx, desc)
}

func (a App) CreateSlot(ctx context.Context, desc string) (int, error) {
	return a.storage.CreateSlot(ctx, desc)
}

func (a App) CreateGroup(ctx context.Context, desc string) (int, error) {
	return a.storage.CreateGroup(ctx, desc)
}

func (a App) GetBannerRotation(ctx context.Context, slotID, sGroupID int) (storage.Banner, error) {
	return a.bs.GetBanner(ctx, slotID, sGroupID)
}

func (a App) UpdateShowStat(ctx context.Context, stat storage.Statistic) error {
	return a.storage.UpdateShowStat(ctx, stat)
}

func (a App) UpdateClickStat(ctx context.Context, stat storage.Statistic) error {
	return a.storage.UpdateClickStat(ctx, stat)
}
