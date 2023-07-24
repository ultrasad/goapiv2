db.createUser(
    {
        user: "kajame",
        pwd: "111111",
        roles: [
            {
                role: "readWrite",
                db: "example"
            }
        ]
    }
);
db.createCollection("test"); //MongoDB creates the database when you first store data in that database