package domain

import (
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
	"time"
)

type BanRecord struct {
	ID       int64     `json:"id"`
	ClubID   int64     `json:"club_id"`
	User     Member    `json:"user"`
	Admin    Member    `json:"admin"`
	Reason   string    `json:"reason,omitempty"`
	BannedAt time.Time `json:"banned_at"`
}

func ToDomainBanRecord(b *clubv1.BanRecord) *BanRecord {
	bannedAt, err := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", b.GetBannedAt())
	if err != nil {
		return nil
	}

	return &BanRecord{
		ID:       b.GetId(),
		ClubID:   b.GetClubId(),
		User:     *UserObjectToMember(b.User),
		Admin:    *UserObjectToMember(b.Admin),
		Reason:   b.GetReason(),
		BannedAt: bannedAt,
	}
}

func ToDomainBanRecordArr(b []*clubv1.BanRecord) []*BanRecord {
	banRecords := make([]*BanRecord, len(b))
	for i, banRecord := range b {
		banRecords[i] = ToDomainBanRecord(banRecord)
	}

	return banRecords
}
