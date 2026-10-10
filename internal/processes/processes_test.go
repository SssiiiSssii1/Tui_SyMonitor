package processes

import "testing"

func TestFilterByNameIsCaseInsensitive(t *testing.T) {
	all := []Info{
		{PID: 1, Name: "chrome"},
		{PID: 2, Name: "Chrome Helper"},
		{PID: 3, Name: "bash"},
	}

	got := Filter(all, "chrome")

	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(got))
	}
	if got[0].PID != 1 || got[1].PID != 2 {
		t.Errorf("expected PIDs 1 and 2, got %d and %d", got[0].PID, got[1].PID)
	}
}

func TestFilterByPID(t *testing.T) {
	all := []Info{
		{PID: 123, Name: "bash"},
		{PID: 456, Name: "top"},
	}

	got := Filter(all, "45")

	if len(got) != 1 || got[0].PID != 456 {
		t.Errorf("expected a single match with PID 456, got %+v", got)
	}
}

func TestFilterEmptyQueryReturnsCopy(t *testing.T) {
	all := []Info{{PID: 1, Name: "one"}, {PID: 2, Name: "two"}}

	got := Filter(all, "")

	if len(got) != len(all) {
		t.Fatalf("expected %d results, got %d", len(all), len(got))
	}

	got[0].PID = 99
	if all[0].PID != 1 {
		t.Error("Filter must not mutate or alias the input slice")
	}
}

func TestSortByCPUDescending(t *testing.T) {
	in := []Info{{PID: 1, CPU: 5}, {PID: 2, CPU: 90}, {PID: 3, CPU: 12}}

	Sort(in, "cpu")

	want := []int32{2, 3, 1}
	for i, pid := range want {
		if in[i].PID != pid {
			t.Errorf("position %d: expected PID %d, got %d", i, pid, in[i].PID)
		}
	}
}

func TestSortByRAMDescending(t *testing.T) {
	in := []Info{{PID: 1, RAM: 1.5}, {PID: 2, RAM: 0.2}, {PID: 3, RAM: 3.5}}

	Sort(in, "ram")

	want := []int32{3, 1, 2}
	for i, pid := range want {
		if in[i].PID != pid {
			t.Errorf("position %d: expected PID %d, got %d", i, pid, in[i].PID)
		}
	}
}

func TestSortByNameAscending(t *testing.T) {
	in := []Info{{PID: 1, Name: "zeta"}, {PID: 2, Name: "Alpha"}, {PID: 3, Name: "beta"}}

	Sort(in, "name")

	want := []int32{2, 3, 1}
	for i, pid := range want {
		if in[i].PID != pid {
			t.Errorf("position %d: expected PID %d, got %d", i, pid, in[i].PID)
		}
	}
}

func TestSortUnknownKeyKeepsOrder(t *testing.T) {
	in := []Info{{PID: 1, CPU: 99}, {PID: 2, CPU: 1}}

	Sort(in, "bogus")

	if in[0].PID != 1 || in[1].PID != 2 {
		t.Error("unknown sort key must leave the order unchanged")
	}
}
