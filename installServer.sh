#!/bin/sh

clear
printf "VianuEdu-Server First-Time Installer v1.0\\n"
if [ "$(id -u)" != "0" ]; then
    printf "Not running as root! Please run as root to complete installation process!\\n"
    exit
fi
printf "Please write the installation path: "
read -r INSTALL_PATH
printf "An administrator user is required to access both the database and some handler functions of the server itself. Credentials are required.\\n"
printf "Please insert the username for the administrator account: "
read -r ADMIN_USER
printf "Please insert the password for the administrator account: "
read -r ADMIN_PASS

if [ ! -d "$DIRECTORY" ]; then
    printf "Directory does not exist! Creating..."
    mkdir "$INSTALL_PATH"
    printf "DONE\\n"
fi

printf "Creating VianuEdu-Server environment..."
mkdir "$INSTALL_PATH/config"

printf "{\\n  \"serverIP\": \"[DATABASE SERVER IP HERE]\",\\n  \"serverPort\": \"[DATABASE SERVER PORT HERE]\",\\n  \"userName\": \"VianuEdu_DataAdmin\",\\n  \"userPass\": \"%s\",\\n  \"databaseName\": \"VianuEdu\"\\n}" "$ADMIN_PASS" > "$INSTALL_PATH/config/DatabaseSettings.json"
printf "{\\n  \"listenPort\": [HTTP SERVER PORT HERE],\\n  \"adminUser\": \"%s\",\\n  \"adminPass\": \"%s\",\\n  \"enableTLS\": false,\\n  \"certFile\": \"[CERTIFICATE FILE PATH HERE]\",\\n  \"keyFile\": \"[PRIVATE KEY FILE PATH HERE]\"\\n}" "$ADMIN_USER" "$ADMIN_PASS" > "$INSTALL_PATH/config/HTTPServer.json"

mkdir "$INSTALL_PATH/errors"
mkdir "$INSTALL_PATH/keys"
mkdir "$INSTALL_PATH/lessons"
mkdir "$INSTALL_PATH/lessons/Geo"
mkdir "$INSTALL_PATH/lessons/Phi"
mkdir "$INSTALL_PATH/lessons/Math"
mkdir "$INSTALL_PATH/lessons/Info"
nrOfGrades=12
for i in $(seq 1 $nrOfGrades)
do
    mkdir "$INSTALL_PATH/lessons/Geo/$i"
    mkdir "$INSTALL_PATH/lessons/Phi/$i"
    mkdir "$INSTALL_PATH/lessons/Math/$i"
    mkdir "$INSTALL_PATH/lessons/Info/$i"
done
mkdir "$INSTALL_PATH/log"
mkdir "$INSTALL_PATH/static"
mkdir "$INSTALL_PATH/templates"
printf "DONE\\n"
printf "Checking if MongoDB Server and Shell are installed..."

if ! service mongod status; then
    printf "\\n [WARN] MongoDB Server and Shell not installed! Installing..."
    apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv 2930ADAE8CAF5059EE73BB4B58712A2291FA4AD5
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu xenial/mongodb-org/3.6 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-3.6.list
    apt-get update
    apt-get install -y mongodb-org
    printf "\\nDone installing MongoDB Server and Shell!\\n"
else
    printf "DONE\\n"
fi

printf "Commencing database initialization...\\n"
printf "Starting MongoDB instance..."
service mongod start
service mongod enable
printf "DONE\\n"
printf "Creating schema necessary for VianuEdu..."
printf "use admin\\n db.createUser({user: \"admin\", pwd:\"%s\", roles:[{role: \"root\", db:\"admin\"}]})\\n use VianuEdu\\n db.createCollection(\"VianuEdu.TestList\")\\n db.createCollection(\"Students.Accounts\")\\n db.createCollection(\"Teachers.Accounts\")\\n db.createCollection(\"Students.SubmittedAnswers\")\\n db.createCollection(\"GeoEdu.Tests\")\\n db.createCollection(\"PhiEdu.Tests\")\\n db.createCollection(\"MathEdu.Tests\")\\n db.createCollection(\"InfoEdu.Tests\")\\n db.createCollection(\"GeoEdu.Grades\")\\n db.createCollection(\"PhiEdu.Grades\")\\n db.createCollection(\"MathEdu.Grades\")\\n db.createCollection(\"InfoEdu.Grades\")\\n db.createUser({user: \"VianuEdu_DataAdmin\", pwd:\"%s\", roles:[{role: \"root\", db:\"VianuEdu\"}]})" "$ADMIN_PASS" "$ADMIN_PASS" > MongoDBScript.json
mongo < MongoDBScript.json
printf "DONE\\n"
printf "Changing database configuration to use authorization on boot..."
awk 'NR==9 {$0="ExecStart=/usr/bin/mongod --auth --config /etc/mongod.conf"}1' /lib/systemd/system/mongod.service
systemctl daemon-reload
service mongod restart
printf "DONE\\n"
printf "Database initialization complete! Installing VianuEdu-Server..."
wget https://github.com/CNITV/VianuEdu-Server/releases/download/v1.0/VianuEdu-Server-v1.0 > /dev/null
mv VianuEdu-Server-v1.0 /usr/bin/VianuEdu-Server
printf "DONE\\n"
printf "INSTALLATION COMPLETE! As of now, you can start the server at any time by calling the following command from the command line: \\n"
printf "    root@localhost#/path/to/install/path$ VianuEdu-Server &\\n"
printf "Make sure that the database is on by starting it and that the command is called from the install directory.\\n"
printf "In addition, make sure that all the configuration files are properly set before starting the server, otherwise it will fail.\\n"