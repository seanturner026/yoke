package argocd

import (
	"embed"
	"fmt"
	"io/fs"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"

	"github.com/yokecd/yoke/pkg/helm"
)

//go:embed argo-cd-*.tgz
var embedded embed.FS

// loadEmbeddedChart parses the embedded file system for the ArgoCD Helm Chart tgz.
func loadEmbeddedChart() (string, []byte, error) {
	matches, err := fs.Glob(embedded, "argo-cd-*.tgz")
	if err != nil {
		return "", nil, err
	}

	if len(matches) != 1 {
		return "", nil, fmt.Errorf("expected exactly one embedded ArgoCD chart, found %d", len(matches))
	}

	archive, err := embedded.ReadFile(matches[0])
	if err != nil {
		return "", nil, err
	}

	return matches[0], archive, nil
}

// RenderChart renders the chart downloaded from https://argoproj.github.io/argo-helm/argo-cd
// See embedded chart file name for version.
func RenderChart(release, namespace string, values map[string]any) ([]*unstructured.Unstructured, error) {
	_, archive, err := loadEmbeddedChart()
	if err != nil {
		return nil, fmt.Errorf("failed to identify zipped archive: %w", err)
	}
	chart, err := helm.LoadChartFromZippedArchive(archive)
	if err != nil {
		return nil, fmt.Errorf("failed to load chart from zipped archive: %w", err)
	}

	return chart.Render(release, namespace, values)
}
