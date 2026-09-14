package storage_files

import (
	"time"

	"github.com/google/uuid"
)

// Row presence is the whole state: there is no status column, because a row exists
// exactly while nobody has committed to keeping the file it names.
type PendingDeletion struct {
	ID        uuid.UUID `gorm:"column:id;primaryKey;type:uuid;default:gen_random_uuid()"`
	StorageID uuid.UUID `gorm:"column:storage_id;not null;type:uuid"`
	FileName  string    `gorm:"column:file_name;not null;type:text"`
	NotBefore time.Time `gorm:"column:not_before;not null"`
	// Retry counter and receipt generation at once: every event that takes the
	// file away from its writer increments it, which invalidates the receipt the
	// writer is holding.
	AttemptCount int       `gorm:"column:attempt_count;not null"`
	LastError    *string   `gorm:"column:last_error;type:text"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
	UpdatedAt    time.Time `gorm:"column:updated_at;not null"`
}

func (PendingDeletion) TableName() string {
	return "storage_pending_deletions"
}

func (p *PendingDeletion) GetReference() StoredFileReference {
	return StoredFileReference{StorageID: p.StorageID, FileName: p.FileName}
}
