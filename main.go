package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {

	lastMonth := time.Now().AddDate(0, -1, 0)

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

	// Snapshot describe input
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
	for idx, item := range resp.Snapshots {
		if checkSnapTime(*item.StartTime, lastMonth) {
			fmt.Printf("[%d] id %s  created on %s. Scheduling for deletion\n", idx, *item.SnapshotId, *item.StartTime)
			// Delete the snapshot
			params := &ec2.DeleteSnapshotInput{
				DryRun:     aws.Bool(false),
				SnapshotId: item.SnapshotId,
			}
			_, err := svc.DeleteSnapshot(params)
			if err != nil {
				fmt.Println(err.Error())
				return
			}

		} else {
			fmt.Printf("[%d] id %s  created on %s Keeping\n", idx, *item.SnapshotId, *item.StartTime)
		}
	}
}

func checkSnapTime(startTime, checkTime time.Time) bool {
	return startTime.Before(checkTime)
}
