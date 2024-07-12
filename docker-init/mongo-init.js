db = db.getSiblingDB('ggnet_db');

db.createUser({
    user: 'ggnet_user',
    pwd: '1q2w3e!',
    roles: [
        {
            role: 'readWrite',
            db: 'ggnet_db',
        },
    ],
});
