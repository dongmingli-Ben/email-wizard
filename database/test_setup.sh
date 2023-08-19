touch ~/.pgpass
chmod 600 ~/.pgpass

echo "hostname:port:database:username:password" >> ~/.pgpass
echo "postgres:5432:email-wizard-data:postgres:123456" >> ~/.pgpass