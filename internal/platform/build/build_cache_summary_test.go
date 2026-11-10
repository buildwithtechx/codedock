package build

import (
	"testing"
)

func TestBuildKitCacheObserver(t *testing.T) {
	observer := NewBuildKitCacheObserver()

	sampleLog := `#1 [internal] load build definition from Dockerfile
#1 DONE 0.0s
#2 [builder 1/3] RUN npm install
#2 CACHED
#3 [builder 2/3] COPY . .
#3 DONE 0.2s
#4 [builder 3/3] RUN npm run build
#4 DONE 1.5s
`
	observer.Observe(sampleLog, "stream1")

	summary := observer.Summary()
	expected := "Docker layer cache: 1 hit(s), 2 miss(es)."
	if summary != expected {
		t.Fatalf("expected %q, got %q", expected, summary)
	}
}

func TestBuildKitCacheObserverEmpty(t *testing.T) {
	observer := NewBuildKitCacheObserver()
	observer.Observe("random plain output\nwithout steps\n", "default")
	if summary := observer.Summary(); summary != "" {
		t.Fatalf("expected empty summary, got %q", summary)
	}
}
