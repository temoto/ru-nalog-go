package umka

import (
	"testing"
	"testing/quick"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCycleAge(t *testing.T) {
	t.Parallel()

	real := time.Now()
	err := quick.Check(func(idt, iopened, iclosed int16) bool {
		dt := real.Add(time.Duration(idt) * time.Second)
		opened := real.Add(time.Duration(iopened) * time.Second)
		closed := real.Add(time.Duration(iclosed) * time.Second)
		st := &Status{Dt: dt.Format(TimeLayout)}
		if iopened%7 != 0 {
			st.CycleOpened = opened.Format(TimeLayout)
		}
		if iclosed%7 != 0 {
			st.CycleClosed = closed.Format(TimeLayout)
		}
		// TODO st.CycleClosed="" if simulating opened
		age, err := st.CycleAge()
		t.Logf("dt=%s opened=%s closed=%s -> age=%v err=%v", st.Dt, st.CycleOpened, st.CycleClosed, age, err)
		if st.CycleOpened == "" {
			assert.NoError(t, err)
			return assert.Equal(t, time.Duration(-1), age)
		}
		if st.CycleClosed != "" {
			if closed.After(dt) {
				return assert.Error(t, err)
			}
			assert.NoError(t, err)
			return assert.Equal(t, -dt.Sub(closed), age)
		}
		if opened.After(dt) {
			return assert.Error(t, err)
		}
		assert.NoError(t, err)
		return assert.Equal(t, dt.Sub(opened), age)
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
}
