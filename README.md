# url-shortener

## The API
#### POST /shortenedURL
Creates a new shortened URL.
The required JSON body requires the url to be shortened and has an optional field for an expiration date in the form of a unix timestamp. 
```
example body:
{
  "url": "https://twitter.com",
  "expirationDate": "1713830400"
}
```
It returns a JSON body with the original URL, the shortened URL, and the expiration date.
```
example response:
{
  "longURL": "https://twitter.com",
  "shortenedURL": "https://cloudflare-url-ownx73g3lq-uw.a.run.app/s/kAUM9w",
  "expirationDate": 1713830400
}
```

#### GET /s/{shortenedURL}
The new shortened URL redirects to the original url, as long as it is not expired.

#### GET /s/{shortenedURL}/usage
Returns a usage object that includes the number of times the urls has been accessed in the last day, week, and overall.
```
example response:
{
  "day": 17,
  "week": 122,
  "allTime": 11000
}
```

#### DELETE /s/{shortened_url}
Deletes the shortenedURL from the url table along with all of the associated entries in the usage table.
Returns a 204.

*There is room for confusion as I am returning a customer readable "shortenedURL" as a response to the post request. However, throughout the code, I use "shortenedURL" as an ID. There is room for improvement on the naming here and something I would look to product or UX for help fixing.


## Infrastructure
This app is deployed using Google Cloud Run.
It is connected to a Google Cloud SQL postgresql database. This database consists of two tables: urls and usage.

urls
```
     Column      |       Type        | Collation | Nullable | Default 
-----------------+-------------------+-----------+----------+---------
 long_url        | character varying |           | not null | 
 shortened_url   | character varying |           | not null | 
 expiration_date | integer           |           |          | 
Indexes:
    "urls_pkey" PRIMARY KEY, btree (shortened_url)
Referenced by:
    TABLE "usage" CONSTRAINT "usage_shortened_url_fkey1" FOREIGN KEY (shortened_url) REFERENCES urls(shortened_url) ON DELETE CASCADE
```

usage
```
    Column     |       Type        | Collation | Nullable | Default 
---------------+-------------------+-----------+----------+---------
 id            | character varying |           | not null | 
 shortened_url | character varying |           | not null | 
 usage_time    | integer           |           | not null | 
Indexes:
    "usage_pkey" PRIMARY KEY, btree (id)
Foreign-key constraints:
    "usage_shortened_url_fkey1" FOREIGN KEY (shortened_url) REFERENCES urls(shortened_url) ON DELETE CASCADE
```

I created the tables and columns manually to save time. 
I also deployed manually using the gcloud CLI.

## Next Steps
Here is a non-exhaustive list of next steps I would take for this project.
1. Make customer facing response and documentation with the shortenedURL less ambiguous (see above).
2. Create the database in terraform to avoid accidental changes.
3. Add testing. All of the methods in the controllers need unit tests, especially around checking expiration dates and calculating the usage stats. There are probably some bugs in there. There should be integration tests around the database methods. I would also add happy path post deployment smoke tests to make sure nothing major broke after a deployment (like breaking the database connection). And testing should run in a pipeline before deployment.
4. Fix error handling. The API currently returns errors that aren't helpful in indicating the problem or if it expected.
5. Add alerts and anomoly detection around increased error messages or decreased success rates. It also needs alerts around the database connections, throughput, and latency, at least.
6. Abstract out the database methods. Right now they are pretty brittle and don't allow for easily adding or modifying queries.
7. Create a pipeline to easily test and deploy code.
8. Rate limiting to prevent customers from destabalizing the application.
9. Determine with product or UX if expired shortened urls should still have usage statistics. If not, add logic to have that route return an appropriate response.

And I have many more ideas. I'm excited to talk to you about this. I really enjoyed working on it and flexing my problem solving muscle. I used this opportunity to learn more about GCP and relearned a lot about postgresql.
Thank you for your time.
