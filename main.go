package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	relativePath, _ := strconv.ParseBool(os.Getenv("IS_RELATIVE_PATH"))
	action := Action{
		Action:       os.Getenv("ACTION"),
		Bucket:       os.Getenv("BUCKET"),
		S3Class:      os.Getenv("S3_CLASS"),
		S3Prefix:     os.Getenv("S3_PREFIX"),
		Key:          fmt.Sprintf("%s.zip", os.Getenv("KEY")),
		RelativePath: relativePath,
		Artifacts:    strings.Split(strings.TrimSpace(os.Getenv("ARTIFACTS")), "\n"),
	}

	switch act := action.Action; act {
	case PutAction:
		if len(action.Artifacts[0]) <= 0 {
			log.Fatal("No artifacts patterns provided")
		}
		log.Print("putting files")

		if err := Zip(action.Key, action.Artifacts, action.RelativePath); err != nil {
			log.Fatal(err)
		}

		if err := PutObject(action.Key, fmt.Sprintf("%s%s", action.S3Prefix, action.Key), action.Bucket, action.S3Class); err != nil {
			log.Fatal(err)
		}
	case GetAction:
		exists, err := ObjectExists(fmt.Sprintf("%s%s", action.S3Prefix, action.Key), action.Bucket)
		if err != nil {
			log.Fatal(err)
		}

		// Get and and unzip if object exists
		if exists {
			if err := GetObject(action.Key, fmt.Sprintf("%s%s", action.S3Prefix, action.Key), action.Bucket); err != nil {
				log.Fatal(err)
			}

			if err := Unzip(action.Key, action.RelativePath); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Printf("No caches found for the following key: %s", action.Key)
		}
	case DeleteAction:
		if err := DeleteObject(fmt.Sprintf("%s%s", action.S3Prefix, action.Key), action.Bucket); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("Action \"%s\" is not allowed. Valid options are: [%s, %s, %s]", act, PutAction, DeleteAction, GetAction)
	}
}
