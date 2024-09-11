package banner

import (
	"testing"

	"github.com/otus-murashko/banners-rotation/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestPickBanner(t *testing.T) {
	tests := []struct {
		name   string
		input  []storage.Statistic
		expect int
	}{
		{
			name: "no shows - no clicks",
			input: []storage.Statistic{
				{BannerID: 1, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 2, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 3, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 4, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 5, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
			},
			expect: 1,
		},
		{
			name: "no clicks, 1st banner 1 time showed",
			input: []storage.Statistic{
				{BannerID: 1, ShowsCount: 1, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 2, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 3, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 4, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 5, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
			},
			expect: 2,
		},
		{
			name: "no clicks, 2 banners 1 time showed",
			input: []storage.Statistic{
				{BannerID: 1, ShowsCount: 1, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 2, ShowsCount: 1, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 3, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 4, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 5, ShowsCount: 0, ClicksCount: 0, SosialGroupID: 1},
			},
			expect: 3,
		},

		{
			name: "all banners have shows, but no clicks",
			input: []storage.Statistic{
				{BannerID: 1, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 2, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 3, ShowsCount: 5, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 4, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 5, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
			},
			expect: 3,
		},
		{
			name: "second banner has clicks, pick up the banner",
			input: []storage.Statistic{
				{BannerID: 1, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 2, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 3, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 4, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
				{BannerID: 5, ShowsCount: 10, ClicksCount: 0, SosialGroupID: 1},
			},
			expect: 2,
		},
		{
			name: "5th banner has less shows but clicks count the same",
			input: []storage.Statistic{
				{BannerID: 1, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 2, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 3, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 4, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 5, ShowsCount: 9, ClicksCount: 3, SosialGroupID: 1},
			},
			expect: 5,
		},
		{
			name: "all banners have the same counter clicks,  but BannerID:4 shows count less than others",
			input: []storage.Statistic{
				{BannerID: 1, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 2, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 3, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 4, ShowsCount: 5, ClicksCount: 3, SosialGroupID: 1},
				{BannerID: 5, ShowsCount: 10, ClicksCount: 3, SosialGroupID: 1},
			},
			expect: 4,
		},
	}
	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			bannerID := getBanner(testCase.input)
			require.Equal(t, testCase.expect, bannerID)
		})
	}
}

func getBanner(inputs []storage.Statistic) int {

	return getBestBannerStat(inputs).BannerID

}
