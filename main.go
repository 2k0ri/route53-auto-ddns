package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

var (
	creds = credentials.NewChainCredentials([]credentials.Provider{&credentials.EnvProvider{}, &credentials.SharedCredentialsProvider{}})
	sess  = session.Must(session.NewSessionWithOptions(session.Options{
		Config:            aws.Config{Credentials: creds},
		SharedConfigState: session.SharedConfigEnable,
	}))
	client = route53.New(sess)
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal(fmt.Sprintf("%s [domain]", os.Args[0]))
	}
	domain := os.Args[1]

	id, err := GetZoneId(domain)
	if err != nil {
		log.Fatal(err)
	}
	rip, err := GetRecordIp(domain, id)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("record's ip: " + rip)
	cip, err := GetCurrentIp()
	if err != nil {
		log.Fatal(err)
	}
	log.Print("current ip: " + cip)
	if rip == cip {
		log.Print("ip not changed")
		return
	}
	err = SetRecord(domain, id, cip)
	if err != nil {
		log.Fatal(err)
	}
	log.Print(fmt.Sprintf("ip changed: %s -> %s", rip, cip))
}

func GetCurrentIp() (string, error) {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	s, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	ip := strings.TrimRight(string(s), "\n")
	return ip, nil
}

func SetRecord(record, zoneid, ip string) error {
	_, err := client.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneid),
		ChangeBatch: &route53.ChangeBatch{Changes: []*route53.Change{
			{
				Action: aws.String(route53.ChangeActionUpsert),
				ResourceRecordSet: &route53.ResourceRecordSet{
					Name: aws.String(record),
					Type: aws.String(route53.RRTypeA),
					TTL:  aws.Int64(3600),
					ResourceRecords: []*route53.ResourceRecord{
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

func GetRecordIp(record, zoneid string) (string, error) {
	out, err := client.ListResourceRecordSets(&route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(zoneid),
		StartRecordName: aws.String(record),
		StartRecordType: aws.String(route53.RRTypeA),
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

func GetZoneId(record string) (string, error) {
	_d := strings.Split(record, ".")
	d := strings.Join(_d[len(_d)-2:], ".") // truncate subdomain
	out, err := client.ListHostedZones(&route53.ListHostedZonesInput{})
	if err != nil {
		return "", err
	}
	for i := range out.HostedZones {
		if strings.Contains(d+".", *out.HostedZones[i].Name) {
			return *out.HostedZones[i].Id, nil
		}
	}
	return "", errors.New(fmt.Sprintf(`No zone matched with "%s"`, d))
}
