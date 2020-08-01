How it works:

The program takes a query parameter specified in the url,
checks that the query contains only letter characters, then
reads the shakespeare-complete.txt and searches throughout that 
text file for words that have the prefix and stores them in a 
map keeping track of their occurrences. 

After scanning the file is complete, the map is sorted by value
and the top 25 most frequent words are concatenated into a string
and returned as output. 

To prevent repeatedly calculating the same queries over and over again,
I included a cache to store queries that have already been completed before
and this will instantly return the result if it exists in the cache.

---------------------------------------------------------------------------------------------------------------------------------

How to use:

Assuming you have access to a unix machine, download the project and go to the
directory of the project on your local machine. 

Run these commands:
go build autocomplete.go
./autocomplete &

The second command will run the program as a process in the background and return a 
PID, it may be useful to remember this PID to kill it later in case you come
across issues while running the program or would like to free up the port. 

To query the program in the terminal please use:

curl http://localhost:9000/autocomplete?term=YourQueryHere

replace YourQueryHere with your query, for example:

curl http://localhost:9000/autocomplete?term=th

After a couple of seconds you should get a result like:

query: results 

for example:

th: the that this thou thy thee they then there their them than these th think thus though therefore those thine three thought thing things thousand

If you do not remember the PID you can look it up by running netstat -tulpn (assuming you
are root) or sudo netstat -tulpn. Here is a link for reference: https://askubuntu.com/questions/1217513/what-does-the-tulpn-option-mean-for-netstat

Then given the PID, to stop the process run:

kill PID


---------------------------------------------------------------------------------------------------------------------------------

Notes:

sample-completions.txt contains results for the prefixes: th, fr, pi, sh, wu, ar, il, ne, se, pl

I included a bash file called run.sh that you can run with:

bash run.sh

please execute this only if the program is running, it is a simple list of the curl requests for the above prefixes

Queries to the program will be logged in autocomplete.txt, please note that killing and restarting the process
and issuing new queries will create a new autocomplete.txt and overwrite the previous one.

In the source code I have included links to resources I referenced and used.

