package main

import (
	"context"
	"github.com/PPACI/microsoft-defender-ATP-exporter/pkg/azureauth"
	"github.com/PPACI/microsoft-defender-ATP-exporter/pkg/exporter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
)

var (
	AzureTenantId     string
	AzureClientId     string
	AzureClientSecret string
)

func init() {
	var ok bool
	AzureTenantId, ok = os.LookupEnv("AZURE_TENANT_ID")
	if !ok {
		log.Fatalln("Missing AZURE_TENANT_ID variable")
	}
	AzureClientId, ok = os.LookupEnv("AZURE_CLIENT_ID")
	if !ok {
		log.Fatalln("Missing AZURE_CLIENT_ID variable")
	}
	AzureClientSecret, ok = os.LookupEnv("AZURE_CLIENT_SECRET")
	if !ok {
		log.Fatalln("Missing AZURE_CLIENT_SECRET variable")
	}
}

func main() {
	// Init Auth client
	log.Println("Init Azure Auth client")
	client := azureauth.NewAuthClient(AzureTenantId, AzureClientId, AzureClientSecret)

	// Start exporter polling job
	ctx := context.Background()
	exporter.StartExporter(client, ctx)

	// Expose prometheus handler
	log.Println("Listen on :8080")
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))

}
