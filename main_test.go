package main

import (
	"log"
	"testing"
)

func TestGetCurrentIP(t *testing.T) {
	ip, err := GetCurrentIP()
	if err != nil {
		t.Fail()
	}
	log.Print(ip)
}

func TestGetRecordIP(t *testing.T) {
	id, err := GetZoneID("h.2k0ri.org")
	if err != nil {
		t.Fail()
	}
	ip, err := GetRecordIP("h.2k0ri.org", id)
	if err != nil {
		t.Fail()
	}
	log.Print(ip)
}

func TestGetZoneID(t *testing.T) {
	id, err := GetZoneID("h.2k0ri.org")
	if err != nil {
		t.Fail()
	}
	log.Print(id)
}
