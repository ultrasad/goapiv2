mongo -- "$MONGO_INITDB_DATABASE" <<EOF
    var rootUser = '$MONGO_INITDB_ROOT_USERNAME';
    var rootPassword = '$MONGO_INITDB_ROOT_PASSWORD';
    var admin = db.getSiblingDB('admin');
    admin.auth(rootUser, rootPassword);
    var user = '$MONGO_INITDB_USERNAME';
    var passwd = '$MONGO_INITDB_PASSWORD';

    db.createUser({
      user: user,
      pwd: passwd,
      db: example,
      roles: ["readWrite"]
    });
    
    db.createCollection("test"); //MongoDB creates the database when you first store data in that database
EOF