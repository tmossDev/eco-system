# Project Change Log

## v1.5.0 - (6 Changes)
- Renamed/reworked Header Component to Navigation Component.
- Added Login Form Component to angular design system.
- Added Login, Forgot Password, User Management and Settings pages to user-management admin-web-app.
- Added Guard to user-management admin-web-app forcing users to login first.
- Added Interceptors for authentication, error handling and api-mocking.
- Added mocking for user-service and auth-service.

## v1.4.1 - (3 Changes)
- Added Icon component to angular design system.
- Added Header component to angular design system which uses Icon component.
- Using Header component in user-management admin-web-app.

## v1.4.0 - (5 Changes)
- Added pnpm for frontend package management.
- Added user-management web-app which is just hello world angular project
- Added angular design system for all angular web-apps in mono repo.
- Added button to angular design system.
- Added storybook for angular design system.

## v1.3.0 - (2 Changes)
- Added user-service controller for user management external calls.
- Added user-service routes for user management external calls using user-service controller.

## v1.2.0 - (4 Changes)
- Added user package for user management code reuse.
- Added user-gateway controller for user management external calls.
- Added user-gateway middleware for user management external calls.
- Added user-gateway routes for user management external calls using user-gateway controller.

## v1.1.0 - (4 Changes)
- Added shared-components for common components.
- Added GoLang common libraries to shared-components/backend/lib.
- Added user-service shell application for user management internal calls.
- Added user-gateway shell application for user management external calls.

## v1.0.0 - (4 Changes)
- Added foundation use case for getting local postgres database up and running.
- Added database Init Script for starting up the shop database from scratch
- Added liquibase setup and seed scripts for postgres database.
- Added a basic docker-compose file for foundation use case.