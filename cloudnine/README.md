# cloudnine

What is this?
-------------

Cloudnine is our sourcecode repository for our submission to govhack 2015. As team 'bigdatabigdreams', our aim was to build a system which could help you discover suburbs that were very similar to what you were familiar with and felt like home - but never knew about. We compare culture, crime data, demographics, and a host of other information so that you can be presented with the perfect match.

Setup
-------------

1. Install a postgresql database on a machine, with a user _cloudnine_ and a database _cloudnine_.
2. Run the generateData.sh script. This will read the dataset files and generate all the database tables and linkages. It will take about 20 minutes.
3. As root, run the site.py script using python. The server should now be running.
