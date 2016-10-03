service postgresql restart
python crimedata_importer.py
python abcdata_importer.py
python crimedata_geocoder.py
python abcdata_tagcreate.py
echo "DB generation complete."
