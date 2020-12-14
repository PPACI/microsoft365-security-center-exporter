package exporter

import (
	"context"
	"github.com/PPACI/microsoft-defender-ATP-exporter/pkg/api/security"
	"github.com/PPACI/microsoft-defender-ATP-exporter/pkg/azureauth"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"time"
)

var (
	secureScoreCurrentGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "microsoft365_security",
		Name:      "secure_score_current",
		Help:      "Secure score",
	}, []string{"tenant_id"})
	secureScoreMaxGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "microsoft365_security",
		Name:      "secure_score_max",
		Help:      "Max Secure score",
	}, []string{"tenant_id"})
	activeUserCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "microsoft365_security",
		Name:      "active_user_count",
		Help:      "Number of active user",
	}, []string{"tenant_id"})
	licensedUserCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "microsoft365_security",
		Name:      "licensed_user_count",
		Help:      "Number of licensed user",
	}, []string{"tenant_id"})
)

func init() {
	prometheus.MustRegister(secureScoreCurrentGauge)
	prometheus.MustRegister(secureScoreMaxGauge)
	prometheus.MustRegister(activeUserCount)
	prometheus.MustRegister(licensedUserCount)
}

func refreshData(authClient *azureauth.AuthClient) {
	refreshSecureScore(authClient)
	log.Println("Refreshed data from API")
}

func refreshSecureScore(authClient *azureauth.AuthClient) {
	secureScore, err := security.GetSecureScore(authClient)
	if err != nil {
		log.Println(err)
	}
	secureScoreCurrentGauge.WithLabelValues(authClient.TenantId).Set(secureScore[0].CurrentScore)
	secureScoreMaxGauge.WithLabelValues(authClient.TenantId).Set(secureScore[0].MaxScore)
	activeUserCount.WithLabelValues(authClient.TenantId).Set(float64(secureScore[0].ActiveUserCount))
	licensedUserCount.WithLabelValues(authClient.TenantId).Set(float64(secureScore[0].LicensedUserCount))
}

func StartExporter(authClient *azureauth.AuthClient, ctx context.Context) {
	ticker := time.NewTicker(30 * time.Second)

	defer ticker.Stop()
	go func() {
		for {
			select {
			case <-ticker.C:
				refreshData(authClient)
			case <-ctx.Done():
				log.Println(ctx.Err())
				return
			}
		}
	}()

	// To start with fresh data
	go refreshData(authClient)
}
