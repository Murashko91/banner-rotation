package storage

import "context"

type Storage interface {
	Connect() error
	Close() error
	GetBannersBySlot(ctx context.Context, slotID int) ([]int, error)
	GetBannersStat(ctx context.Context, slotID int, groupID int, bannerIDs []int) ([]Statistic, error)
	AddBannerToSlot(ctx context.Context, bannerID int, slotID int) error
	DeleteBannerFromSlot(ctx context.Context, bannerID int, slotID int) error
	CreateBanner(ctx context.Context, desc string) (int, error)
	CreateSlot(ctx context.Context, desc string) (int, error)
	CreateGroup(ctx context.Context, desc string) (int, error)
	UpdateShowStat(ctx context.Context, stat Statistic) error
	UpdateClickStat(ctx context.Context, stat Statistic) error
}

type Banner struct {
	ID    int    `db:"id"`
	Descr string `db:"descr"`
}

type Slot struct {
	ID    int    `db:"id"`
	Descr string `db:"descr"`
}

type SosialGroup struct {
	ID    int    `db:"id"`
	Descr string `db:"desc	r"`
}

type Rotation struct {
	BannerID int `db:"banner"`
	SlotID   int `db:"slot"`
}

type Statistic struct {
	BannerID      int `db:"banner"`
	SlotID        int `db:"slot"`
	ClicksCount   int `db:"clicks"`
	ShowsCount    int `db:"shows"`
	SocialGroupID int `db:"s_group"`
}
