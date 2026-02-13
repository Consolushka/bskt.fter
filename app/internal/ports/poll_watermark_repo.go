package ports

import (
	"IMP/app/internal/core/poll_watermarks"
)

type PollWatermarkRepo interface {
	FirstOrCreate(model poll_watermarks.PollWatermarkModel) (poll_watermarks.PollWatermarkModel, error)
	Update(model poll_watermarks.PollWatermarkModel) (poll_watermarks.PollWatermarkModel, error)
}
