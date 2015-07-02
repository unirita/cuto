package converter

import "testing"

func TestNewJob(t *testing.T) {
	j := NewJob("test")

	if j.Name != "test" {
		t.Errorf("Job name[%s] is wrong.", j.Name)
	}
}

func TestJobNext_Job(t *testing.T) {
	j1 := NewJob("test1")
	j2 := NewJob("test2")

	j1.SetNext(j2)
	if j1.Next() == nil {
		t.Fatal("Failed to set next job.")
	}

	nextJob, ok := j1.Next().(*Job)
	if !ok {
		t.Fatal("Next element is not job.")
	}
	if nextJob.Name != "test2" {
		t.Errorf("Next job's name[%s] is wrong. Expected [%s]", nextJob.Name, "test2")
	}
}

func TestJobNext_Parallel(t *testing.T) {
	j := NewJob("test")
	p := NewParallel()

	j.SetNext(p)
	if j.Next() == nil {
		t.Fatal("Failed to set next job.")
	}

	_, ok := j.Next().(*Parallel)
	if !ok {
		t.Fatal("Next element is not parallel.")
	}
}

func TestNewParallel(t *testing.T) {
	p := NewParallel()

	if p.PathHeads == nil {
		t.Fatal("PathHeads is not set.")
	}
	if len(p.PathHeads) != 0 {
		t.Errorf("PathHeads size[%d] must be 0.", len(p.PathHeads))
	}
}

func TestParallelNext_Job(t *testing.T) {
	p := NewParallel()
	j := NewJob("test1")

	p.SetNext(j)
	if p.Next() == nil {
		t.Fatal("Failed to set next job.")
	}

	nextJob, ok := p.Next().(*Job)
	if !ok {
		t.Fatal("Next element is not job.")
	}
	if nextJob.Name != "test1" {
		t.Errorf("Next job's name[%s] is wrong. Expected [%s]", nextJob.Name, "test1")
	}
}

func TestParallelNext_Parallel(t *testing.T) {
	p1 := NewParallel()
	p2 := NewParallel()

	p1.SetNext(p2)
	if p1.Next() == nil {
		t.Fatal("Failed to set next job.")
	}

	_, ok := p1.Next().(*Parallel)
	if !ok {
		t.Fatal("Next element is not parallel.")
	}
}

func TestAddPathHead(t *testing.T) {
	p := NewParallel()
	p.AddPathHead(NewJob("test1"))
	p.AddPathHead(NewJob("test2"))

	if len(p.PathHeads) != 2 {
		t.Fatalf("PathHeads size[%d] must be 2.", len(p.PathHeads))
	}
	if p.PathHeads[0].Name != "test1" {
		t.Errorf("PathHeads[0] name[%s] is wrong. Expected [%s]", p.PathHeads[0], "test1")
	}
	if p.PathHeads[1].Name != "test2" {
		t.Errorf("PathHeads[1] name[%s] is wrong. Expected [%s]", p.PathHeads[1], "test2")
	}
}
