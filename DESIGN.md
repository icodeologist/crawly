# Web crawler v1
### Design architecture
- Starts with single crawler
- Fetches all the links (a[href])
- Checks duplicates with a in memory map - map[url]exist? true or false
- Then adds to queue which handles normalization of urls and proceeds crawling new links

### Flaws
- Each time only single link is crawling - Takes too long
- Map looses data while restarting
- Politeness violation -  [politeness policy](https://en.wikipedia.org/wiki/Web_crawler)
- Poor error handling - No keeping track of failed links

### tech stack
- htmlquery(parsing)
- net/url && urlx(normalization)
- queue go package
