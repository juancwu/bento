# Planing for Bento

Basic planning or roadmap for bento.

## Must have
- [ ] Be serve over HTTP request
- [ ] Must have authentication - Basic Authentication ATM
- [ ] Create new projects
- [ ] Each project can have many different compartments or environments
- [ ] Each environment can have many key-value pairs
- [ ] Update key-value pairs
- [ ] Delete key-value pairs
- [ ] Delete compartments
- [ ] Delete projects
- [ ] Download an entire `.env` with all the key-value pairs for a selected compartment of a project
- [ ] Have teams - people in the same team can have access to the kv pairs in the team projects

## Database Design

[Lucid App](https://lucid.app/lucidchart/e09306ad-5afd-476d-9a0e-9ff2f718b653/edit?invitationId=inv_79248f13-1881-494b-a7b4-cf56239081bc)

## Hosting

Bento will be hosted in fly.io and that should be enough at the moment.

## Database

Bento will be using PlanetScale database services.

Each value stored in the database should be encrypted by a secret set by the user.

## Authentication & Authorization

To ease the management for this backend service, we need to have different levels of authorization based on roles or access levels.

For authentication, since right now the only concern is to be able to use a cli that calls this service to get the key-value pairs,
there is no need for a more sophiscicated authentication pattern.

The ultimate goal is to implement oauth2 or something similar for authentication. This can be sync with the FE client
to make things easier for authenticaiton.

## Teams

Been thinking that adding the ability to add users to the same team or organization can be useful when working in a project
with multiple people in it and need the same base `.env` files.

## Authorization Based On Roles

### Owner

- All access to the team
- Add user to team
- Remove user from team
- Rename team
- Delete team
- Grant/Revoke permission to access read/write operations to certain compartments/environments

### Admin
- All of owner's priviledges except for rename/delete team

### Member
- Can have read/write priviledges depending on what is granted to the member
