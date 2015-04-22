package parser

import (
	"strings"
	"testing"
)

func TestParseNetworkFile_ファイルが存在しない場合はエラーが発生する(t *testing.T) {
	if _, err := ParseNetworkFile("noexistfilepath"); err == nil {
		t.Error("エラーが発生しなかった。")
	}
}

func TestParseNetwork_ネットワーク定義XMLをパースできる(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
  <process id="sample" name="Sample" isExecutable="true">
    <startEvent id="startevent1" name="Start"></startEvent>
    <endEvent id="endevent1" name="End"></endEvent>
    <serviceTask id="servicetask1" name="JOB1"></serviceTask>
    <parallelGateway id="parallelgateway1" name="Parallel Gateway"></parallelGateway>
    <serviceTask id="servicetask2" name="JOB2"></serviceTask>
    <sequenceFlow id="flow1" sourceRef="startevent1" targetRef="servicetask1"></sequenceFlow>
    <serviceTask id="servicetask3" name="JOB3"></serviceTask>
    <sequenceFlow id="flow2" sourceRef="servicetask1" targetRef="parallelgateway1"></sequenceFlow>
    <sequenceFlow id="flow3" sourceRef="parallelgateway1" targetRef="servicetask2"></sequenceFlow>
    <sequenceFlow id="flow4" sourceRef="parallelgateway1" targetRef="servicetask3"></sequenceFlow>
    <parallelGateway id="parallelgateway2" name="Parallel Gateway"></parallelGateway>
    <sequenceFlow id="flow5" sourceRef="servicetask2" targetRef="parallelgateway2"></sequenceFlow>
    <sequenceFlow id="flow6" sourceRef="servicetask3" targetRef="parallelgateway2"></sequenceFlow>
    <sequenceFlow id="flow7" sourceRef="parallelgateway2" targetRef="endevent1"></sequenceFlow>
  </process>
</definitions>`

	r := strings.NewReader(xml)
	proc, err := ParseNetwork(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生: %s", err)
	}

	se := proc.Start
	if len(se) != 1 {
		t.Fatalf("startEventが%d個にも関わらず、%d個取得された", 1, len(se))
	}
	if se[0].ID != "startevent1" {
		t.Errorf("startEventのidは%sのはずが、%sが取得された", "startevent1", se[0].ID)
	}

	ee := proc.End
	if len(ee) != 1 {
		t.Fatalf("endEventが%d個にも関わらず、%d個取得された", 1, len(ee))
	}
	if ee[0].ID != "endevent1" {
		t.Errorf("endEventのidは%sのはずが、%sが取得された", "endevent1", ee[0].ID)
	}

	st := proc.Task
	if len(st) != 3 {
		t.Fatalf("serviceTaskが%d個にも関わらず、%d個取得された", 3, len(st))
	}
	if st[0].ID != "servicetask1" {
		t.Errorf("1つめのserviceTaskのidは%sのはずが、%sが取得された", "servicetask1", st[0].ID)
	}
	if st[1].ID != "servicetask2" {
		t.Errorf("2つめのserviceTaskのidは%sのはずが、%sが取得された", "servicetask2", st[1].ID)
	}

	pg := proc.Gateway
	if len(pg) != 2 {
		t.Fatalf("parallelGatewayが%d個にも関わらず、%d個取得された", 2, len(pg))
	}
	if pg[0].ID != "parallelgateway1" {
		t.Errorf("1つめのparallelgatewayのidは%sのはずが、%sが取得された", "parallelgateway1", pg[0].ID)
	}
	if pg[1].ID != "parallelgateway2" {
		t.Errorf("2つめのparallelgatewayのidは%sのはずが、%sが取得された", "parallelgateway2", pg[1].ID)
	}

	sf := proc.Flow
	if len(sf) != 7 {
		t.Fatalf("sequenceFlowが%d個にも関わらず、%d個取得された", 7, len(sf))
	}
	if sf[0].From != "startevent1" {
		t.Errorf("1つめのsequenceFlowのsourceRefは%sのはずが、%sが取得された", "startevent1", sf[0].From)
	}
	if sf[0].To != "servicetask1" {
		t.Errorf("1つめのsequenceFlowのtargetRefは%sのはずが、%sが取得された", "servicetask1", sf[0].To)
	}
	if sf[1].From != "servicetask1" {
		t.Errorf("1つめのsequenceFlowのsourceRefは%sのはずが、%sが取得された", "servicetask1", sf[1].From)
	}
	if sf[1].To != "parallelgateway1" {
		t.Errorf("1つめのsequenceFlowのtargetRefは%sのはずが、%sが取得された", "parallelgateway1", sf[1].To)
	}
}

func TestParseNetwork_必要最低限の要素があればエラーを吐かない(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
  <process>
    <startEvent/>
    <endEvent/>
  </process>
</definitions>`

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err != nil {
		t.Fatalf("想定外のエラーが発生した: %s", err)
	}
}

func TestParseNetwork_XMLの書式エラー時にエラーを吐く(t *testing.T) {
	xml := "bad_xml"

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err == nil {
		t.Fatal("エラーが返されなかった。")
	}
}

func TestParseNetwork_process要素が無い場合にエラーを吐く(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
</definitions>`

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err == nil {
		t.Fatal("エラーが返されなかった。")
	}
}

func TestParseNetwork_process要素が多すぎる場合にエラーを吐く(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
  <process/>
  <process/>
</definitions>`

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err == nil {
		t.Fatal("エラーが返されなかった。")
	}
}

func TestParseNetwork_startEvent要素が無い場合にエラーを吐く(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
  <process>
    <endEvent/>
  </process>
</definitions>`

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err == nil {
		t.Fatal("エラーが返されなかった。")
	}
}

func TestParseNetwork_startEvent要素が多すぎる場合にエラーを吐く(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
  <process>
    <startEvent/>
    <startEvent/>
    <endEvent/>
  </process>
</definitions>`

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err == nil {
		t.Fatal("エラーが返されなかった。")
	}
}

func TestParseNetwork_endEvent要素が無い場合にエラーを吐く(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
  <process>
    <startEvent/>
  </process>
</definitions>`

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err == nil {
		t.Fatal("エラーが返されなかった。")
	}
}

func TestParseNetwork_endEvent要素が多すぎる場合にエラーを吐く(t *testing.T) {
	xml := `
<?xml version="1.0" encoding="UTF-8"?>
<definitions>
  <process>
    <startEvent/>
    <endEvent/>
    <endEvent/>
  </process>
</definitions>`

	r := strings.NewReader(xml)
	_, err := ParseNetwork(r)
	if err == nil {
		t.Fatal("エラーが返されなかった。")
	}
}
