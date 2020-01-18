package ru_nalog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocAppend(t *testing.T) {
	d := NewDoc(0, FDCheck)
	assert.NoError(t, d.AppendNew(1054, 1).Err())
	assert.NoError(t, d.AppendNew(1055, 2).Err())
	assert.NoError(t, d.AppendNew(1008, "e@ma.il").Err())
	assert.NoError(t, d.AppendNew(1036, 102030).Err())
	row := d.AppendNew(1059, nil)
	assert.NoError(t, row.AppendNew(1023, 1).Err())
	assert.NoError(t, row.AppendNew(1030, "item").Err())
	assert.NoError(t, row.AppendNew(1079, 7).Err())
	assert.NoError(t, row.AppendNew(1199, 6).Err())
	assert.NoError(t, row.AppendNew(1212, 1).Err())
	assert.NoError(t, row.AppendNew(1214, 1).Err())
	assert.Equal(t, "Doc(#0 Type=3 Props=[(#1054 1) (#1055 2) (#1008 e@ma.il) (#1036 102030) (#1059 [(#1023 1) (#1030 item) (#1079 7) (#1199 6) (#1212 1) (#1214 1)])])", d.String())
}
