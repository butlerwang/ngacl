<!-----------------------------

- File Name : README.md

- Purpose :

- Creation Date : 06-29-2017

- Last Modified : Thu 29 Jun 2017 09:56:59 PM UTC

- Created By : Kiyor

------------------------------->

### Require Header

- `ccna-acl-tkmd5-secret`

### Option Header

- `ccna-acl-tkmd5-options`

#### Option Value

- `ignore_expire` value 1 or 0, default=0

## Example Nginx Config

```text
{{ if .IsEdge }}
set $ccna_acl_method "tkmd5";
set $ccna_acl_failover 403;
set $ccna_acl_tkmd5_secret "F0d9pIHAG6FPlLkr";
set $ccna_acl_tkmd5_options "ignore_expire=1";
access_by_lua_file conf/lua/acl.lua;
{{ end }}
```
