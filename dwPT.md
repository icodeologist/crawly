How crawler works:
  Start with seed url
  then GET req that url
  Get the html data and parse it 
  Gather the data 
  Put it to database
  Fetch other links from the page
  add those links to queue
  Then recheck if they have already visited 
  if already visited skip
  else crawl and repeat the same

Maintain max depth of 3 and max pages of 20 limit for now
