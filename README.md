# Nanny - Pet Sitting Service API

> A comprehensive REST API backend service for a pet sitting platform that was written fully with Go, and this platform links pet owners with pet sitters in Kazakhstan providing a full features booking/review system. The team tried their best to do this!

## Project Overview
The Nanny is a real world solution for pet care services nowadays, because there is no reliable system for this that enables:
1. Pet owners to find and book qualified sitters
2. Pet sitters to offer their services and manage bookings
3. Administrators to oversee platform operations and approve sitters
4. Secure authentication and role-based access control
5. Review and rating system for quality assurance



## Team Members

| Name | Role | Responsibilities |
| --- | --- | --- |
| Armankyzy Anara | Team Lead | Project coordination, authentication module development, input validation, repository organization, initial database schema design, code review and quality control, infrastructure management, environment setup |
| Alimzhankyzy Nuray | Scrum Master | Messaging system implementation, sprint planning, owner bookings fixes, team workflow coordination, JSON response standardization, unit test development, testing strategy and implementation |
| Kalikhan Arukhan | QA Engineer | Makefile configuration, test coverage setup, JWT authentication implementation, backend/frontend separation, code quality assurance |
| Sabukhi Raziyev | Core Backend Developer | Docker containerization, golang-migrate setup, database migrations, deployment configuration |


## Core Features
 User Management: Registration and authentication for owners, sitters, and admins
 Pet Profiles: Complete pet information management
 Service Listings: Sitters can create and manage service offerings
 Booking System: Full booking lifecycle (create, confirm, cancel, complete)
 Review System: Rating and feedback for completed services
 Admin Dashboard: User management and sitter approval workflow
## Technical Features
 JWT Authentication: Secure token-based authentication
 Role-Based Access Control: Owner, Sitter, and Admin roles
 Background Workers: Automated expired booking cleanup
 Graceful Shutdown: Proper context handling and shutdown
 Database Migrations: Schema versioning with golang-migrate
 Comprehensive Testing: Unit tests with 70%+ coverage
 Docker Support: Containerized deployment
 CORS Support: Cross-origin resource sharing

## Layer Architecture
Each module follows a three-layer architecture:
Handler Layer (HTTP)
Service Layer (Business Logic)
Repository Layer (Database)
Database schema
We used PostgreSQL with 9 interconnected tables:
## Core Entities
users - User accounts (owners, sitters, admins)
pets - Pet profiles owned by users
sitters - Extended sitter profiles with experience and location
services - Services offered by sitters
bookings - Service bookings between owners and sitters
reviews - Ratings and feedback for completed bookings
payments - Payment tracking (future feature)
chats - Chat sessions for bookings
messages - Messages within chats
# Entity Relationships
users (1) -----> (*) pets
users (1) -----> (1) sitters
sitters (1) ----> (*) services
bookings (*) ----> (1) users (owner)
bookings (*) ----> (1) sitters
bookings (*) ----> (1) pets
bookings (*) ----> (1) services
bookings (1) ----> (*) reviews
bookings (1) ----> (1) chats
chats (1) -------> (*) messages
## Background Workers
The application includes a background worker that runs concurrently:
Expired Booking Cleanup Worker
Automatically cancels bookings that remain in "pending" status for more than 24 hours after their scheduled start time.
Runs every hour
Uses context for graceful cancellation
Logs all operations
Handles database errors gracefully
## Docker Deployment
Services
The docker-compose.yml defines two services:
postgres - PostgreSQL 15 database
Port: 5432
Auto-initialization with schema.sql
Health checks enabled
Persistent data volume
backend - Go API server
Port: 8080
Waits for database to be healthy
Auto-restarts on failure
Multi-stage build for small image size
Project Statistics
Total Lines of Code: app 6,000 lines of Go
Modules: 6 main modules (auth, pets, bookings, reviews, services, admin)
API Endpoints: 38 endpoints
Database Tables: 9 tables with relationships
Test Files: 13 comprehensive test suites
Test Coverage: 90%+
Docker Image Size: app 20MB
Adding a New Module
If you want to add a new module,
Create module directory in internal/
Implement three layers: handler, service, repository
Define interfaces in service
Add tests for each layer
Register routes in main.go

This is an academic project for Golang Application Development course at KBTU.
# Contact
For questions or issues, please contact the team members:

1. Armankyzy Anara (Team Lead)
2. Alimzhankyzy Nuray (Scrum Master)
3. Sabukhi Raziyev (Core Backend Developer)
4. Kalikhan Arukhan (QA Engineer)




Course: Golang Application Development
Instructor: Salikh Pernebek
University: Kazakh-British Technical University
Academic Year: 2025-2026, fall semester
