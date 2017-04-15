package main

import (
	"flag"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	// flags
	volumeID := flag.String("volumeid", "volume id", "Volume Id")

	flag.Parse()

	// Create an EC2 Service
	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "default",
	})

	if err != nil {
		fmt.Println("Error creating session", err)
		return
	}

	svc := ec2.New(sess)

	params := &ec2.DescribeSnapshotsInput{
		DryRun: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("volume-id"),
				Values: []*string{
					aws.String(*volumeID),
				},
			},
		},
	}

	resp, err := svc.DescribeSnapshots(params)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// print the response
	fmt.Println(resp)
}
