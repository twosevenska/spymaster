db.auth('admin-user', 'adm1n')

spy = db.getSiblingDB("spymaster");

spy.createUser(
    {
        user: "spymaster",
        pwd: "face1t",
        roles: [
            {
                role: "readWrite",
                db: "spymaster",
            }
        ]
    }
);
spy.createCollection('users')

spy.users.insert({
    first_name: 'Yellow',
    last_name: 'King',
    nickname: 'hastur',
    password: 'Carcosa',
    email: 'hastur@lost.space',
    country: 'UK',
    created_at: '1',
    updated_at: '',
})

test = db.getSiblingDB("test_spymaster");

test.createUser(
    {
        user: "test",
        pwd: "test",
        roles: [
            {
                role: "readWrite",
                db: "test_spymaster",
            }
        ]
    }
);

test.createCollection('users')