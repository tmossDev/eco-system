DO $$
DECLARE
    AdminRoleId INTEGER;
    SystemAutoRoleId INTEGER;
    ModeratorRoleId INTEGER;
    UserRoleId INTEGER;
    UserSystemAutoId INTEGER;
    AdminUserId INTEGER;
    ModeratorUserId INTEGER;
    BlikkiesUserId INTEGER;
BEGIN
    -- Select role IDs
    SELECT id INTO AdminRoleId FROM roles WHERE "name" = 'ROLE_ADMIN';
    SELECT id INTO SystemAutoRoleId FROM roles WHERE "name" = 'ROLE_SYSTEM_AUTO';
    SELECT id INTO ModeratorRoleId FROM roles WHERE "name" = 'ROLE_MODERATOR';
    SELECT id INTO UserRoleId FROM roles WHERE "name" = 'ROLE_USER';

    SELECT id INTO SystemAutoRoleId FROM users WHERE "account_code" = 'S001';


    -- Insert additional users and get their IDs
    INSERT INTO "users" (
        "account_code",
        "email",
        "first_name",
        "last_name",
        "password",
        "phone",
        "username",
        "user_created_id",
        "user_updated_id"
    )
    VALUES (
        'A001',
        'admin@test.com',
        'Jane',
        'Doe',
        '$2a$10$up93efLWCZXf6iBGq.7XeerwJp3HzzaQUMBIIKzZN.rUtOIVYerbi',
        '0605281163',
        'Admin',
        UserSystemAutoId,
        UserSystemAutoId
    )
    RETURNING id INTO AdminUserId;

    INSERT INTO "users" (
        "account_code",
        "email",
        "first_name",
        "last_name",
        "password",
        "phone",
        "username",
        "user_created_id",
        "user_updated_id"
    )
    VALUES (
        'A002',
        'moderator@test.com',
        'Michael',
        'Rogers',
        '$2a$10$up93efLWCZXf6iBGq.7XeerwJp3HzzaQUMBIIKzZN.rUtOIVYerbi',
        '0605281162',
        'Moderator',
        UserSystemAutoId,
        UserSystemAutoId
    )
    RETURNING id INTO ModeratorUserId;

    INSERT INTO "users" (
        "account_code",
        "email",
        "first_name",
        "last_name",
        "password",
        "phone",
        "username",
        "user_created_id",
        "user_updated_id"
    )
    VALUES (
        'A003',
        'blikkies@test.com',
        'Blikkies',
        'Blignaut',
        '$2a$10$up93efLWCZXf6iBGq.7XeerwJp3HzzaQUMBIIKzZN.rUtOIVYerbi',
        '0605281161',
        'Blikkies',
        UserSystemAutoId,
        UserSystemAutoId
    )
    RETURNING id INTO BlikkiesUserId;
END $$;
