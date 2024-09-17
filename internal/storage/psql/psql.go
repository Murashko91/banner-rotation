package psql

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/otus-murashko/banners-rotation/internal/storage"
)

type Storage struct {
	info StorageInfo
	db   *sqlx.DB
}

type StorageInfo struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

func New(psqlInfo StorageInfo) *Storage {
	return &Storage{info: psqlInfo}
}

func (s *Storage) Connect() error {
	db, err := sqlx.Open("pgx", getPsqlString(s.info))
	if err != nil {
		return err
	}
	s.db = db
	return nil
}

func (s *Storage) Close() error {
	s.db.Close()
	return nil
}

func (s *Storage) GetBannersBySlot(ctx context.Context, slotID int) ([]int, error) {
	sql := `SELECT banner
	FROM rotation 
	WHERE slot = $1`

	rows, err := s.db.QueryContext(ctx, sql, slotID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	banners := make([]int, 0)
	errorsStr := make([]string, 0)
	for rows.Next() {
		var qBannerID int

		err := rows.Scan(&qBannerID)
		if err != nil {
			errorsStr = append(errorsStr, err.Error())
			continue
		}

		banners = append(banners, qBannerID)
	}
	return banners, getQueryError(errorsStr)
}

func (s *Storage) GetBannersStat(ctx context.Context, slotID int, groupID int, bannerIDs []int) (
	[]storage.Statistic, error,
) {
	sql := `SELECT banner, slot, clicks, shows, s_group
	FROM statistic 
	WHERE slot = $1 AND s_group = $2 AND banner = any($3)`

	rows, err := s.db.QueryxContext(ctx, sql, slotID, groupID, pq.Array(bannerIDs))
	if err != nil {
		return []storage.Statistic{}, err
	}
	defer rows.Close()

	stats := make([]storage.Statistic, 0)
	errorsStr := make([]string, 0)
	for rows.Next() {
		var qStat storage.Statistic

		err := rows.StructScan(&qStat)
		if err != nil {
			errorsStr = append(errorsStr, err.Error())
			continue
		}

		stats = append(stats, qStat)
	}
	return stats, getQueryError(errorsStr)
}

func (s *Storage) UpdateShowStat(ctx context.Context, stat storage.Statistic) error {
	sql := `UPDATE statistic SET 
			shows = shows + 1
 			where banner = $1 AND slot = $2 AND s_group = $3;`

	_, err := s.db.ExecContext(ctx, sql, stat.BannerID, stat.SlotID, stat.SocialGroupID)

	return err
}

func (s *Storage) UpdateClickStat(ctx context.Context, stat storage.Statistic) error {
	sql := `UPDATE statistic SET 
			clicks = clicks + 1
 			where banner = $1 AND slot = $2 AND s_group = $3;`

	_, err := s.db.ExecContext(ctx, sql, stat.BannerID, stat.SlotID, stat.SocialGroupID)

	return err
}

func (s *Storage) AddBannerToSlot(ctx context.Context, bannerID int, slotID int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	sql := `INSERT INTO rotation(banner, slot)
		 	VALUES($1, $2) ON CONFLICT (banner, slot) DO NOTHING `

	// Insert to Slot
	_, err = s.db.ExecContext(ctx, sql, bannerID, slotID)
	if err != nil {
		return err
	}

	// Get all social groups

	sql = `SELECT id
		   FROM social_group`
	rows, err := s.db.QueryContext(ctx, sql)
	if err != nil {
		return err
	}
	defer rows.Close()

	groupIDs := make([]int, 0)
	errorsStr := make([]string, 0)
	for rows.Next() {
		var groupID int

		err := rows.Scan(&groupID)
		if err != nil {
			errorsStr = append(errorsStr, err.Error())
			continue
		}

		groupIDs = append(groupIDs, groupID)
	}

	if len(errorsStr) > 0 {
		return getQueryError(errorsStr)
	}

	if len(groupIDs) == 0 {
		return fmt.Errorf("no groups created in DB")
	}

	sql = "INSERT INTO statistic(banner, slot, s_group) VALUES " +
		createInsertStatValues(groupIDs) +
		"ON CONFLICT (banner, slot, s_group) DO NOTHING "

	// Create empty statistic for all sosial groups

	_, err = s.db.ExecContext(ctx, sql, bannerID, slotID)
	if err != nil {
		return err
	}

	err = tx.Commit()

	return err
}

func createInsertStatValues(groupIDs []int) string {
	sb := strings.Builder{}

	for i, groupID := range groupIDs {
		sb.Write([]byte(fmt.Sprintf("($1, $2, %d)", groupID)))
		if i < len(groupIDs)-1 {
			sb.Write([]byte(", "))
		}
	}

	return sb.String()
}

func (s *Storage) DeleteBannerFromSlot(ctx context.Context, bannerID int, slotID int) error {
	sql := `DELETE from rotation where banner = $1 and slot = $2;`

	_, err := s.db.ExecContext(ctx, sql, bannerID, slotID)
	if err != nil {
		return err
	}

	return err
}

func (s *Storage) CreateBanner(ctx context.Context, desc string) (int, error) {
	return createInstance(ctx, s.db, "banner", desc)
}

func (s *Storage) CreateSlot(ctx context.Context, desc string) (int, error) {
	return createInstance(ctx, s.db, "slot", desc)
}

func (s *Storage) CreateGroup(ctx context.Context, desc string) (int, error) {
	return createInstance(ctx, s.db, "social_group", desc)
}

func createInstance(ctx context.Context, db *sqlx.DB, tNmae, desc string) (int, error) {
	sql := fmt.Sprintf("INSERT INTO %s(descr) VALUES($1) RETURNING id", tNmae)

	lastInsertID := 0

	row := db.QueryRowContext(ctx, sql, desc)

	if row.Err() != nil {
		return 0, row.Err()
	}

	err := row.Scan(&lastInsertID)

	return lastInsertID, err
}

func getPsqlString(dbConfig StorageInfo) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName)
}

func getQueryError(errorsStr []string) error {
	if len(errorsStr) == 0 {
		return nil
	}

	return fmt.Errorf("get banners error: , %v", strings.Join(errorsStr, ";"))
}
