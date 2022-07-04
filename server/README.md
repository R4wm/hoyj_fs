#
- create script to make object storage full path to download

- download all files from object storage to new mounted drive @ /var/www/html/helpersofyourjoy/storage/media

- create labels for all files and place them in /var/www/html/helpersofyourjoy/storage/labels

- cron job to scrape all labels and push into redis

- api to get that payload from redis
  - api should take args to set filters and get new json payload


---
### Maintenance
- cron job to check all labels correspond to correct files
