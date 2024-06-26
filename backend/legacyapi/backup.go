package api

// BackupStorageBackend is the storage backend of a backup.
type BackupStorageBackend string

const (
	// BackupStorageBackendLocal is the local storage backend for a backup.
	BackupStorageBackendLocal BackupStorageBackend = "LOCAL"
	// BackupStorageBackendS3 is the AWS S3 storage backend for a backup. Not used yet.
	BackupStorageBackendS3 BackupStorageBackend = "S3"
	// BackupStorageBackendGCS is the Google Cloud Storage (GCS) storage backend for a backup. Not used yet.
	BackupStorageBackendGCS BackupStorageBackend = "GCS"
	// BackupStorageBackendOSS is the AliCloud Object Storage Service (OSS) storage backend for a backup. Not used yet.
	BackupStorageBackendOSS BackupStorageBackend = "OSS"
)

// BinlogInfo is the binlog coordination for MySQL.
type BinlogInfo struct {
	FileName string `json:"fileName"`
	Position int64  `json:"position"`
}

// IsEmpty return true if the BinlogInfo is empty.
func (b BinlogInfo) IsEmpty() bool {
	return b == BinlogInfo{}
}

// BackupPayload contains backup related database specific info, it differs for different database types.
// It is encoded in JSON and stored in the backup table.
type BackupPayload struct {
	// MySQL related fields
	// BinlogInfo is recorded when taking the backup.
	// It is recorded within the same transaction as the dump so that the binlog position is consistent with the dump.
	// Please refer to https://github.com/bytebase/bytebase/blob/main/docs/design/pitr-mysql.md#full-backup for details.
	BinlogInfo BinlogInfo `json:"binlogInfo"`
}
