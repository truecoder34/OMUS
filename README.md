# OMUS
OMUS - One More URL Shortener;
For now - planned only API.

## Functions:
1) Generate UNIQUE aliases for provided URL; 
2) Should redirect users to original URL via short LINK;
3) Short Link should have a lifetime. Short Link expires after EOL;
4) Collect and Store statistics about short link: visits count, re-/generate count;

## Not Functions: 
1) Service should be async : handle numerous requests ; 
2) Real-time forwarding ;
3) Absolute random short links generation ; 

## API Endpoints: 
1) [POST] create/ - create new short link and put it into DB;
2) [GET] - check if short link already exists for specified URL;
3) [GET] - get statistics about short link : EOL date, redirects, calls;
4) [GET] - redirect to original URL by passing a short one;

## Tech Stack
1) Initially will be Postgress USED
2) Try to use Redis instead
3) RabbitMQ msg-broker to log events of API

