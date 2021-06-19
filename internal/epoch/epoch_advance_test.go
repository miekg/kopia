package epoch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/kopia/kopia/repo/blob"
)

func TestShouldAdvanceEpoch(t *testing.T) {
	t0 := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

	var lotsOfMetadata []blob.Metadata

	lotsOfMetadata = append(lotsOfMetadata, blob.Metadata{
		Timestamp: t0, Length: 1,
	})

	for i := 0; i < DefaultParameters.EpochAdvanceOnCountThreshold; i++ {
		lotsOfMetadata = append(lotsOfMetadata, blob.Metadata{
			Timestamp: t0.Add(DefaultParameters.MinEpochDuration),
			Length:    1,
		})
	}

	cases := []struct {
		desc string
		bms  []blob.Metadata
		want bool
	}{
		{
			desc: "zero blobs",
			bms:  []blob.Metadata{},
			want: false,
		},
		{
			desc: "one blob",
			bms: []blob.Metadata{
				{Timestamp: t0, Length: 1},
			},
			want: false,
		},
		{
			desc: "two blobs, not enough time passed, size enough to advance",
			bms: []blob.Metadata{
				{Timestamp: t0.Add(DefaultParameters.MinEpochDuration - 1), Length: DefaultParameters.EpochAdvanceOnTotalSizeBytesThreshold},
				{Timestamp: t0, Length: 1},
			},
			want: false,
		},
		{
			desc: "two blobs, enough time passed, total size enough to advance",
			bms: []blob.Metadata{
				{Timestamp: t0, Length: 1},
				{Timestamp: t0.Add(DefaultParameters.MinEpochDuration), Length: DefaultParameters.EpochAdvanceOnTotalSizeBytesThreshold},
			},
			want: true,
		},
		{
			desc: "two blobs, enough time passed, total size not enough to advance",
			bms: []blob.Metadata{
				{Timestamp: t0, Length: 1},
				{Timestamp: t0.Add(DefaultParameters.MinEpochDuration), Length: DefaultParameters.EpochAdvanceOnTotalSizeBytesThreshold - 2},
			},
			want: false,
		},
		{
			desc: "enough time passed, count not enough to advance",
			bms: []blob.Metadata{
				{Timestamp: t0, Length: 1},
				{Timestamp: t0.Add(DefaultParameters.MinEpochDuration), Length: 1},
			},
			want: false,
		},
		{
			desc: "enough time passed, count enough to advance",
			bms:  lotsOfMetadata,
			want: true,
		},
	}

	for _, tc := range cases {
		require.Equal(t, tc.want,
			shouldAdvance(tc.bms, DefaultParameters.MinEpochDuration, DefaultParameters.EpochAdvanceOnCountThreshold, DefaultParameters.EpochAdvanceOnTotalSizeBytesThreshold),
			tc.desc)
	}
}
