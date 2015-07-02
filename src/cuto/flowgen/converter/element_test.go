package converter

import "testing"

func initIndexes() {
	jIndex = 1
	gIndex = 1
}

func TestNewJob(t *testing.T) {
	initIndexes()
	j := NewJob("test")
	if j.ID() != "job1" {
		t.Errorf("Job ID[%s] is wrong.", j.ID())
	}
	if j.Name() != "test" {
		t.Errorf("Job name[%s] is wrong.", j.Name())
	}
}

func TestJobNext_Job(t *testing.T) {
	initIndexes()
	j := NewJob("test1")

	j.SetNext(NewJob("test2"))
	if j.Next() == nil {
		t.Fatal("Failed to set next job.")
	}
	if j.Next().ID() != "job2" {
		t.Errorf("Next job's ID[%s] is wrong. Expected [%s]", j.Next().ID(), "job2")
	}
	if _, ok := j.Next().(*Job); !ok {
		t.Fatal("Next element is not job.")
	}
}

func TestJobNext_Gateway(t *testing.T) {
	initIndexes()
	j := NewJob("test")
	p := NewGateway()

	j.SetNext(p)
	if j.Next() == nil {
		t.Fatal("Failed to set next job.")
	}
	if j.Next().ID() != "gw1" {
		t.Errorf("Next gateway's ID[%s] is wrong. Expected [%s]", j.Next().ID(), "gw1")
	}
	if _, ok := j.Next().(*Gateway); !ok {
		t.Fatal("Next element is not gateway.")
	}
}

func TestNewGateway(t *testing.T) {
	initIndexes()
	p := NewGateway()
	if p.ID() != "gw1" {
		t.Errorf("Gateway ID[%s] is wrong.", p.ID())
	}
	if p.PathHeads == nil {
		t.Fatal("PathHeads is not set.")
	}
	if len(p.PathHeads) != 0 {
		t.Errorf("PathHeads size[%d] must be %d.", len(p.PathHeads), 0)
	}
}

func TestGatewayNext_Job(t *testing.T) {
	initIndexes()
	p := NewGateway()

	p.SetNext(NewJob("test"))
	if p.Next() == nil {
		t.Fatal("Failed to set next job.")
	}
	if p.Next().ID() != "job1" {
		t.Errorf("Next job's ID[%s] is wrong. Expected [%s]", p.Next().ID(), "job1")
	}
	if _, ok := p.Next().(*Job); !ok {
		t.Fatal("Next element is not job.")
	}
}

func TestGatewayNext_Gateway(t *testing.T) {
	initIndexes()
	p := NewGateway()

	p.SetNext(NewGateway())
	if p.Next() == nil {
		t.Fatal("Failed to set next job.")
	}
	if p.Next().ID() != "gw2" {
		t.Errorf("Next gateway's ID[%s] is wrong. Expected [%s]", p.Next().ID(), "gw2")
	}
	if _, ok := p.Next().(*Gateway); !ok {
		t.Fatal("Next element is not parallel.")
	}
}

func TestAddPathHead(t *testing.T) {
	initIndexes()
	p := NewGateway()
	p.AddPathHead(NewJob("test1"))
	p.AddPathHead(NewJob("test2"))

	if len(p.PathHeads) != 2 {
		t.Fatalf("PathHeads size[%d] must be %d.", len(p.PathHeads), 2)
	}
	if p.PathHeads[0].ID() != "job1" {
		t.Errorf("PathHeads[0] name[%s] is wrong. Expected [%s]", p.PathHeads[0], "test1")
	}
	if p.PathHeads[1].ID() != "job2" {
		t.Errorf("PathHeads[1] name[%s] is wrong. Expected [%s]", p.PathHeads[1], "test2")
	}
}
