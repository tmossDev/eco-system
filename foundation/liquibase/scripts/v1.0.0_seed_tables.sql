DO $$
DECLARE
    AdminRoleId INTEGER;
    SystemAutoRoleId INTEGER;
    ModeratorRoleId INTEGER;
    UserRoleId INTEGER;
    UserSystemAutoId INTEGER;
BEGIN
    -- Insert roles and get their IDs
    INSERT INTO roles("name")
    VALUES('ROLE_ADMIN')
    RETURNING id INTO AdminRoleId;

    INSERT INTO roles("name")
    VALUES('ROLE_SYSTEM_AUTO')
    RETURNING id INTO SystemAutoRoleId;

    INSERT INTO roles("name")
    VALUES('ROLE_MODERATOR')
    RETURNING id INTO ModeratorRoleId;

    INSERT INTO roles("name")
    VALUES('ROLE_USER')
    RETURNING id INTO UserRoleId;

    -- Insert user and get the user ID
    INSERT INTO "users" (
        "account_code",
        "email",
        "first_name",
        "last_name",
        "password",
        "phone",
        "username"
    )
    VALUES (
        'S001',
        'systen@test.com',
        'System',
        'Auto',
        '$2a$10$up93efLWCZXf6iBGq.7XeerwJp3HzzaQUMBIIKzZN.rUtOIVYerbi',
        '0605281163',
        'System Auto'
    )
    RETURNING id INTO UserSystemAutoId;

    -- Insert user roles
    INSERT INTO user_roles(
        "user_id",
        "role_id"
    )
    VALUES(
        UserSystemAutoId,
        SystemAutoRoleId
    );