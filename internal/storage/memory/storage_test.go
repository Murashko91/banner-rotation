package memory

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBannerStorage(t *testing.T) {
	count := 100

	t.Run("test banner db Creation", func(t *testing.T) {
		memory := NewMemoryStorage()

		wg := &sync.WaitGroup{}
		wg.Add(count)

		// create Banners
		for i := 0; i < count; i++ {
			go func() {
				bannerID, err := memory.CreateBanner(context.Background(), fmt.Sprintf("description %d", i))
				wg.Done()
				if err != nil {
					require.Nilf(t, err, err.Error())
				}
				require.NotEqualf(t, bannerID, 0, fmt.Sprintf("id should be created for newly created banner record: %d", i))
			}()
		}
		wg.Wait()
	})

	t.Run("test slot db Creation", func(t *testing.T) {
		memory := NewMemoryStorage()

		wg := &sync.WaitGroup{}
		wg.Add(count)

		// create Slots
		for i := 0; i < count; i++ {
			go func() {
				slotID, err := memory.CreateSlot(context.Background(), fmt.Sprintf("description %d", i))
				wg.Done()
				if err != nil {
					require.Nilf(t, err, err.Error())
				}
				require.NotEqualf(t, slotID, 0, fmt.Sprintf("id should be created for newly created slot record: %d", i))
			}()
		}
		wg.Wait()
	})

	t.Run("test group db Creation", func(t *testing.T) {
		memory := NewMemoryStorage()

		wg := &sync.WaitGroup{}
		wg.Add(count)

		// create Social groups
		for i := 0; i < count; i++ {
			go func() {
				groupID, err := memory.CreateGroup(context.Background(), fmt.Sprintf("description %d", i))
				wg.Done()
				if err != nil {
					require.Nilf(t, err, err.Error())
				}
				require.NotEqualf(t, groupID, 0, fmt.Sprintf("id should be created for newly created slot record: %d", i))
			}()
		}
		wg.Wait()
	})

	t.Run("test add banner to rotation test", func(t *testing.T) {
		memory := NewMemoryStorage()

		bannerCounter := 10

		groupID, err := memory.CreateGroup(context.Background(), "description 1")

		require.NoErrorf(t, err, "unexpected error")
		slotID, err := memory.CreateSlot(context.Background(), "description 1")
		require.NoErrorf(t, err, "unexpected error")
		wg := &sync.WaitGroup{}
		wg.Add(bannerCounter)

		for i := 0; i < bannerCounter; i++ {
			go func() {
				err := memory.AddBannerToSlot(context.Background(), i+1, 1)
				wg.Done()
				if err != nil {
					require.Nilf(t, err, err.Error())
				}
			}()
		}
		wg.Wait()

		// get banners by slot test
		bannerIDs, err := memory.GetBannersBySlot(context.Background(), 1)
		require.NoErrorf(t, err, "unexpected error")
		require.Equalf(t, len(bannerIDs), bannerCounter, "unexpected length")

		// Test Get Banner statistics
		stats, err := memory.GetBannersStat(context.Background(), slotID, groupID, bannerIDs)
		require.NoErrorf(t, err, "unexpected error")
		require.Equalf(t, len(stats), bannerCounter, "unexpected length")

		// Test clics and shows
		for _, stat := range stats {
			require.Equalf(t, stat.ShowsCount, 0, "unexpected shows")
			require.Equalf(t, stat.ClicksCount, 0, "unexpected clicks")

			err = memory.UpdateClickStat(context.Background(), stat)
			require.NoErrorf(t, err, "unexpected error")
			err = memory.UpdateClickStat(context.Background(), stat)
			require.NoErrorf(t, err, "unexpected error")
			err = memory.UpdateShowStat(context.Background(), stat)
			require.NoErrorf(t, err, "unexpected error")
		}

		stats, err = memory.GetBannersStat(context.Background(), slotID, groupID, bannerIDs)
		require.NoErrorf(t, err, "unexpected error")

		// Test clics and shows after update
		for _, stat := range stats {
			require.Equalf(t, stat.ShowsCount, 1, "unexpected shows")
			require.Equalf(t, stat.ClicksCount, 2, "unexpected clicks")
		}

		// Test remove banner from slot
		err = memory.DeleteBannerFromSlot(context.Background(), bannerIDs[0], slotID)
		require.NoErrorf(t, err, "unexpected error")
		bannerIDs, err = memory.GetBannersBySlot(context.Background(), 1)
		require.NoErrorf(t, err, "unexpected error")
		require.Equalf(t, len(bannerIDs), bannerCounter-1, "unexpected length")
	})
}
