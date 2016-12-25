package main

import (
	"log"
	"testing"
)

func TestGetCurrentIp(t *testing.T) {
	ip, err := GetCurrentIp()
	if err != nil {
		t.Fail()
	}
	log.Print(ip)
}

func TestGetRecordIp(t *testing.T) {
	id, err := GetZoneId("h.2k0ri.org")
	if err != nil {
		t.Fail()
	}
	ip, err := GetRecordIp("h.2k0ri.org", id)
	if err != nil {
		t.Fail()
	}
	log.Print(ip)
}

func TestGetZoneId(t *testing.T) {
	id, err := GetZoneId("h.2k0ri.org")
	if err != nil {
		t.Fail()
	}
	log.Print(id)
}
