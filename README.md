# ldap2carddav
An Carddav Server using LDAP as its backend.
It allows read only access.

## Installation
1. Grab the binary from the releases page
2. Create a config file (see provided) and save it as `ldapcarddav.yml` either in `/etc/ldapcarddav.yml` or place it next to the binary.
3. Create a systemd service (an example is provided) and adjust the path
4. Start the service
