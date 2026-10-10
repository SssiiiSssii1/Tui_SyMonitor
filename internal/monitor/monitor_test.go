package monitor

import "testing"

func TestHistoryTrimsToMaxHistory(t *testing.T) {
	var h History

	for i := 0; i < maxHistory+5; i++ {
		h.Add(Snapshot{CPU: float64(i), RAM: float64(i)})
	}

	if len(h.CPU) != maxHistory {
		t.Errorf("CPU history has %d samples, want %d", len(h.CPU), maxHistory)
	}
	if len(h.RAM) != maxHistory {
		t.Errorf("RAM history has %d samples, want %d", len(h.RAM), maxHistory)
	}

	if h.CPU[0] != 5 {
		t.Errorf("oldest CPU sample = %v, want 5 (the first 5 samples dropped)", h.CPU[0])
	}
	if h.RAM[len(h.RAM)-1] != maxHistory+4 {
		t.Errorf("newest RAM sample = %v, want %d", h.RAM[len(h.RAM)-1], maxHistory+4)
	}
}

func TestHistoryStartEmpty(t *testing.T) {
	var h History

	h.Add(Snapshot{CPU: 10, RAM: 20})

	if len(h.CPU) != 1 || len(h.RAM) != 1 {
		t.Fatalf("expected one sample, got CPU=%d RAM=%d", len(h.CPU), len(h.RAM))
	}
	if h.CPU[0] != 10 || h.RAM[0] != 20 {
		t.Errorf("unexpected values: CPU=%v RAM=%v", h.CPU[0], h.RAM[0])
	}
}
