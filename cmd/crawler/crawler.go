package main

import "leetcode-submission-crawler/logger"

/*
There will be two way to trigger for an update:
  1. by cron expression: we will provide a cron expression to trigger the update job
  2. by rest api: we will provide a rest api to trigger the update job

  The crawler app will receive a channel, that contains the command to trigger the update job.
*/

func main() {
	logger.InitLogger()

	defer logger.Sync()
}
