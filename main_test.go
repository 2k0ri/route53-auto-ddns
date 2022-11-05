package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"log"
	"os"
	"testing"
)

var (
	ctx    context.Context
	client *route53.Client
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}
	client = route53.NewFromConfig(cfg)
	os.Exit(m.Run())
}

func TestGetCurrentIP(t *testing.T) {
	ip, err := GetCurrentIP()
	if err != nil {
		t.Fail()
	}
	log.Print(ip)
}

func TestGetRecordIP(t *testing.T) {
	id, err := GetZoneID(ctx, client, "h.2k0ri.org")
	if err != nil {
		t.Fail()
	}
	ip, err := GetRecordIP(ctx, client, "h.2k0ri.org", id)
	if err != nil {
		t.Fail()
	}
	log.Print(ip)
}

func TestGetZoneID(t *testing.T) {
	id, err := GetZoneID(ctx, client, "h.2k0ri.org")
	if err != nil {
		t.Fail()
	}
	log.Print(id)
}
