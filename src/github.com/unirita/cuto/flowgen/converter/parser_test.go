package converter

import "testing"

func TestParseString(t *testing.T) {
	s := `
test1 ->
test2 ->
[
	test3,
	test4 -> test5,
	test6
]->test7
`

	head, err := ParseString(s)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if head == nil {
		t.Fatalf("Head job is not exist.")
	}
	headJob, ok := head.(*Job)
	if !ok {
		t.Fatalf("Head is not job element.")
	}
	if headJob.Name() != "test1" {
		t.Errorf("Unexpected head job name[%s]", headJob.Name())
	}

	second := head.Next()
	if second == nil {
		t.Fatalf("Second job is not exist.")
	}
	secondJob, ok := second.(*Job)
	if !ok {
		t.Fatalf("Second is not job element.")
	}
	if secondJob.Name() != "test2" {
		t.Errorf("Unexpected second job name[%s]", secondJob.Name())
	}

	third := second.Next()
	if third == nil {
		t.Fatalf("Third job is not exist.")
	}
	_, ok = third.(*Gateway)
	if !ok {
		t.Fatalf("Third element is not gateway element.")
	}

	forth := third.Next()
	if forth == nil {
		t.Fatalf("Forth job is not exist.")
	}
	forthJob, ok := forth.(*Job)
	if !ok {
		t.Fatalf("Forth is not job element.")
	}
	if forthJob.Name() != "test7" {
		t.Errorf("Unexpected forth job name[%s]", forthJob.Name())
	}
}

func TestParseString_OnlyOneJob(t *testing.T) {
	initIndexes()
	s := "test1"

	head, err := ParseString(s)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if head == nil {
		t.Fatalf("Head job is not exist.")
	}
	if head.ID() != "job1" {
		t.Errorf("Unexpected job id[%s].", head.ID())
	}
	if head.Next() != nil {
		t.Error("Head element has unexpected next element.")
	}
}

func TestParseString_EmptyGateway(t *testing.T) {
	s := "test1->[]->test2"
	_, err := ParseString(s)
	if err == nil {
		t.Error("Error was not detected.")
	}
}

func TestExtractGateways_EmptyPath(t *testing.T) {
	s := "test1->[test2,,test3]->test4"
	_, err := ParseString(s)
	if err == nil {
		t.Error("Error was not detected.")
	}
}

func TestExtractGateways(t *testing.T) {
	s := "test1->[test2,test3]->[test4,test5->test6]->test7"
	afterStr, gws, err := extractGateways(s)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if afterStr != "test1->:gw->:gw->test7" {
		t.Errorf("Gateway expression did not replace correctly. Result[%s]", afterStr)
	}
	if len(gws) != 2 {
		t.Fatalf("Number of gateway[%d] is not expected value[%d].", len(gws), 2)
	}
}

func TestParseJob(t *testing.T) {
	j, err := parseJob("test")
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if j.Name() != "test" {
		t.Errorf("Unexpected job name[%s].", j.Name())
	}
}

func TestParseJob_DetectIrregalCharacters(t *testing.T) {
	_, err := parseJob(`te\st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `\`)
	}
	_, err = parseJob(`te/st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `/`)
	}
	_, err = parseJob(`te:st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `:`)
	}
	_, err = parseJob(`te*st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `*`)
	}
	_, err = parseJob(`te?st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `?`)
	}
	_, err = parseJob(`te<st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `<`)
	}
	_, err = parseJob(`te>st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `>`)
	}
	_, err = parseJob(`te$st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `$`)
	}
	_, err = parseJob(`te&st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `&`)
	}
	_, err = parseJob(`te,st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `,`)
	}
	_, err = parseJob(`te-st`)
	if err == nil {
		t.Errorf("Irregal character[%s] is not detected as error.", `-`)
	}
}

func TestParseGateway_SinglePath(t *testing.T) {
	s := "[test1]"
	p, err := parseGateway(s)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if p == nil {
		t.Fatalf("Return value is nil.")
	}
	if len(p.PathHeads) != 1 {
		t.Fatalf("PathHeads size[%d] must be %d.", len(p.PathHeads), 1)
	}
	if p.PathHeads[0].Name() != "test1" {
		t.Errorf("Unexpected job name[%s].", p.PathHeads[0].Name())
	}
}

func TestParseGateway_SinglePath_MultiJob(t *testing.T) {
	s := "[test1->test2->test3]"
	p, err := parseGateway(s)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if p == nil {
		t.Fatalf("Return value is nil.")
	}
	if len(p.PathHeads) != 1 {
		t.Fatalf("PathHeads size[%d] must be %d.", len(p.PathHeads), 1)
	}

	head := p.PathHeads[0]
	if head == nil {
		t.Fatalf("Head job is not exist.")
	}
	if head.Name() != "test1" {
		t.Errorf("Unexpected head job name[%s]", head.Name())
	}

	second := head.Next()
	if second == nil {
		t.Fatalf("Second job is not exist.")
	}
	secondJob, ok := second.(*Job)
	if !ok {
		t.Fatalf("Second is not job element.")
	}
	if secondJob.Name() != "test2" {
		t.Errorf("Unexpected second job name[%s]", secondJob.Name())
	}

	third := second.Next()
	if third == nil {
		t.Fatalf("Third job is not exist.")
	}
	thirdJob, ok := third.(*Job)
	if !ok {
		t.Fatalf("Third is not job element.")
	}
	if thirdJob.Name() != "test3" {
		t.Errorf("Unexpected third job name[%s]", thirdJob.Name())
	}
}

func TestParseGateway_MultiPath(t *testing.T) {
	s := "[test1,test2,test3]"
	p, err := parseGateway(s)
	if err != nil {
		t.Fatalf("Unexpected error occured: %s", err)
	}
	if p == nil {
		t.Fatalf("Return value is nil.")
	}
	if len(p.PathHeads) != 3 {
		t.Fatalf("PathHeads size[%d] must be %d.", len(p.PathHeads), 3)
	}
	if p.PathHeads[0].Name() != "test1" {
		t.Errorf("Unexpected job name[%s].", p.PathHeads[0].Name())
	}
	if p.PathHeads[1].Name() != "test2" {
		t.Errorf("Unexpected job name[%s].", p.PathHeads[1].Name())
	}
	if p.PathHeads[2].Name() != "test3" {
		t.Errorf("Unexpected job name[%s].", p.PathHeads[2].Name())
	}
}
