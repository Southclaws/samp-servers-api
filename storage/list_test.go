package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Southclaws/samp-servers-api/types"
)

func TestManager_GetServers(t *testing.T) {
	type args struct {
		page   int
		size   types.PageSize
		sort   types.SortOrder
		by     types.SortColumn
		filter []types.FilterAttribute
	}
	tests := []struct {
		name        string
		args        args
		wantServers []types.ServerCore
		wantErr     bool
	}{
		{
			"v no sort",
			args{1, 0, "", "", []types.FilterAttribute{}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"s4.example.com", "test server 4", 50, 50, "rivershell", "Polish", true, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
				{"s2.example.com", "test server 2", 0, 100, "Grand Larceny", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"v desc",
			args{1, 0, "asc", "", []types.FilterAttribute{}},
			[]types.ServerCore{
				{"s2.example.com", "test server 2", 0, 100, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
				{"s4.example.com", "test server 4", 50, 50, "rivershell", "Polish", true, "0.3.7-R2"},
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"v pass",
			args{1, 0, "", "", []types.FilterAttribute{types.FilterPassword}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
				{"s2.example.com", "test server 2", 0, 100, "Grand Larceny", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"v empty",
			args{1, 0, "", "", []types.FilterAttribute{types.FilterEmpty}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"s4.example.com", "test server 4", 50, 50, "rivershell", "Polish", true, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"v full",
			args{1, 0, "", "", []types.FilterAttribute{types.FilterFull}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
				{"s2.example.com", "test server 2", 0, 100, "Grand Larceny", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"v pass empty",
			args{1, 0, "", "", []types.FilterAttribute{types.FilterPassword, types.FilterEmpty}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"v pass full",
			args{1, 0, "", "", []types.FilterAttribute{types.FilterPassword, types.FilterFull}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
				{"s2.example.com", "test server 2", 0, 100, "Grand Larceny", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"v empty full",
			args{1, 0, "", "", []types.FilterAttribute{types.FilterEmpty, types.FilterFull}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"limit to 1",
			args{1, 1, "", "", []types.FilterAttribute{types.FilterPassword, types.FilterFull}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"get second page",
			args{2, 1, "", "", []types.FilterAttribute{types.FilterPassword, types.FilterFull}},
			[]types.ServerCore{
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
			},
			false,
		},
		{
			"get multiple per page",
			args{1, 2, "", "", []types.FilterAttribute{types.FilterPassword, types.FilterFull}},
			[]types.ServerCore{
				{"s3.example.com", "test server 3", 948, 1000, "Grand Larceny", "English", false, "0.3.7-R2"},
				{"ss.southcla.ws", "Scavenge and Survive Official", 4, 32, "Scavenge & Survive by Southclaws", "English", false, "0.3.7-R2"},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotServers, err := mgr.GetServers(tt.args.page, tt.args.size, tt.args.sort, tt.args.by, tt.args.filter)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantServers, gotServers)
		})
	}
}
