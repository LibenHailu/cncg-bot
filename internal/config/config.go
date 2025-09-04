package config

import (
	"os"
)

type Config struct {
	Telegram struct {
		BotToken  string `mapstructure:"bot_token"`
		ChannelID string `mapstructure:"channel_id"`
		ParseMode string `mapstructure:"parse_mode"`
	}
	Scheduler struct {
		CronSpec  string `mapstructure:"cron_spec"`
		BatchSize int    `mapstructure:"batch_size"`
	}
	Filters struct {
		MaxAgeDays int     `mapstructure:"max_age_days"`
		MinScore   float64 `mapstructure:"min_score"`
	}
	Keywords struct {
		Positive []string
		Negative []string
	}
	Sources []struct {
		Name   string
		Type   string
		URL    string
		Weight float64
		Tags   []string
	}
	DBPath string
}

func Load() Config {
	var cfg Config

	if cfg.Telegram.BotToken == "" {
		cfg.Telegram.BotToken = os.Getenv("TOKEN")
	}

	cfg.Telegram.ChannelID = os.Getenv("CHANNEL_ID")
	cfg.Telegram.ParseMode = "MarkdownV2"

	// Scheduler
	cfg.Scheduler.CronSpec = "0 9 * * *" // daily at 09:00
	cfg.Scheduler.BatchSize = 10

	// Filters
	cfg.Filters.MaxAgeDays = 21
	cfg.Filters.MinScore = 0.6

	// Keywords
	cfg.Keywords.Positive = []string{
		"kubernetes", "cncf", "k8s", "envoy", "istio", "linkerd",
		"prometheus", "opentelemetry", "otel", "ebpf", "wasm",
		"argo", "flux", "helm", "cilium", "kyverno", "trivy",
		"sigstore", "telemetry", "observability", "supply chain",
		"slsa", "containerd", "cri-o",
	}
	cfg.Keywords.Negative = []string{"marketing webinar", "press release", "sponsored"}

	// Sources
	cfg.Sources = []struct {
		Name   string
		Type   string
		URL    string
		Weight float64
		Tags   []string
	}{
		{"CNCF Blog", "rss", "https://www.cncf.io/feed/", 1.0, []string{"news", "foundation", "projects"}},
		{"Kubernetes Blog", "rss", "https://kubernetes.io/feed.xml", 1.0, []string{"kubernetes", "sig"}},
		{"OpenTelemetry Blog", "rss", "https://opentelemetry.io/blog/index.xml", 0.9, []string{"otel", "observability"}},
		{"Envoy Proxy", "rss", "https://blog.envoyproxy.io/feed", 0.8, []string{"envoy", "proxy"}},
		{"Cadence Workflow", "rss", "https://cadenceworkflow.io/blog/atom.xml", 0.7, []string{"automation", "configuration"}},
		{"DevStream", "rss", "https://blog.devstream.io/index.xml", 0.7, []string{"devstream"}},
		{"Kairos Blog", "rss", "https://kairos.io/blog/index.xml", 0.75, []string{"kairos", "edge", "linux"}},
		{"KCL Blog", "rss", "https://www.kcl-lang.io/blog/rss.xml", 0.72, []string{"kcl", "policy", "language"}},
		{"Kratix Blog", "rss", "https://docs.kratix.io/blog/rss.xml", 0.73, []string{"kratix", "platform"}},
		{"KubeDL Blog", "rss", "https://kubedl.io/blog/rss.xml", 0.74, []string{"kubedl", "ml", "kubernetes"}},
		{"KubeEdge Blog", "rss", "https://kubeedge.io/blog/rss.xml", 0.82, []string{"kubeedge", "edge", "iot"}},
		{"Konstruct Blog", "rss", "https://blog.konstruct.io/rss/", 0.71, []string{"konstruct", "infrastructure"}},
		{"KusionStack Blog", "rss", "https://medium.com/feed/@kusionstack", 0.72, []string{"kusionstack", "iac"}},
		{"ManageIQ Blog", "rss", "https://www.manageiq.org/feed.xml", 0.77, []string{"manageiq", "automation"}},
		{"OpenStack Blog", "rss", "https://www.openstack.org/blog/feed/", 0.92, []string{"openstack", "cloud"}},
		{"OpenTofu Blog", "rss", "https://opentofu.org/blog/rss.xml", 0.86, []string{"opentofu", "terraform", "iac"}},
		{"OpenYurt Blog", "rss", "https://openyurt.io/blog/rss.xml", 0.78, []string{"openyurt", "edge"}},
		{"Pulumi Blog", "rss", "https://www.pulumi.com/blog/rss.xml", 0.91, []string{"pulumi", "iac", "cloud"}},
		{"Salt Project Blog", "rss", "https://saltproject.io/blog/index.xml", 0.87, []string{"saltstack", "automation"}},
		{"Shifu Technical Blogs", "rss", "https://shifu.dev/technical-blogs/rss.xml", 0.71, []string{"shifu", "edge"}},
		{"Shifu News", "rss", "https://shifu.dev/news/rss.xml", 0.71, []string{"shifu", "news"}},
		{"HashiCorp Blog", "rss", "https://www.hashicorp.com/blog/feed.xml", 0.98, []string{"hashicorp", "terraform", "vault"}},
		{"Updatecli Blog", "rss", "https://www.updatecli.io/blog/index.xml", 0.72, []string{"updatecli", "automation"}},
		{"Dragonfly (D7y) Blog", "rss", "https://d7y.io/blog/rss.xml", 0.76, []string{"dragonfly", "distribution"}},
		{"Cerbos Blog", "rss", "https://www.cerbos.dev/rss", 0.74, []string{"cerbos", "authorization"}},
		{"Confidential Containers Blog", "rss", "https://confidentialcontainers.org/blog/index.xml", 0.77, []string{"confidential-containers", "security"}},
		{"Falco Blog", "rss", "https://falco.org/blog/feed.xml", 0.88, []string{"falco", "security"}},
		{"in-toto Blog", "rss", "https://in-toto.io/blog/index.xml", 0.76, []string{"in-toto", "supply-chain"}},
		{"Keycloak Blog", "rss", "https://www.keycloak.org/rss.xml", 0.89, []string{"keycloak", "identity"}},
		{"KubeArmor Blog", "rss", "https://kubearmor.io/blog/rss.xml", 0.77, []string{"kubearmor", "security"}},
		{"Kubescape New Feed", "rss", "https://kubescape.io/feed_rss_created.xml", 0.82, []string{"kubescape", "security"}},
		{"Kubescape Updates", "rss", "https://kubescape.io/feed_rss_updated.xml", 0.82, []string{"kubescape", "security"}},
		{"Kubewarden Blog", "rss", "https://www.kubewarden.io/blog/index.xml", 0.8, []string{"kubewarden", "policy"}},
		{"Kyverno Blog", "rss", "https://kyverno.io/blog/index.xml", 0.87, []string{"kyverno", "policy"}},
		{"AccuKnox Blog", "rss", "https://accuknox.com/feed", 0.76, []string{"accuknox", "security"}},
		{"Open Policy Containers Blog", "rss", "https://openpolicycontainers.com/blog/rss.xml", 0.72, []string{"policy", "containers"}},
		{"OpenFGA Blog", "rss", "https://openfga.dev/blog/rss.xml", 0.73, []string{"openfga", "authorization"}},
		{"Paladin Cloud Blog", "rss", "https://paladincloud.io/feed/", 0.74, []string{"paladincloud", "security"}},
		{"Paralus Blog", "rss", "https://www.paralus.io/blog/rss.xml", 0.73, []string{"paralus", "kubernetes", "networking"}},
		{"Ratify Blog", "rss", "https://ratify.dev/blog/rss.xml", 0.7, []string{"ratify", "supply-chain", "security"}},
		{"SecureCodeBox Blog", "rss", "https://www.securecodebox.io/blog/rss.xml", 0.7, []string{"security", "testing", "kubernetes"}},
		{"Sigstore Blog", "rss", "https://blog.sigstore.dev/index.xml", 0.9, []string{"sigstore", "security", "supply-chain"}},
		{"Sonobuoy Blog", "rss", "https://sonobuoy.io/blog/feed.xml", 0.9, []string{"sonobuoy", "testing", "kubernetes"}},
		{"StackRox Blog", "rss", "https://www.stackrox.io/rss.xml", 0.9, []string{"security", "containers", "kubernetes"}},
		{"Teleport Blog", "rss", "https://goteleport.com/blog/rss.xml", 0.8, []string{"teleport", "security", "identity"}},
		{"Varmor Blog", "rss", "https://www.varmor.org/blog/rss.xml", 0.7, []string{"varmor", "security", "kubernetes"}},
		{"Pinniped Blog", "rss", "https://pinniped.dev/blog/index.xml", 0.7, []string{"pinniped", "security", "identity"}},
		{"JuiceFS Blog", "rss", "https://juicefs.com/en/blog/latest/feed/", 0.8, []string{"storage", "distributed-systems", "filesystem"}},
		{"MinIO Blog", "rss", "https://blog.min.io/rss/", 0.9, []string{"minio", "storage", "s3"}},
		{"MooseFS Blog", "rss", "https://moosefs.com/feed/", 0.7, []string{"storage", "filesystem"}},
		{"ORAS Blog", "rss", "https://oras.land/blog/rss.xml", 0.7, []string{"oras", "artifacts", "registry"}},
		{"Rook Blog", "rss", "https://blog.rook.io/feed", 0.8, []string{"rook", "storage", "ceph"}},
		{"Velero Blog", "rss", "https://velero.io/blog/index.xml", 0.8, []string{"backup", "recovery", "kubernetes"}},
		{"CNCF Vineyard Blog", "rss", "https://medium.com/feed/cncf-vineyard", 0.7, []string{"vineyard", "data", "cncf"}},
		{"Layer5 Blog", "rss", "https://layer5.io/blog/feed.xml", 0.7, []string{"service-mesh", "istio", "kubernetes"}},
		{"Layer5 News", "rss", "https://layer5.io/news/feed.xml", 0.7, []string{"service-mesh", "community", "cncf"}},
		{"DeisLabs Blog", "rss", "https://deislabs.io/posts/index.xml", 0.7, []string{"deislabs", "cloud-native", "experiments"}},
		{"Kuasar Blog", "rss", "https://kuasar.io/blog/index.xml", 0.7, []string{"wasm", "containers", "runtime"}},
		{"Podman Desktop Blog", "rss", "https://podman-desktop.io/blog/rss.xml", 0.8, []string{"podman", "desktop", "containers"}},
		{"Podman Blog", "rss", "https://blog.podman.io/feed/", 0.8, []string{"podman", "containers", "linux"}},
		{"Antrea Blog", "rss", "https://antrea.io/blog/feed.xml", 0.7, []string{"antrea", "networking", "kubernetes"}},
		{"Cilium Blog", "rss", "https://cilium.io/blog/rss.xml", 0.9, []string{"cilium", "ebpf", "networking"}},
		{"Kube-OVN Blog", "rss", "https://www.kube-ovn.io/news/rss.xml", 0.7, []string{"networking", "ovn", "kubernetes"}},
		{"Red Hat Blog", "rss", "https://www.redhat.com/en/rss/blog", 0.9, []string{"redhat", "open-source", "enterprise"}},
		{"Crossplane Blog", "rss", "https://blog.crossplane.io/rss/", 0.8, []string{"crossplane", "infrastructure", "control-plane"}},
		{"Fluid Blog", "rss", "https://fluid-cloudnative.github.io/blog/rss.xml", 0.7, []string{"fluid", "data", "kubernetes"}},
		{"Project Hami Blog", "rss", "https://project-hami.io/blog/rss.xml", 0.7, []string{"hami", "gpu", "scheduling"}},
		{"Karmada Blog", "rss", "https://karmada.io/blog/rss.xml", 0.8, []string{"karmada", "multi-cluster", "kubernetes"}},
		{"KCP Blog", "rss", "https://www.kcp.io/blog/index.xml", 0.7, []string{"kcp", "multi-tenancy", "kubernetes"}},
		{"k0s Project Blog", "rss", "https://medium.com/feed/k0sproject", 0.7, []string{"k0s", "lightweight", "kubernetes"}},
		{"KEDA Blog", "rss", "https://keda.sh/blog/index.xml", 0.8, []string{"keda", "autoscaling", "kubernetes"}},
		{"Kestra Blog", "rss", "https://kestra.io/rss.xml", 0.7, []string{"kestra", "orchestration", "workflow"}},
		{"Knative Blog (Created)", "rss", "https://knative.dev/blog/feed_rss_created.xml", 0.9, []string{"knative", "serverless", "kubernetes"}},
		{"Knative Blog (Updated)", "rss", "https://knative.dev/blog/feed_rss_updated.xml", 0.9, []string{"knative", "serverless", "kubernetes"}},
		{"Koordinator Blog", "rss", "https://koordinator.sh/blog/rss.xml", 0.7, []string{"koordinator", "scheduling", "performance"}},
		{"Kube-Green Blog", "rss", "https://kube-green.dev/blog/rss.xml", 0.7, []string{"kube-green", "sustainability", "kubernetes"}},
	}

	// DB path
	cfg.DBPath = os.Getenv("DB_PATH")

	return cfg

}
