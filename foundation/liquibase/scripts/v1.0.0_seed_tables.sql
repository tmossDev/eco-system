DO $$
DECLARE
    AdminRoleId BIGINT;
    SystemAutoRoleId BIGINT;
    ModeratorRoleId BIGINT;
    UserRoleId BIGINT;
    UserSystemAutoId BIGINT;
BEGIN
    -- Insert roles and get their IDs
    INSERT INTO roles("name")
    VALUES('ROLE_ADMIN')
    ON CONFLICT ("name") DO UPDATE SET "name" = EXCLUDED."name"
    RETURNING id INTO AdminRoleId;

    INSERT INTO roles("name")
    VALUES('ROLE_SYSTEM_AUTO')
    ON CONFLICT ("name") DO UPDATE SET "name" = EXCLUDED."name"
    RETURNING id INTO SystemAutoRoleId;

    INSERT INTO roles("name")
    VALUES('ROLE_MODERATOR')
    ON CONFLICT ("name") DO UPDATE SET "name" = EXCLUDED."name"
    RETURNING id INTO ModeratorRoleId;

    INSERT INTO roles("name")
    VALUES('ROLE_USER')
    ON CONFLICT ("name") DO UPDATE SET "name" = EXCLUDED."name"
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
        'system@test.com',
        'System',
        'Auto',
        '$2a$10$up93efLWCZXf6iBGq.7XeerwJp3HzzaQUMBIIKzZN.rUtOIVYerbi',
        '0605281163',
        'System Auto'
    )
    ON CONFLICT ("account_code") DO UPDATE
    SET
        "email" = EXCLUDED."email",
        "first_name" = EXCLUDED."first_name",
        "last_name" = EXCLUDED."last_name",
        "password" = EXCLUDED."password",
        "phone" = EXCLUDED."phone",
        "username" = EXCLUDED."username",
        "updated_at" = now()
    RETURNING id INTO UserSystemAutoId;

    -- Insert user roles
    INSERT INTO user_roles(
        "user_id",
        "role_id"
    )
    VALUES(
        UserSystemAutoId,
        SystemAutoRoleId
    )
    ON CONFLICT ("user_id", "role_id") DO NOTHING;
END $$;
