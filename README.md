# gonforce
Gonforce ( Read as Go Enforce ) enforces the given package follows the guidelines defined in "gonforce.yaml" file.

Sample files may look something like the following

1. Blacklist
```
package: github.com/hasitpbhatt/gonforce

default:
  type: blacklist
  imports:
    - fmt
    - gopkg.in
  except:
    - gopkg.in/yaml.v2
```

Here we are ensuring the package can't use any imports starting with "gopkg.in" or "fmt", but with an exception "gopkg.in/yaml.v2" which would not be allowed otherwise due to "gopkg.in" being one of the blacklisted imports.

2. Whitelist
```
package: github.com/hasitpbhatt/gonforce

default:
  type: whitelist
  imports:
    - fmt
    - gopkg.in
  except:
    - gopkg.in/models
```

Here we are ensuring only imports that can be used should always start with "fmt" or "gopkg.in", but we don't want to allow "gopkg.in/models" due to some specific reason. Any other imports will automatically fail the gonforce validation.
