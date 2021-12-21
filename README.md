# ldap2carddav
A CardDAV Server using LDAP as its backend.
It allows read only access.

## Installation
1. Grab the binary from the releases page.
2. Create a config file and save it as `ldap2carddav.yml` either in `/etc/ldap2carddav.yml` or place it next to the binary.
3. Create a systemd service (an example is provided) and adjust the path.
4. Start the service.

## Configuration
The configuration file should be named as `ldap2carddav.yml` and either placed next to the executable or in `/etc/`. The file should look like the following:
```yaml
---
ldap:
  url: ldaps://ldap.example.com
  binddn: cn=carddavservice,ou=Services,dc=example,dc=com
  bindpw: password
  basedn: ou=People,dc=example,dc=com

carddav:
  clear_old_entries: true
  address: 0.0.0.0
```

### Configurable parameters
#### LDAP
| Variable          |        Default         | Description                                                                                                                                                                                                                               |
|-------------------|:----------------------:|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `url`             |                        | Connection URL for LDAP server. For example: `ldaps://ldap.example.com`                                                                                                                                                                   |
| `binddn`          |                        | Distinguished name (DN) of the LDAP user used to bind to the server when searching for entries. For example: `cn=carddavservice,ou=Services,dc=example,dc=com`                                                                            |
| `bindpw`          |                        | The bind password to use in conjunction with the bind DN. Leave it empty for unauthenticated bind.                                                                                                                                        |
| `basedn`          |                        | The DN of the branch of the directory where all searches should start from. At the very least, this must be the top of your directory tree, but could also specify a subtree in the directory. For example: `ou=People,dc=example,dc=com` |
| `unique_id_field` |         `uid`          | The field storing the user unique identifier. Usually `uid` or `sAMAccountName`.                                                                                                                                                          |
| `phone_field`     |        `mobile`        | The field storing the user phone number. Usually `mobile` or `telephoneNumber`.                                                                                                                                                           |
| `avatar_field`    |      `jpegPhoto`       | The field storing the user avatar picture. Usually `jpegPhoto` or `thumbnailPhoto`.                                                                                                                                                       |
| `filter`          | `(objectClass=Person)` | A valid LDAP search filter.                                                                                                                                                                                                               |
| `scrapetime`      |        `6000s`         | Waiting time between refreshing objects from LDAP.                                                                                                                                                                                        |

#### CardDAV
| Variable            |       Default       | Description                                                                                         |
|---------------------|:-------------------:|-----------------------------------------------------------------------------------------------------|
| `storage_path`      | `/srv/ldap2carddav` | Paths where vCards will be stored.                                                                  |
| `subdirectory`      |       `cards`       | Name of the subdirectory where vCards will be stored.                                               |
| `address_book_name` | `LDAP address book` | Name of the address book.                                                                           |
| `clear_old_entries` |       `true`        | Clear all previous vCards when starting the service.                                                |
| `address`           |     `127.0.0.1`     | The IP address on which the application listens. Use `0.0.0.0` if you want listen on all addresses. |
| `port`              |       `8000`        | The network port on which the application listens.                                                  |
