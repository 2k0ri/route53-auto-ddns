package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	route53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("%s [domain]", os.Args[0]))
	}
	domain := os.Args[1]

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	client := route53.NewFromConfig(cfg)

	id, err := GetZoneID(ctx, client, domain)
	if err != nil {
		log.Fatal(err)
	}
	rip, err := GetRecordIP(ctx, client, domain, id)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("record's ip: " + rip)
	cip, err := GetCurrentIP()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("current ip: " + cip)
	if rip == cip {
		log.Print("ip not changed")
		return
	}
	err = SetRecord(ctx, client, domain, id, cip)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(fmt.Sprintf("ip changed: %s -> %s", rip, cip))
}

func GetCurrentIP() (string, error) {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	s, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimRight(string(s), "\n")
	return ip, nil
}

func GetZoneID(ctx context.Context, client *route53.Client, record string) (string, error) {
	_d := strings.Split(record, ".")
	d := strings.Join(_d[len(_d)-2:], ".") // truncate subdomain
	out, err := client.ListHostedZones(ctx, &route53.ListHostedZonesInput{})
	if err != nil {
		return "", err
	}
	for i := range out.HostedZones {
		if strings.Contains(d+".", *out.HostedZones[i].Name) {
			return *out.HostedZones[i].Id, nil
		}
	}
	return "", fmt.Errorf(`No zone matched with "%s"`, d)
}

func SetRecord(ctx context.Context, client *route53.Client, record, zoneid, ip string) error {
	_, err := client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneid),
		ChangeBatch: &route53types.ChangeBatch{Changes: []route53types.Change{
			{
				Action: route53types.ChangeActionUpsert,
				ResourceRecordSet: &route53types.ResourceRecordSet{
					Name: aws.String(record),
					Type: route53types.RRTypeA,
					TTL:  aws.Int64(3600),
					ResourceRecords: []route53types.ResourceRecord{
						{Value: aws.String(ip)},
					},
				},
			},
		}},
	})
	if err != nil {
		return err
	}
	return nil
}

func GetRecordIP(ctx context.Context, client *route53.Client, record, zoneid string) (string, error) {
	out, err := client.ListResourceRecordSets(ctx, &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneid),
		StartRecordName: aws.String(record),
		StartRecordType: route53types.RRTypeA,
	})
	if err != nil {
		return "", err
	}
	if len(out.ResourceRecordSets) == 0 {
		log.Print(fmt.Sprintf(`No record matched with "%s"`, record))
		return "", nil
	}
	ip := *out.ResourceRecordSets[0].ResourceRecords[0].Value
	return ip, nil
}
