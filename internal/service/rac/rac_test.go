package rac

import (
	"os"
	"testing"
)

func TestParsingClusterInfo(t *testing.T) {
	d, err := os.ReadFile("testdata/cluster_info_response")
	if err != nil {
		t.Fatal(err)
	}

	want := "4dc663d6-21c2-44b9-859d-7b262dbe70a6"

	clusterInfo, err := parseClusterInfoResponse(d)
	if err != nil {
		t.Fatal(err)
	}

	if got := clusterInfo["cluster"]; got != want {
		t.Fatalf("clusterId is unexpected, want %s; got: %s", want, got)
	}

}

func TestParsingInfobaseSummaryList(t *testing.T) {
	d, err := os.ReadFile("testdata/infobase_summary_list")
	if err != nil {
		t.Fatal(err)
	}

	name := "kzn02_copy"
	want := Infobase{
		id: "254b9d66-5d0c-4e12-8b41-61815f7ce6dc",
	}

	ls, err := parseInfobaseListResponse(d)
	if err != nil {
		t.Fatal(err)
	}

	got, ok := ls[name]

	if !ok {
		t.Fatalf("not found infobase %s", name)
	}

	if got.id != want.id {
		t.Fatalf("unexpected id, want %s; got: %s", want.id, got.id)
	}

}

func TestParsingSessionList(t *testing.T) {
	d, err := os.ReadFile("testdata/session_list")
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"3b7d1af6-4d33-4363-989a-bf2a7beed809",
		"5697970c-0331-452c-b64b-b273f44af189",
		"696d034b-5396-45ce-b9e9-fef69ff8ad54",
	}

	ls, err := parseSessoinListResponse(d)
	if err != nil {
		t.Fatal(err)
	}

	for i, got := range ls {
		if got != want[i] {
			t.Fatalf("unexpected session id, want %s; got: %s", want[i], got)
		}
	}
}
