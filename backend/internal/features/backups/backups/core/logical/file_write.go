package backups_core_logical

import (
	"context"
	"io"

	storage_files "databasus-backend/internal/features/storages/files"
)

// BackgroundFileWrite is an upload running beside the dump that feeds it through a
// pipe. Every engine drives it the same way: peek at the error while the dump is
// still producing, then take the receipt once the writer is done.
type BackgroundFileWrite struct {
	// Errors stays writable because the engines put a peeked value back and their
	// cancellation cleanup drains it.
	Errors   chan error
	Receipts chan storage_files.WriteReceipt
}

// StartBackgroundFileWrite sends the receipt before the error, so a caller holding
// the error knows the receipt is already waiting.
func StartBackgroundFileWrite(
	ctx context.Context,
	fileStore BackupFileStore,
	reference storage_files.StoredFileReference,
	source *io.PipeReader,
	onFailure context.CancelCauseFunc,
) BackgroundFileWrite {
	write := BackgroundFileWrite{
		Errors:   make(chan error, 1),
		Receipts: make(chan storage_files.WriteReceipt, 1),
	}

	go func() {
		receipt, err := fileStore.WriteFile(ctx, reference, source)
		if err != nil {
			_ = source.CloseWithError(err)
			onFailure(err)
		}

		write.Receipts <- receipt
		write.Errors <- err
	}()

	return write
}
