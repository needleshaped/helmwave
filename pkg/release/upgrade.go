package release

import (
	log "github.com/sirupsen/logrus"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	helm "helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
	"os"
	"path/filepath"
)

func (rel *Config) upgrade(helm *helm.EnvSettings) (*release.Release, error) {
	client := rel.newUpgrade()

	locateChart, err := client.ChartPathOptions.LocateChart(rel.Chart.Name, helm)
	if err != nil {
		return nil, err
	}

	valOpts := &values.Options{ValueFiles: rel.values}
	vals, err := valOpts.MergeValues(getter.All(helm))
	if err != nil {
		return nil, err
	}

	ch, err := loader.Load(locateChart)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	if req := ch.Metadata.Dependencies; req != nil {
		if err := action.CheckDependencies(ch, req); err != nil {
			return nil, err
		}
	}

	if !(ch.Metadata.Type == "" || ch.Metadata.Type == "application") {
		log.Warnf("%s charts are not installable \n", ch.Metadata.Type)
	}

	if ch.Metadata.Deprecated {
		log.Warn("⚠️ This locateChart is deprecated")
	}

	if !rel.isInstalled() {
		if rel.dryRun {
			log.Debugf("Templating %q ", rel.UniqName())
		} else {
			log.Debugf("🧐 Release %q does not exist. Installing it now.", rel.UniqName())
		}

		return rel.newInstall().Run(ch, vals)
	}

	return client.Run(rel.Name, ch, vals)

}

func (rel *Config) chartDepsUpd(settings *helm.EnvSettings) error {
	client := action.NewDependency()
	man := &downloader.Manager{
		Out:              os.Stdout,
		ChartPath:        filepath.Clean(rel.Chart.Name),
		Keyring:          client.Keyring,
		SkipUpdate:       client.SkipRefresh,
		Getters:          getter.All(settings),
		RepositoryConfig: settings.RepositoryConfig,
		RepositoryCache:  settings.RepositoryCache,
		Debug:            settings.Debug,
	}
	if client.Verify {
		man.Verify = downloader.VerifyAlways
	}
	return man.Update()
}
